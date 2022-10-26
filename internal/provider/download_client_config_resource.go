package provider

import (
	"context"
	"fmt"
	"strconv"

	"github.com/devopsarr/terraform-provider-sonarr/internal/helpers"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"golift.io/starr/sonarr"
)

const downloadClientConfigResourceName = "download_client_config"

// Ensure provider defined types fully satisfy framework interfaces.
var _ resource.Resource = &DownloadClientConfigResource{}
var _ resource.ResourceWithImportState = &DownloadClientConfigResource{}

func NewDownloadClientConfigResource() resource.Resource {
	return &DownloadClientConfigResource{}
}

// DownloadClientConfigResource defines the download client config implementation.
type DownloadClientConfigResource struct {
	client *sonarr.Sonarr
}

// DownloadClientConfig describes the download client config data model.
type DownloadClientConfig struct {
	DownloadClientWorkingFolders    types.String `tfsdk:"download_client_working_folders"`
	ID                              types.Int64  `tfsdk:"id"`
	EnableCompletedDownloadHandling types.Bool   `tfsdk:"enable_completed_download_handling"`
	AutoRedownloadFailed            types.Bool   `tfsdk:"auto_redownload_failed"`
}

func (r *DownloadClientConfigResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_" + downloadClientConfigResourceName
}

func (r *DownloadClientConfigResource) GetSchema(ctx context.Context) (tfsdk.Schema, diag.Diagnostics) {
	return tfsdk.Schema{
		MarkdownDescription: "[subcategory:Download Clients]: #\nDownload Client Config resource.\nFor more information refer to [Download Client](https://wiki.servarr.com/sonarr/settings#completed-download-handling) documentation.",
		Attributes: map[string]tfsdk.Attribute{
			"id": {
				MarkdownDescription: "Download Client Config ID.",
				Computed:            true,
				Type:                types.Int64Type,
				PlanModifiers: tfsdk.AttributePlanModifiers{
					resource.UseStateForUnknown(),
				},
			},
			"enable_completed_download_handling": {
				MarkdownDescription: "Enable Completed Download Handling flag.",
				Required:            true,
				Type:                types.BoolType,
			},
			"auto_redownload_failed": {
				MarkdownDescription: "Auto Redownload Failed flag.",
				Required:            true,
				Type:                types.BoolType,
			},
			"download_client_working_folders": {
				MarkdownDescription: "Download Client Working Folders.",
				Computed:            true,
				Type:                types.StringType,
			},
		},
	}, nil
}

func (r *DownloadClientConfigResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	// Prevent panic if the provider has not been configured.
	if req.ProviderData == nil {
		return
	}

	client, ok := req.ProviderData.(*sonarr.Sonarr)
	if !ok {
		resp.Diagnostics.AddError(
			helpers.UnexpectedResourceConfigureType,
			fmt.Sprintf("Expected *sonarr.Sonarr, got: %T. Please report this issue to the provider developers.", req.ProviderData),
		)

		return
	}

	r.client = client
}

func (r *DownloadClientConfigResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	// Retrieve values from plan
	var config *DownloadClientConfig

	resp.Diagnostics.Append(req.Plan.Get(ctx, &config)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Build Create resource
	data := config.read()
	data.ID = 1

	// Create new DownloadClientConfig
	response, err := r.client.UpdateDownloadClientConfigContext(ctx, data)
	if err != nil {
		resp.Diagnostics.AddError(helpers.ClientError, fmt.Sprintf("Unable to create %s, got error: %s", downloadClientConfigResourceName, err))

		return
	}

	tflog.Trace(ctx, "created "+downloadClientConfigResourceName+": "+strconv.Itoa(int(response.ID)))
	// Generate resource state struct
	config.write(response)
	resp.Diagnostics.Append(resp.State.Set(ctx, &config)...)
}

func (r *DownloadClientConfigResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	// Get current state
	var config *DownloadClientConfig

	resp.Diagnostics.Append(req.State.Get(ctx, &config)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Get downloadClientConfig current value
	response, err := r.client.GetDownloadClientConfigContext(ctx)
	if err != nil {
		resp.Diagnostics.AddError(helpers.ClientError, fmt.Sprintf("Unable to read %s, got error: %s", downloadClientConfigResourceName, err))

		return
	}

	tflog.Trace(ctx, "read "+downloadClientConfigResourceName+": "+strconv.Itoa(int(response.ID)))
	// Map response body to resource schema attribute
	config.write(response)
	resp.Diagnostics.Append(resp.State.Set(ctx, &config)...)
}

func (r *DownloadClientConfigResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	// Get plan values
	var config *DownloadClientConfig

	resp.Diagnostics.Append(req.Plan.Get(ctx, &config)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Build Update resource
	data := config.read()

	// Update DownloadClientConfig
	response, err := r.client.UpdateDownloadClientConfigContext(ctx, data)
	if err != nil {
		resp.Diagnostics.AddError(helpers.ClientError, fmt.Sprintf("Unable to update %s, got error: %s", downloadClientConfigResourceName, err))

		return
	}

	tflog.Trace(ctx, "updated "+downloadClientConfigResourceName+": "+strconv.Itoa(int(response.ID)))
	// Generate resource state struct
	config.write(response)
	resp.Diagnostics.Append(resp.State.Set(ctx, config)...)
}

func (r *DownloadClientConfigResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	// DownloadClientConfig cannot be really deleted just removing configuration
	tflog.Trace(ctx, "decoupled "+downloadClientConfigResourceName+": 1")
	resp.State.RemoveResource(ctx)
}

func (r *DownloadClientConfigResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	// resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
	tflog.Trace(ctx, "imported "+downloadClientConfigResourceName+": 1")
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("id"), 1)...)
}

func (c *DownloadClientConfig) write(downloadClientConfig *sonarr.DownloadClientConfig) {
	c.EnableCompletedDownloadHandling = types.Bool{Value: downloadClientConfig.EnableCompletedDownloadHandling}
	c.AutoRedownloadFailed = types.Bool{Value: downloadClientConfig.AutoRedownloadFailed}
	c.ID = types.Int64{Value: downloadClientConfig.ID}
	c.DownloadClientWorkingFolders = types.String{Value: downloadClientConfig.DownloadClientWorkingFolders}
}

func (c *DownloadClientConfig) read() *sonarr.DownloadClientConfig {
	return &sonarr.DownloadClientConfig{
		EnableCompletedDownloadHandling: c.EnableCompletedDownloadHandling.Value,
		AutoRedownloadFailed:            c.AutoRedownloadFailed.Value,
		ID:                              c.ID.Value,
		DownloadClientWorkingFolders:    c.DownloadClientWorkingFolders.Value,
	}
}

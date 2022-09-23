package provider

import (
	"context"
	"fmt"
	"strconv"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"golift.io/starr/sonarr"
)

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
	EnableCompletedDownloadHandling types.Bool   `tfsdk:"enable_completed_download_handling"`
	AutoRedownloadFailed            types.Bool   `tfsdk:"auto_redownload_failed"`
	ID                              types.Int64  `tfsdk:"id"`
	DownloadClientWorkingFolders    types.String `tfsdk:"download_client_working_folders"`
}

func (r *DownloadClientConfigResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_download_client_config"
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
			UnexpectedResourceConfigureType,
			fmt.Sprintf("Expected *sonarr.Sonarr, got: %T. Please report this issue to the provider developers.", req.ProviderData),
		)

		return
	}

	r.client = client
}

func (r *DownloadClientConfigResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	// Retrieve values from plan
	var plan DownloadClientConfig

	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Build Create resource
	data := readDownloadClientConfig(&plan)
	data.ID = 1

	// Create new DownloadClientConfig
	response, err := r.client.UpdateDownloadClientConfigContext(ctx, data)
	if err != nil {
		resp.Diagnostics.AddError(ClientError, fmt.Sprintf("Unable to create downloadClientConfig, got error: %s", err))

		return
	}

	tflog.Trace(ctx, "created download_client_config: "+strconv.Itoa(int(response.ID)))
	// Generate resource state struct
	result := writeDownloadClientConfig(response)
	resp.Diagnostics.Append(resp.State.Set(ctx, result)...)
}

func (r *DownloadClientConfigResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	// Get current state
	var state DownloadClientConfig

	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Get downloadClientConfig current value
	response, err := r.client.GetDownloadClientConfigContext(ctx)
	if err != nil {
		resp.Diagnostics.AddError(ClientError, fmt.Sprintf("Unable to read downloadClientConfig, got error: %s", err))

		return
	}

	tflog.Trace(ctx, "read download_client_config: "+strconv.Itoa(int(response.ID)))
	// Map response body to resource schema attribute
	result := writeDownloadClientConfig(response)
	resp.Diagnostics.Append(resp.State.Set(ctx, result)...)
}

func (r *DownloadClientConfigResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	// Get plan values
	var plan DownloadClientConfig

	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Build Update resource
	data := readDownloadClientConfig(&plan)

	// Update DownloadClientConfig
	response, err := r.client.UpdateDownloadClientConfigContext(ctx, data)
	if err != nil {
		resp.Diagnostics.AddError(ClientError, fmt.Sprintf("Unable to update downloadClientConfig, got error: %s", err))

		return
	}

	tflog.Trace(ctx, "updated download_client_config: "+strconv.Itoa(int(response.ID)))
	// Generate resource state struct
	result := writeDownloadClientConfig(response)
	resp.Diagnostics.Append(resp.State.Set(ctx, result)...)
}

func (r *DownloadClientConfigResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	// DownloadClientConfig cannot be really deleted just removing configuration
	tflog.Trace(ctx, "decoupled download_client_config: 1")
	resp.State.RemoveResource(ctx)
}

func (r *DownloadClientConfigResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	// resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
	tflog.Trace(ctx, "imported download_client_config: 1")
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("id"), 1)...)
}

func writeDownloadClientConfig(downloadClientConfig *sonarr.DownloadClientConfig) *DownloadClientConfig {
	return &DownloadClientConfig{
		EnableCompletedDownloadHandling: types.Bool{Value: downloadClientConfig.EnableCompletedDownloadHandling},
		AutoRedownloadFailed:            types.Bool{Value: downloadClientConfig.AutoRedownloadFailed},
		ID:                              types.Int64{Value: downloadClientConfig.ID},
		DownloadClientWorkingFolders:    types.String{Value: downloadClientConfig.DownloadClientWorkingFolders},
	}
}

func readDownloadClientConfig(downloadClientConfig *DownloadClientConfig) *sonarr.DownloadClientConfig {
	return &sonarr.DownloadClientConfig{
		EnableCompletedDownloadHandling: downloadClientConfig.EnableCompletedDownloadHandling.Value,
		AutoRedownloadFailed:            downloadClientConfig.AutoRedownloadFailed.Value,
		ID:                              downloadClientConfig.ID.Value,
		DownloadClientWorkingFolders:    downloadClientConfig.DownloadClientWorkingFolders.Value,
	}
}

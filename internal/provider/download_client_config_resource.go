package provider

import (
	"context"
	"strconv"

	"github.com/devopsarr/sonarr-go/sonarr"
	"github.com/devopsarr/terraform-provider-sonarr/internal/helpers"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

const downloadClientConfigResourceName = "download_client_config"

// Ensure provider defined types fully satisfy framework interfaces.
var (
	_ resource.Resource                = &DownloadClientConfigResource{}
	_ resource.ResourceWithImportState = &DownloadClientConfigResource{}
)

func NewDownloadClientConfigResource() resource.Resource {
	return &DownloadClientConfigResource{}
}

// DownloadClientConfigResource defines the download client config implementation.
type DownloadClientConfigResource struct {
	client *sonarr.APIClient
}

// DownloadClientConfig describes the download client config data model.
type DownloadClientConfig struct {
	DownloadClientWorkingFolders    types.String `tfsdk:"download_client_working_folders"`
	ID                              types.Int64  `tfsdk:"id"`
	EnableCompletedDownloadHandling types.Bool   `tfsdk:"enable_completed_download_handling"`
	AutoRedownloadFailed            types.Bool   `tfsdk:"auto_redownload_failed"`
}

func (r *DownloadClientConfigResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_" + downloadClientConfigResourceName
}

func (r *DownloadClientConfigResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "<!-- subcategory:Download Clients -->Download Client Config resource.\nFor more information refer to [Download Client](https://wiki.servarr.com/sonarr/settings#completed-download-handling) documentation.",
		Attributes: map[string]schema.Attribute{
			"id": schema.Int64Attribute{
				MarkdownDescription: "Download Client Config ID.",
				Computed:            true,
				PlanModifiers: []planmodifier.Int64{
					int64planmodifier.UseStateForUnknown(),
				},
			},
			"enable_completed_download_handling": schema.BoolAttribute{
				MarkdownDescription: "Enable Completed Download Handling flag.",
				Required:            true,
			},
			"auto_redownload_failed": schema.BoolAttribute{
				MarkdownDescription: "Auto Redownload Failed flag.",
				Required:            true,
			},
			"download_client_working_folders": schema.StringAttribute{
				MarkdownDescription: "Download Client Working Folders.",
				Computed:            true,
			},
		},
	}
}

func (r *DownloadClientConfigResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	if client := helpers.ResourceConfigure(ctx, req, resp); client != nil {
		r.client = client
	}
}

func (r *DownloadClientConfigResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	// Retrieve values from plan
	var config *DownloadClientConfig

	resp.Diagnostics.Append(req.Plan.Get(ctx, &config)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Build Create resource
	request := config.read()
	request.SetId(1)

	// Create new DownloadClientConfig
	response, _, err := r.client.DownloadClientConfigAPI.UpdateDownloadClientConfig(ctx, strconv.Itoa(int(request.GetId()))).DownloadClientConfigResource(*request).Execute()
	if err != nil {
		resp.Diagnostics.AddError(helpers.ClientError, helpers.ParseClientError(helpers.Create, downloadClientConfigResourceName, err))

		return
	}

	tflog.Trace(ctx, "created "+downloadClientConfigResourceName+": "+strconv.Itoa(int(response.GetId())))
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
	response, _, err := r.client.DownloadClientConfigAPI.GetDownloadClientConfig(ctx).Execute()
	if err != nil {
		resp.Diagnostics.AddError(helpers.ClientError, helpers.ParseClientError(helpers.Read, downloadClientConfigResourceName, err))

		return
	}

	tflog.Trace(ctx, "read "+downloadClientConfigResourceName+": "+strconv.Itoa(int(response.GetId())))
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
	request := config.read()

	// Update DownloadClientConfig
	response, _, err := r.client.DownloadClientConfigAPI.UpdateDownloadClientConfig(ctx, strconv.Itoa(int(request.GetId()))).DownloadClientConfigResource(*request).Execute()
	if err != nil {
		resp.Diagnostics.AddError(helpers.ClientError, helpers.ParseClientError(helpers.Update, downloadClientConfigResourceName, err))

		return
	}

	tflog.Trace(ctx, "updated "+downloadClientConfigResourceName+": "+strconv.Itoa(int(response.GetId())))
	// Generate resource state struct
	config.write(response)
	resp.Diagnostics.Append(resp.State.Set(ctx, config)...)
}

func (r *DownloadClientConfigResource) Delete(ctx context.Context, _ resource.DeleteRequest, resp *resource.DeleteResponse) {
	// DownloadClientConfig cannot be really deleted just removing configuration
	tflog.Trace(ctx, "decoupled "+downloadClientConfigResourceName+": 1")
	resp.State.RemoveResource(ctx)
}

func (r *DownloadClientConfigResource) ImportState(ctx context.Context, _ resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	tflog.Trace(ctx, "imported "+downloadClientConfigResourceName+": 1")
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("id"), 1)...)
}

func (c *DownloadClientConfig) write(downloadClientConfig *sonarr.DownloadClientConfigResource) {
	c.EnableCompletedDownloadHandling = types.BoolValue(downloadClientConfig.GetEnableCompletedDownloadHandling())
	c.AutoRedownloadFailed = types.BoolValue(downloadClientConfig.GetAutoRedownloadFailed())
	c.ID = types.Int64Value(int64(downloadClientConfig.GetId()))
	c.DownloadClientWorkingFolders = types.StringValue(downloadClientConfig.GetDownloadClientWorkingFolders())
}

func (c *DownloadClientConfig) read() *sonarr.DownloadClientConfigResource {
	config := sonarr.NewDownloadClientConfigResource()
	config.SetEnableCompletedDownloadHandling(c.EnableCompletedDownloadHandling.ValueBool())
	config.SetAutoRedownloadFailed(c.AutoRedownloadFailed.ValueBool())
	config.SetId(int32(c.ID.ValueInt64()))
	config.SetDownloadClientWorkingFolders(c.DownloadClientWorkingFolders.ValueString())

	return config
}

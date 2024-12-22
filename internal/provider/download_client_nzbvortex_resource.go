package provider

import (
	"context"
	"strconv"

	"github.com/devopsarr/sonarr-go/sonarr"
	"github.com/devopsarr/terraform-provider-sonarr/internal/helpers"
	"github.com/hashicorp/terraform-plugin-framework-validators/int64validator"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

const (
	downloadClientNzbvortexResourceName   = "download_client_nzbvortex"
	downloadClientNzbvortexImplementation = "Nzbvortex"
	downloadClientNzbvortexConfigContract = "NzbvortexSettings"
	downloadClientNzbvortexProtocol       = "usenet"
)

// Ensure provider defined types fully satisfy framework interfaces.
var (
	_ resource.Resource                = &DownloadClientNzbvortexResource{}
	_ resource.ResourceWithImportState = &DownloadClientNzbvortexResource{}
)

func NewDownloadClientNzbvortexResource() resource.Resource {
	return &DownloadClientNzbvortexResource{}
}

// DownloadClientNzbvortexResource defines the download client implementation.
type DownloadClientNzbvortexResource struct {
	client *sonarr.APIClient
	auth   context.Context
}

// DownloadClientNzbvortex describes the download client data model.
type DownloadClientNzbvortex struct {
	Tags                     types.Set    `tfsdk:"tags"`
	Name                     types.String `tfsdk:"name"`
	Host                     types.String `tfsdk:"host"`
	URLBase                  types.String `tfsdk:"url_base"`
	APIKey                   types.String `tfsdk:"api_key"`
	TvCategory               types.String `tfsdk:"tv_category"`
	RecentTvPriority         types.Int64  `tfsdk:"recent_tv_priority"`
	OlderTvPriority          types.Int64  `tfsdk:"older_tv_priority"`
	Priority                 types.Int64  `tfsdk:"priority"`
	Port                     types.Int64  `tfsdk:"port"`
	ID                       types.Int64  `tfsdk:"id"`
	Enable                   types.Bool   `tfsdk:"enable"`
	RemoveFailedDownloads    types.Bool   `tfsdk:"remove_failed_downloads"`
	RemoveCompletedDownloads types.Bool   `tfsdk:"remove_completed_downloads"`
}

func (d DownloadClientNzbvortex) toDownloadClient() *DownloadClient {
	return &DownloadClient{
		Tags:                     d.Tags,
		Name:                     d.Name,
		Host:                     d.Host,
		URLBase:                  d.URLBase,
		APIKey:                   d.APIKey,
		TvCategory:               d.TvCategory,
		RecentTvPriority:         d.RecentTvPriority,
		OlderTvPriority:          d.OlderTvPriority,
		Priority:                 d.Priority,
		Port:                     d.Port,
		ID:                       d.ID,
		Enable:                   d.Enable,
		RemoveFailedDownloads:    d.RemoveFailedDownloads,
		RemoveCompletedDownloads: d.RemoveCompletedDownloads,
		Implementation:           types.StringValue(downloadClientNzbvortexImplementation),
		ConfigContract:           types.StringValue(downloadClientNzbvortexConfigContract),
		Protocol:                 types.StringValue(downloadClientNzbvortexProtocol),
	}
}

func (d *DownloadClientNzbvortex) fromDownloadClient(client *DownloadClient) {
	d.Tags = client.Tags
	d.Name = client.Name
	d.Host = client.Host
	d.URLBase = client.URLBase
	d.APIKey = client.APIKey
	d.TvCategory = client.TvCategory
	d.RecentTvPriority = client.RecentTvPriority
	d.OlderTvPriority = client.OlderTvPriority
	d.Priority = client.Priority
	d.Port = client.Port
	d.ID = client.ID
	d.Enable = client.Enable
	d.RemoveFailedDownloads = client.RemoveFailedDownloads
	d.RemoveCompletedDownloads = client.RemoveCompletedDownloads
}

func (r *DownloadClientNzbvortexResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_" + downloadClientNzbvortexResourceName
}

func (r *DownloadClientNzbvortexResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "<!-- subcategory:Download Clients -->\nDownload Client Nzbvortex resource.\nFor more information refer to [Download Client](https://wiki.servarr.com/sonarr/settings#download-clients) and [Nzbvortex](https://wiki.servarr.com/sonarr/supported#nzbvortex).",
		Attributes: map[string]schema.Attribute{
			"enable": schema.BoolAttribute{
				MarkdownDescription: "Enable flag.",
				Optional:            true,
				Computed:            true,
			},
			"remove_completed_downloads": schema.BoolAttribute{
				MarkdownDescription: "Remove completed downloads flag.",
				Optional:            true,
				Computed:            true,
			},
			"remove_failed_downloads": schema.BoolAttribute{
				MarkdownDescription: "Remove failed downloads flag.",
				Optional:            true,
				Computed:            true,
			},
			"priority": schema.Int64Attribute{
				MarkdownDescription: "Priority.",
				Optional:            true,
				Computed:            true,
			},
			"name": schema.StringAttribute{
				MarkdownDescription: "Download Client name.",
				Required:            true,
			},
			"tags": schema.SetAttribute{
				MarkdownDescription: "List of associated tags.",
				Optional:            true,
				Computed:            true,
				ElementType:         types.Int64Type,
			},
			"id": schema.Int64Attribute{
				MarkdownDescription: "Download Client ID.",
				Computed:            true,
				PlanModifiers: []planmodifier.Int64{
					int64planmodifier.UseStateForUnknown(),
				},
			},
			// Field values
			"port": schema.Int64Attribute{
				MarkdownDescription: "Port.",
				Optional:            true,
				Computed:            true,
			},
			"recent_tv_priority": schema.Int64Attribute{
				MarkdownDescription: "Recent TV priority. `-1` Low, `0` Normal, `1` High.",
				Optional:            true,
				Computed:            true,
				Validators: []validator.Int64{
					int64validator.OneOf(-1, 0, 1),
				},
			},
			"older_tv_priority": schema.Int64Attribute{
				MarkdownDescription: "Older TV priority. `-1` Low, `0` Normal, `1` High.",
				Optional:            true,
				Computed:            true,
				Validators: []validator.Int64{
					int64validator.OneOf(-1, 0, 1),
				},
			},
			"host": schema.StringAttribute{
				MarkdownDescription: "host.",
				Optional:            true,
				Computed:            true,
			},
			"url_base": schema.StringAttribute{
				MarkdownDescription: "Base URL.",
				Optional:            true,
				Computed:            true,
			},
			"api_key": schema.StringAttribute{
				MarkdownDescription: "API key.",
				Required:            true,
				Sensitive:           true,
			},
			"tv_category": schema.StringAttribute{
				MarkdownDescription: "TV category.",
				Optional:            true,
				Computed:            true,
			},
		},
	}
}

func (r *DownloadClientNzbvortexResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	if auth, client := resourceConfigure(ctx, req, resp); client != nil {
		r.client = client
		r.auth = auth
	}
}

func (r *DownloadClientNzbvortexResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	// Retrieve values from plan
	var client *DownloadClientNzbvortex

	resp.Diagnostics.Append(req.Plan.Get(ctx, &client)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Create new DownloadClientNzbvortex
	request := client.read(ctx, &resp.Diagnostics)

	response, _, err := r.client.DownloadClientAPI.CreateDownloadClient(r.auth).DownloadClientResource(*request).Execute()
	if err != nil {
		resp.Diagnostics.AddError(helpers.ClientError, helpers.ParseClientError(helpers.Create, downloadClientNzbvortexResourceName, err))

		return
	}

	tflog.Trace(ctx, "created "+downloadClientNzbvortexResourceName+": "+strconv.Itoa(int(response.GetId())))
	// Generate resource state struct
	client.write(ctx, response, &resp.Diagnostics)
	resp.Diagnostics.Append(resp.State.Set(ctx, &client)...)
}

func (r *DownloadClientNzbvortexResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	// Get current state
	var client DownloadClientNzbvortex

	resp.Diagnostics.Append(req.State.Get(ctx, &client)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Get DownloadClientNzbvortex current value
	response, _, err := r.client.DownloadClientAPI.GetDownloadClientById(r.auth, int32(client.ID.ValueInt64())).Execute()
	if err != nil {
		resp.Diagnostics.AddError(helpers.ClientError, helpers.ParseClientError(helpers.Read, downloadClientNzbvortexResourceName, err))

		return
	}

	tflog.Trace(ctx, "read "+downloadClientNzbvortexResourceName+": "+strconv.Itoa(int(response.GetId())))
	// Map response body to resource schema attribute
	client.write(ctx, response, &resp.Diagnostics)
	resp.Diagnostics.Append(resp.State.Set(ctx, &client)...)
}

func (r *DownloadClientNzbvortexResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	// Get plan values
	var client *DownloadClientNzbvortex

	resp.Diagnostics.Append(req.Plan.Get(ctx, &client)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Update DownloadClientNzbvortex
	request := client.read(ctx, &resp.Diagnostics)

	response, _, err := r.client.DownloadClientAPI.UpdateDownloadClient(r.auth, request.GetId()).DownloadClientResource(*request).Execute()
	if err != nil {
		resp.Diagnostics.AddError(helpers.ClientError, helpers.ParseClientError(helpers.Update, downloadClientNzbvortexResourceName, err))

		return
	}

	tflog.Trace(ctx, "updated "+downloadClientNzbvortexResourceName+": "+strconv.Itoa(int(response.GetId())))
	// Generate resource state struct
	client.write(ctx, response, &resp.Diagnostics)
	resp.Diagnostics.Append(resp.State.Set(ctx, &client)...)
}

func (r *DownloadClientNzbvortexResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var ID int64

	resp.Diagnostics.Append(req.State.GetAttribute(ctx, path.Root("id"), &ID)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Delete DownloadClientNzbvortex current value
	_, err := r.client.DownloadClientAPI.DeleteDownloadClient(r.auth, int32(ID)).Execute()
	if err != nil {
		resp.Diagnostics.AddError(helpers.ClientError, helpers.ParseClientError(helpers.Delete, downloadClientNzbvortexResourceName, err))

		return
	}

	tflog.Trace(ctx, "deleted "+downloadClientNzbvortexResourceName+": "+strconv.Itoa(int(ID)))
	resp.State.RemoveResource(ctx)
}

func (r *DownloadClientNzbvortexResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	helpers.ImportStatePassthroughIntID(ctx, path.Root("id"), req, resp)
	tflog.Trace(ctx, "imported "+downloadClientNzbvortexResourceName+": "+req.ID)
}

func (d *DownloadClientNzbvortex) write(ctx context.Context, downloadClient *sonarr.DownloadClientResource, diags *diag.Diagnostics) {
	genericDownloadClient := d.toDownloadClient()
	genericDownloadClient.write(ctx, downloadClient, diags)
	d.fromDownloadClient(genericDownloadClient)
}

func (d *DownloadClientNzbvortex) read(ctx context.Context, diags *diag.Diagnostics) *sonarr.DownloadClientResource {
	return d.toDownloadClient().read(ctx, diags)
}

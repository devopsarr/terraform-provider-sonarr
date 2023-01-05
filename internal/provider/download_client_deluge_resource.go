package provider

import (
	"context"
	"fmt"
	"strconv"

	"github.com/devopsarr/sonarr-go/sonarr"
	"github.com/devopsarr/terraform-provider-sonarr/tools"
	"github.com/hashicorp/terraform-plugin-framework-validators/int64validator"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

const (
	downloadClientDelugeResourceName   = "download_client_deluge"
	downloadClientDelugeImplementation = "Deluge"
	downloadClientDelugeConfigContract = "DelugeSettings"
	downloadClientDelugeProtocol       = "torrent"
)

// Ensure provider defined types fully satisfy framework interfaces.
var (
	_ resource.Resource                = &DownloadClientDelugeResource{}
	_ resource.ResourceWithImportState = &DownloadClientDelugeResource{}
)

func NewDownloadClientDelugeResource() resource.Resource {
	return &DownloadClientDelugeResource{}
}

// DownloadClientDelugeResource defines the download client implementation.
type DownloadClientDelugeResource struct {
	client *sonarr.APIClient
}

// DownloadClientDeluge describes the download client data model.
type DownloadClientDeluge struct {
	Tags                     types.Set    `tfsdk:"tags"`
	Name                     types.String `tfsdk:"name"`
	Host                     types.String `tfsdk:"host"`
	URLBase                  types.String `tfsdk:"url_base"`
	Username                 types.String `tfsdk:"username"`
	Password                 types.String `tfsdk:"password"`
	TvCategory               types.String `tfsdk:"tv_category"`
	TvDirectory              types.String `tfsdk:"tv_directory"`
	TvImportedCategory       types.String `tfsdk:"tv_imported_category"`
	RecentTvPriority         types.Int64  `tfsdk:"recent_tv_priority"`
	OlderTvPriority          types.Int64  `tfsdk:"older_tv_priority"`
	Priority                 types.Int64  `tfsdk:"priority"`
	Port                     types.Int64  `tfsdk:"port"`
	ID                       types.Int64  `tfsdk:"id"`
	AddPaused                types.Bool   `tfsdk:"add_paused"`
	UseSsl                   types.Bool   `tfsdk:"use_ssl"`
	Enable                   types.Bool   `tfsdk:"enable"`
	RemoveFailedDownloads    types.Bool   `tfsdk:"remove_failed_downloads"`
	RemoveCompletedDownloads types.Bool   `tfsdk:"remove_completed_downloads"`
}

func (d DownloadClientDeluge) toDownloadClient() *DownloadClient {
	return &DownloadClient{
		Tags:                     d.Tags,
		Name:                     d.Name,
		Host:                     d.Host,
		URLBase:                  d.URLBase,
		Username:                 d.Username,
		Password:                 d.Password,
		TvCategory:               d.TvCategory,
		TvDirectory:              d.TvDirectory,
		TvImportedCategory:       d.TvImportedCategory,
		RecentTvPriority:         d.RecentTvPriority,
		OlderTvPriority:          d.OlderTvPriority,
		Priority:                 d.Priority,
		Port:                     d.Port,
		ID:                       d.ID,
		AddPaused:                d.AddPaused,
		UseSsl:                   d.UseSsl,
		Enable:                   d.Enable,
		RemoveFailedDownloads:    d.RemoveFailedDownloads,
		RemoveCompletedDownloads: d.RemoveCompletedDownloads,
	}
}

func (d *DownloadClientDeluge) fromDownloadClient(client *DownloadClient) {
	d.Tags = client.Tags
	d.Name = client.Name
	d.Host = client.Host
	d.URLBase = client.URLBase
	d.Username = client.Username
	d.Password = client.Password
	d.TvCategory = client.TvCategory
	d.TvDirectory = client.TvDirectory
	d.TvImportedCategory = client.TvImportedCategory
	d.RecentTvPriority = client.RecentTvPriority
	d.OlderTvPriority = client.OlderTvPriority
	d.Priority = client.Priority
	d.Port = client.Port
	d.ID = client.ID
	d.AddPaused = client.AddPaused
	d.UseSsl = client.UseSsl
	d.Enable = client.Enable
	d.RemoveFailedDownloads = client.RemoveFailedDownloads
	d.RemoveCompletedDownloads = client.RemoveCompletedDownloads
}

func (r *DownloadClientDelugeResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_" + downloadClientDelugeResourceName
}

func (r *DownloadClientDelugeResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "<!-- subcategory:Download Clients -->Download Client Deluge resource.\nFor more information refer to [Download Client](https://wiki.servarr.com/sonarr/settings#download-clients) and [Deluge](https://wiki.servarr.com/sonarr/supported#deluge).",
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
			"add_paused": schema.BoolAttribute{
				MarkdownDescription: "Add paused flag.",
				Optional:            true,
				Computed:            true,
			},
			"use_ssl": schema.BoolAttribute{
				MarkdownDescription: "Use SSL flag.",
				Optional:            true,
				Computed:            true,
			},
			"port": schema.Int64Attribute{
				MarkdownDescription: "Port.",
				Optional:            true,
				Computed:            true,
			},
			"recent_tv_priority": schema.Int64Attribute{
				MarkdownDescription: "Recent TV priority. `0` Last, `1` First.",
				Optional:            true,
				Computed:            true,
				Validators: []validator.Int64{
					int64validator.OneOf(0, 1),
				},
			},
			"older_tv_priority": schema.Int64Attribute{
				MarkdownDescription: "Older TV priority. `0` Last, `1` First.",
				Optional:            true,
				Computed:            true,
				Validators: []validator.Int64{
					int64validator.OneOf(0, 1),
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
			"username": schema.StringAttribute{
				MarkdownDescription: "Username.",
				Optional:            true,
				Computed:            true,
			},
			"password": schema.StringAttribute{
				MarkdownDescription: "Password.",
				Optional:            true,
				Computed:            true,
				Sensitive:           true,
			},
			"tv_category": schema.StringAttribute{
				MarkdownDescription: "TV category.",
				Optional:            true,
				Computed:            true,
			},
			"tv_directory": schema.StringAttribute{
				MarkdownDescription: "TV directory.",
				Optional:            true,
				Computed:            true,
			},
			"tv_imported_category": schema.StringAttribute{
				MarkdownDescription: "TV imported category.",
				Optional:            true,
				Computed:            true,
			},
		},
	}
}

func (r *DownloadClientDelugeResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	// Prevent panic if the provider has not been configured.
	if req.ProviderData == nil {
		return
	}

	client, ok := req.ProviderData.(*sonarr.APIClient)
	if !ok {
		resp.Diagnostics.AddError(
			tools.UnexpectedResourceConfigureType,
			fmt.Sprintf("Expected *sonarr.APIClient, got: %T. Please report this issue to the provider developers.", req.ProviderData),
		)

		return
	}

	r.client = client
}

func (r *DownloadClientDelugeResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	// Retrieve values from plan
	var client *DownloadClientDeluge

	resp.Diagnostics.Append(req.Plan.Get(ctx, &client)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Create new DownloadClientDeluge
	request := client.read(ctx)

	response, _, err := r.client.DownloadClientApi.CreateDownloadclient(ctx).DownloadClientResource(*request).Execute()
	if err != nil {
		resp.Diagnostics.AddError(tools.ClientError, fmt.Sprintf("Unable to create %s, got error: %s", downloadClientDelugeResourceName, err))

		return
	}

	tflog.Trace(ctx, "created "+downloadClientDelugeResourceName+": "+strconv.Itoa(int(response.GetId())))
	// Generate resource state struct
	client.write(ctx, response)
	resp.Diagnostics.Append(resp.State.Set(ctx, &client)...)
}

func (r *DownloadClientDelugeResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	// Get current state
	var client DownloadClientDeluge

	resp.Diagnostics.Append(req.State.Get(ctx, &client)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Get DownloadClientDeluge current value
	response, _, err := r.client.DownloadClientApi.GetDownloadclientById(ctx, int32(client.ID.ValueInt64())).Execute()
	if err != nil {
		resp.Diagnostics.AddError(tools.ClientError, fmt.Sprintf("Unable to read %s, got error: %s", downloadClientDelugeResourceName, err))

		return
	}

	tflog.Trace(ctx, "read "+downloadClientDelugeResourceName+": "+strconv.Itoa(int(response.GetId())))
	// Map response body to resource schema attribute
	client.write(ctx, response)
	resp.Diagnostics.Append(resp.State.Set(ctx, &client)...)
}

func (r *DownloadClientDelugeResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	// Get plan values
	var client *DownloadClientDeluge

	resp.Diagnostics.Append(req.Plan.Get(ctx, &client)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Update DownloadClientDeluge
	request := client.read(ctx)

	response, _, err := r.client.DownloadClientApi.UpdateDownloadclient(ctx, strconv.Itoa(int(request.GetId()))).DownloadClientResource(*request).Execute()
	if err != nil {
		resp.Diagnostics.AddError(tools.ClientError, fmt.Sprintf("Unable to update %s, got error: %s", downloadClientDelugeResourceName, err))

		return
	}

	tflog.Trace(ctx, "updated "+downloadClientDelugeResourceName+": "+strconv.Itoa(int(response.GetId())))
	// Generate resource state struct
	client.write(ctx, response)
	resp.Diagnostics.Append(resp.State.Set(ctx, &client)...)
}

func (r *DownloadClientDelugeResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var client *DownloadClientDeluge

	resp.Diagnostics.Append(req.State.Get(ctx, &client)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Delete DownloadClientDeluge current value
	_, err := r.client.DownloadClientApi.DeleteDownloadclient(ctx, int32(client.ID.ValueInt64())).Execute()
	if err != nil {
		resp.Diagnostics.AddError(tools.ClientError, fmt.Sprintf("Unable to read %s, got error: %s", downloadClientDelugeResourceName, err))

		return
	}

	tflog.Trace(ctx, "deleted "+downloadClientDelugeResourceName+": "+strconv.Itoa(int(client.ID.ValueInt64())))
	resp.State.RemoveResource(ctx)
}

func (r *DownloadClientDelugeResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	// resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
	id, err := strconv.Atoi(req.ID)
	if err != nil {
		resp.Diagnostics.AddError(
			tools.UnexpectedImportIdentifier,
			fmt.Sprintf("Expected import identifier with format: ID. Got: %q", req.ID),
		)

		return
	}

	tflog.Trace(ctx, "imported "+downloadClientDelugeResourceName+": "+strconv.Itoa(id))
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("id"), id)...)
}

func (d *DownloadClientDeluge) write(ctx context.Context, downloadClient *sonarr.DownloadClientResource) {
	genericDownloadClient := DownloadClient{
		Enable:                   types.BoolValue(downloadClient.GetEnable()),
		RemoveCompletedDownloads: types.BoolValue(downloadClient.GetRemoveCompletedDownloads()),
		RemoveFailedDownloads:    types.BoolValue(downloadClient.GetRemoveFailedDownloads()),
		Priority:                 types.Int64Value(int64(downloadClient.GetPriority())),
		ID:                       types.Int64Value(int64(downloadClient.GetId())),
		Name:                     types.StringValue(downloadClient.GetName()),
	}
	genericDownloadClient.Tags, _ = types.SetValueFrom(ctx, types.Int64Type, downloadClient.Tags)
	genericDownloadClient.writeFields(ctx, downloadClient.Fields)
	d.fromDownloadClient(&genericDownloadClient)
}

func (d *DownloadClientDeluge) read(ctx context.Context) *sonarr.DownloadClientResource {
	var tags []*int32

	tfsdk.ValueAs(ctx, d.Tags, &tags)

	client := sonarr.NewDownloadClientResource()
	client.SetEnable(d.Enable.ValueBool())
	client.SetRemoveCompletedDownloads(d.RemoveCompletedDownloads.ValueBool())
	client.SetRemoveFailedDownloads(d.RemoveFailedDownloads.ValueBool())
	client.SetPriority(int32(d.Priority.ValueInt64()))
	client.SetId(int32(d.ID.ValueInt64()))
	client.SetConfigContract(downloadClientDelugeConfigContract)
	client.SetImplementation(downloadClientDelugeImplementation)
	client.SetName(d.Name.ValueString())
	client.SetProtocol(downloadClientDelugeProtocol)
	client.SetTags(tags)
	client.SetFields(d.toDownloadClient().readFields(ctx))

	return client
}

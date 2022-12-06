package provider

import (
	"context"
	"fmt"
	"strconv"

	"github.com/devopsarr/terraform-provider-sonarr/tools"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"golift.io/starr/sonarr"
)

const (
	downloadClientTorrentBlackholeResourceName   = "download_client_torrent_blackhole"
	DownloadClientTorrentBlackholeImplementation = "TorrentBlackhole"
	DownloadClientTorrentBlackholeConfigContrat  = "TorrentBlackholeSettings"
	DownloadClientTorrentBlackholeProtocol       = "torrent"
)

// Ensure provider defined types fully satisfy framework interfaces.
var _ resource.Resource = &DownloadClientTorrentBlackholeResource{}
var _ resource.ResourceWithImportState = &DownloadClientTorrentBlackholeResource{}

func NewDownloadClientTorrentBlackholeResource() resource.Resource {
	return &DownloadClientTorrentBlackholeResource{}
}

// DownloadClientTorrentBlackholeResource defines the download client implementation.
type DownloadClientTorrentBlackholeResource struct {
	client *sonarr.Sonarr
}

// DownloadClientTorrentBlackhole describes the download client data model.
type DownloadClientTorrentBlackhole struct {
	Tags                     types.Set    `tfsdk:"tags"`
	Name                     types.String `tfsdk:"name"`
	TorrentFolder            types.String `tfsdk:"torrent_folder"`
	WatchFolder              types.String `tfsdk:"watch_folder"`
	MagnetFileExtension      types.String `tfsdk:"magnet_file_extension"`
	Priority                 types.Int64  `tfsdk:"priority"`
	ID                       types.Int64  `tfsdk:"id"`
	Enable                   types.Bool   `tfsdk:"enable"`
	RemoveFailedDownloads    types.Bool   `tfsdk:"remove_failed_downloads"`
	RemoveCompletedDownloads types.Bool   `tfsdk:"remove_completed_downloads"`
	SaveMagnetFiles          types.Bool   `tfsdk:"save_magnet_files"`
	ReadOnly                 types.Bool   `tfsdk:"read_only"`
}

func (d DownloadClientTorrentBlackhole) toDownloadClient() *DownloadClient {
	return &DownloadClient{
		Tags:                     d.Tags,
		Name:                     d.Name,
		TorrentFolder:            d.TorrentFolder,
		WatchFolder:              d.WatchFolder,
		MagnetFileExtension:      d.MagnetFileExtension,
		Priority:                 d.Priority,
		ID:                       d.ID,
		Enable:                   d.Enable,
		RemoveFailedDownloads:    d.RemoveFailedDownloads,
		RemoveCompletedDownloads: d.RemoveCompletedDownloads,
		SaveMagnetFiles:          d.SaveMagnetFiles,
		ReadOnly:                 d.ReadOnly,
	}
}

func (d *DownloadClientTorrentBlackhole) fromDownloadClient(client *DownloadClient) {
	d.Tags = client.Tags
	d.Name = client.Name
	d.TorrentFolder = client.TorrentFolder
	d.WatchFolder = client.WatchFolder
	d.MagnetFileExtension = client.MagnetFileExtension
	d.Priority = client.Priority
	d.ID = client.ID
	d.Enable = client.Enable
	d.RemoveFailedDownloads = client.RemoveFailedDownloads
	d.RemoveCompletedDownloads = client.RemoveCompletedDownloads
	d.SaveMagnetFiles = client.SaveMagnetFiles
	d.ReadOnly = client.ReadOnly
}

func (r *DownloadClientTorrentBlackholeResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_" + downloadClientTorrentBlackholeResourceName
}

func (r *DownloadClientTorrentBlackholeResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "<!-- subcategory:Download Clients -->Download Client Torrent Blackhole resource.\nFor more information refer to [Download Client](https://wiki.servarr.com/sonarr/settings#download-clients) and [TorrentBlackhole](https://wiki.servarr.com/sonarr/supported#torrentblackhole).",
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
			"save_magnet_files": schema.BoolAttribute{
				MarkdownDescription: "Save magnet files flag.",
				Optional:            true,
				Computed:            true,
			},
			"read_only": schema.BoolAttribute{
				MarkdownDescription: "Read only flag.",
				Optional:            true,
				Computed:            true,
			},
			"torrent_folder": schema.StringAttribute{
				MarkdownDescription: "Torrent folder.",
				Required:            true,
			},
			"watch_folder": schema.StringAttribute{
				MarkdownDescription: "Watch folder flag.",
				Required:            true,
			},
			"magnet_file_extension": schema.StringAttribute{
				MarkdownDescription: "Magnet file extension.",
				Optional:            true,
				Computed:            true,
			},
		},
	}
}

func (r *DownloadClientTorrentBlackholeResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	// Prevent panic if the provider has not been configured.
	if req.ProviderData == nil {
		return
	}

	client, ok := req.ProviderData.(*sonarr.Sonarr)
	if !ok {
		resp.Diagnostics.AddError(
			tools.UnexpectedResourceConfigureType,
			fmt.Sprintf("Expected *sonarr.Sonarr, got: %T. Please report this issue to the provider developers.", req.ProviderData),
		)

		return
	}

	r.client = client
}

func (r *DownloadClientTorrentBlackholeResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	// Retrieve values from plan
	var client *DownloadClientTorrentBlackhole

	resp.Diagnostics.Append(req.Plan.Get(ctx, &client)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Create new DownloadClientTorrentBlackhole
	request := client.read(ctx)

	response, err := r.client.AddDownloadClientContext(ctx, request)
	if err != nil {
		resp.Diagnostics.AddError(tools.ClientError, fmt.Sprintf("Unable to create %s, got error: %s", downloadClientTorrentBlackholeResourceName, err))

		return
	}

	tflog.Trace(ctx, "created "+downloadClientTorrentBlackholeResourceName+": "+strconv.Itoa(int(response.ID)))
	// Generate resource state struct
	client.write(ctx, response)
	resp.Diagnostics.Append(resp.State.Set(ctx, &client)...)
}

func (r *DownloadClientTorrentBlackholeResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	// Get current state
	var client DownloadClientTorrentBlackhole

	resp.Diagnostics.Append(req.State.Get(ctx, &client)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Get DownloadClientTorrentBlackhole current value
	response, err := r.client.GetDownloadClientContext(ctx, client.ID.ValueInt64())
	if err != nil {
		resp.Diagnostics.AddError(tools.ClientError, fmt.Sprintf("Unable to read %s, got error: %s", downloadClientTorrentBlackholeResourceName, err))

		return
	}

	tflog.Trace(ctx, "read "+downloadClientTorrentBlackholeResourceName+": "+strconv.Itoa(int(response.ID)))
	// Map response body to resource schema attribute
	client.write(ctx, response)
	resp.Diagnostics.Append(resp.State.Set(ctx, &client)...)
}

func (r *DownloadClientTorrentBlackholeResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	// Get plan values
	var client *DownloadClientTorrentBlackhole

	resp.Diagnostics.Append(req.Plan.Get(ctx, &client)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Update DownloadClientTorrentBlackhole
	request := client.read(ctx)

	response, err := r.client.UpdateDownloadClientContext(ctx, request)
	if err != nil {
		resp.Diagnostics.AddError(tools.ClientError, fmt.Sprintf("Unable to update %s, got error: %s", downloadClientTorrentBlackholeResourceName, err))

		return
	}

	tflog.Trace(ctx, "updated "+downloadClientTorrentBlackholeResourceName+": "+strconv.Itoa(int(response.ID)))
	// Generate resource state struct
	client.write(ctx, response)
	resp.Diagnostics.Append(resp.State.Set(ctx, &client)...)
}

func (r *DownloadClientTorrentBlackholeResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var client *DownloadClientTorrentBlackhole

	resp.Diagnostics.Append(req.State.Get(ctx, &client)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Delete DownloadClientTorrentBlackhole current value
	err := r.client.DeleteDownloadClientContext(ctx, client.ID.ValueInt64())
	if err != nil {
		resp.Diagnostics.AddError(tools.ClientError, fmt.Sprintf("Unable to read %s, got error: %s", downloadClientTorrentBlackholeResourceName, err))

		return
	}

	tflog.Trace(ctx, "deleted "+downloadClientTorrentBlackholeResourceName+": "+strconv.Itoa(int(client.ID.ValueInt64())))
	resp.State.RemoveResource(ctx)
}

func (r *DownloadClientTorrentBlackholeResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	// resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
	id, err := strconv.Atoi(req.ID)
	if err != nil {
		resp.Diagnostics.AddError(
			tools.UnexpectedImportIdentifier,
			fmt.Sprintf("Expected import identifier with format: ID. Got: %q", req.ID),
		)

		return
	}

	tflog.Trace(ctx, "imported "+downloadClientTorrentBlackholeResourceName+": "+strconv.Itoa(id))
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("id"), id)...)
}

func (d *DownloadClientTorrentBlackhole) write(ctx context.Context, downloadClient *sonarr.DownloadClientOutput) {
	genericDownloadClient := DownloadClient{
		Enable:                   types.BoolValue(downloadClient.Enable),
		RemoveCompletedDownloads: types.BoolValue(downloadClient.RemoveCompletedDownloads),
		RemoveFailedDownloads:    types.BoolValue(downloadClient.RemoveFailedDownloads),
		Priority:                 types.Int64Value(int64(downloadClient.Priority)),
		ID:                       types.Int64Value(downloadClient.ID),
		Name:                     types.StringValue(downloadClient.Name),
		Tags:                     types.SetNull(types.Int64Type),
	}
	if !d.Tags.IsNull() || !(len(downloadClient.Tags) == 0) {
		genericDownloadClient.Tags, _ = types.SetValueFrom(ctx, types.Int64Type, downloadClient.Tags)
	}

	genericDownloadClient.writeFields(ctx, downloadClient.Fields)
	d.fromDownloadClient(&genericDownloadClient)
}

func (d *DownloadClientTorrentBlackhole) read(ctx context.Context) *sonarr.DownloadClientInput {
	var tags []int

	tfsdk.ValueAs(ctx, d.Tags, &tags)

	return &sonarr.DownloadClientInput{
		Enable:                   d.Enable.ValueBool(),
		RemoveCompletedDownloads: d.RemoveCompletedDownloads.ValueBool(),
		RemoveFailedDownloads:    d.RemoveFailedDownloads.ValueBool(),
		Priority:                 int(d.Priority.ValueInt64()),
		ID:                       d.ID.ValueInt64(),
		ConfigContract:           DownloadClientTorrentBlackholeConfigContrat,
		Implementation:           DownloadClientTorrentBlackholeImplementation,
		Name:                     d.Name.ValueString(),
		Protocol:                 DownloadClientTorrentBlackholeProtocol,
		Tags:                     tags,
		Fields:                   d.toDownloadClient().readFields(ctx),
	}
}

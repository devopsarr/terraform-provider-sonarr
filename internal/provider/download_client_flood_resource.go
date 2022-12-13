package provider

import (
	"context"
	"fmt"
	"strconv"

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
	"golift.io/starr/sonarr"
)

const (
	downloadClientFloodResourceName   = "download_client_flood"
	downloadClientFloodImplementation = "Flood"
	downloadClientFloodConfigContrat  = "FloodSettings"
	DownloadClientFloodProtocol       = "torrent"
)

// Ensure provider defined types fully satisfy framework interfaces.
var (
	_ resource.Resource                = &DownloadClientFloodResource{}
	_ resource.ResourceWithImportState = &DownloadClientFloodResource{}
)

func NewDownloadClientFloodResource() resource.Resource {
	return &DownloadClientFloodResource{}
}

// DownloadClientFloodResource defines the download client implementation.
type DownloadClientFloodResource struct {
	client *sonarr.Sonarr
}

// DownloadClientFlood describes the download client data model.
type DownloadClientFlood struct {
	Tags                     types.Set    `tfsdk:"tags"`
	FieldTags                types.Set    `tfsdk:"field_tags"`
	AdditionalTags           types.Set    `tfsdk:"additional_tags"`
	PostImportTags           types.Set    `tfsdk:"post_import_tags"`
	Name                     types.String `tfsdk:"name"`
	Host                     types.String `tfsdk:"host"`
	URLBase                  types.String `tfsdk:"url_base"`
	Username                 types.String `tfsdk:"username"`
	Password                 types.String `tfsdk:"password"`
	Destination              types.String `tfsdk:"destination"`
	TvCategory               types.String `tfsdk:"tv_category"`
	TvDirectory              types.String `tfsdk:"tv_directory"`
	RecentTvPriority         types.Int64  `tfsdk:"recent_tv_priority"`
	OlderTvPriority          types.Int64  `tfsdk:"older_tv_priority"`
	Priority                 types.Int64  `tfsdk:"priority"`
	Port                     types.Int64  `tfsdk:"port"`
	ID                       types.Int64  `tfsdk:"id"`
	StartOnAdd               types.Bool   `tfsdk:"start_on_add"`
	UseSsl                   types.Bool   `tfsdk:"use_ssl"`
	Enable                   types.Bool   `tfsdk:"enable"`
	RemoveFailedDownloads    types.Bool   `tfsdk:"remove_failed_downloads"`
	RemoveCompletedDownloads types.Bool   `tfsdk:"remove_completed_downloads"`
}

func (d DownloadClientFlood) toDownloadClient() *DownloadClient {
	return &DownloadClient{
		Tags:                     d.Tags,
		FieldTags:                d.FieldTags,
		PostImportTags:           d.PostImportTags,
		AdditionalTags:           d.AdditionalTags,
		Name:                     d.Name,
		Host:                     d.Host,
		URLBase:                  d.URLBase,
		Username:                 d.Username,
		Password:                 d.Password,
		Destination:              d.Destination,
		TvCategory:               d.TvCategory,
		TvDirectory:              d.TvDirectory,
		RecentTvPriority:         d.RecentTvPriority,
		OlderTvPriority:          d.OlderTvPriority,
		Priority:                 d.Priority,
		Port:                     d.Port,
		ID:                       d.ID,
		StartOnAdd:               d.StartOnAdd,
		UseSsl:                   d.UseSsl,
		Enable:                   d.Enable,
		RemoveFailedDownloads:    d.RemoveFailedDownloads,
		RemoveCompletedDownloads: d.RemoveCompletedDownloads,
	}
}

func (d *DownloadClientFlood) fromDownloadClient(client *DownloadClient) {
	d.Tags = client.Tags
	d.FieldTags = client.FieldTags
	d.AdditionalTags = client.AdditionalTags
	d.PostImportTags = client.PostImportTags
	d.Name = client.Name
	d.Host = client.Host
	d.URLBase = client.URLBase
	d.Username = client.Username
	d.Password = client.Password
	d.Destination = client.Destination
	d.TvCategory = client.TvCategory
	d.TvDirectory = client.TvDirectory
	d.RecentTvPriority = client.RecentTvPriority
	d.OlderTvPriority = client.OlderTvPriority
	d.Priority = client.Priority
	d.Port = client.Port
	d.ID = client.ID
	d.StartOnAdd = client.StartOnAdd
	d.UseSsl = client.UseSsl
	d.Enable = client.Enable
	d.RemoveFailedDownloads = client.RemoveFailedDownloads
	d.RemoveCompletedDownloads = client.RemoveCompletedDownloads
}

func (r *DownloadClientFloodResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_" + downloadClientFloodResourceName
}

func (r *DownloadClientFloodResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "<!-- subcategory:Download Clients -->Download Client Flood resource.\nFor more information refer to [Download Client](https://wiki.servarr.com/sonarr/settings#download-clients) and [Flood](https://wiki.servarr.com/sonarr/supported#flood).",
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
			"start_on_add": schema.BoolAttribute{
				MarkdownDescription: "Start on add flag.",
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
			"destination": schema.StringAttribute{
				MarkdownDescription: "Destination.",
				Optional:            true,
				Computed:            true,
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
			"field_tags": schema.SetAttribute{
				MarkdownDescription: "Field tags.",
				Optional:            true,
				Computed:            true,
				ElementType:         types.StringType,
			},
			"additional_tags": schema.SetAttribute{
				MarkdownDescription: "Additional tags, `0` TitleSlug, `1` Quality, `2` Language, `3` ReleaseGroup, `4` Year, `5` Indexer, `6` Network.",
				Optional:            true,
				Computed:            true,
				ElementType:         types.Int64Type,
			},
			"post_import_tags": schema.SetAttribute{
				MarkdownDescription: "Post import tags.",
				Optional:            true,
				Computed:            true,
				ElementType:         types.StringType,
			},
		},
	}
}

func (r *DownloadClientFloodResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *DownloadClientFloodResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	// Retrieve values from plan
	var client *DownloadClientFlood

	resp.Diagnostics.Append(req.Plan.Get(ctx, &client)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Create new DownloadClientFlood
	request := client.read(ctx)

	response, err := r.client.AddDownloadClientContext(ctx, request)
	if err != nil {
		resp.Diagnostics.AddError(tools.ClientError, fmt.Sprintf("Unable to create %s, got error: %s", downloadClientFloodResourceName, err))

		return
	}

	tflog.Trace(ctx, "created "+downloadClientFloodResourceName+": "+strconv.Itoa(int(response.ID)))
	// Generate resource state struct
	client.write(ctx, response)
	// panic(fmt.Sprint(response.Fields[7], response.Fields[8]))
	resp.Diagnostics.Append(resp.State.Set(ctx, &client)...)
}

func (r *DownloadClientFloodResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	// Get current state
	var client DownloadClientFlood

	resp.Diagnostics.Append(req.State.Get(ctx, &client)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Get DownloadClientFlood current value
	response, err := r.client.GetDownloadClientContext(ctx, client.ID.ValueInt64())
	if err != nil {
		resp.Diagnostics.AddError(tools.ClientError, fmt.Sprintf("Unable to read %s, got error: %s", downloadClientFloodResourceName, err))

		return
	}

	tflog.Trace(ctx, "read "+downloadClientFloodResourceName+": "+strconv.Itoa(int(response.ID)))
	// Map response body to resource schema attribute
	client.write(ctx, response)
	resp.Diagnostics.Append(resp.State.Set(ctx, &client)...)
}

func (r *DownloadClientFloodResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	// Get plan values
	var client *DownloadClientFlood

	resp.Diagnostics.Append(req.Plan.Get(ctx, &client)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Update DownloadClientFlood
	request := client.read(ctx)

	response, err := r.client.UpdateDownloadClientContext(ctx, request)
	if err != nil {
		resp.Diagnostics.AddError(tools.ClientError, fmt.Sprintf("Unable to update %s, got error: %s", downloadClientFloodResourceName, err))

		return
	}

	tflog.Trace(ctx, "updated "+downloadClientFloodResourceName+": "+strconv.Itoa(int(response.ID)))
	// Generate resource state struct
	client.write(ctx, response)
	resp.Diagnostics.Append(resp.State.Set(ctx, &client)...)
}

func (r *DownloadClientFloodResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var client *DownloadClientFlood

	resp.Diagnostics.Append(req.State.Get(ctx, &client)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Delete DownloadClientFlood current value
	err := r.client.DeleteDownloadClientContext(ctx, client.ID.ValueInt64())
	if err != nil {
		resp.Diagnostics.AddError(tools.ClientError, fmt.Sprintf("Unable to read %s, got error: %s", downloadClientFloodResourceName, err))

		return
	}

	tflog.Trace(ctx, "deleted "+downloadClientFloodResourceName+": "+strconv.Itoa(int(client.ID.ValueInt64())))
	resp.State.RemoveResource(ctx)
}

func (r *DownloadClientFloodResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	// resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
	id, err := strconv.Atoi(req.ID)
	if err != nil {
		resp.Diagnostics.AddError(
			tools.UnexpectedImportIdentifier,
			fmt.Sprintf("Expected import identifier with format: ID. Got: %q", req.ID),
		)

		return
	}

	tflog.Trace(ctx, "imported "+downloadClientFloodResourceName+": "+strconv.Itoa(id))
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("id"), id)...)
}

func (d *DownloadClientFlood) write(ctx context.Context, downloadClient *sonarr.DownloadClientOutput) {
	genericDownloadClient := DownloadClient{
		Enable:                   types.BoolValue(downloadClient.Enable),
		RemoveCompletedDownloads: types.BoolValue(downloadClient.RemoveCompletedDownloads),
		RemoveFailedDownloads:    types.BoolValue(downloadClient.RemoveFailedDownloads),
		Priority:                 types.Int64Value(int64(downloadClient.Priority)),
		ID:                       types.Int64Value(downloadClient.ID),
		Name:                     types.StringValue(downloadClient.Name),
		AdditionalTags:           types.SetValueMust(types.Int64Type, nil),
		FieldTags:                types.SetValueMust(types.StringType, nil),
		PostImportTags:           types.SetValueMust(types.StringType, nil),
	}
	genericDownloadClient.Tags, _ = types.SetValueFrom(ctx, types.Int64Type, downloadClient.Tags)
	genericDownloadClient.writeFields(ctx, downloadClient.Fields)
	d.fromDownloadClient(&genericDownloadClient)
}

func (d *DownloadClientFlood) read(ctx context.Context) *sonarr.DownloadClientInput {
	var tags []int

	tfsdk.ValueAs(ctx, d.Tags, &tags)

	return &sonarr.DownloadClientInput{
		Enable:                   d.Enable.ValueBool(),
		RemoveCompletedDownloads: d.RemoveCompletedDownloads.ValueBool(),
		RemoveFailedDownloads:    d.RemoveFailedDownloads.ValueBool(),
		Priority:                 int(d.Priority.ValueInt64()),
		ID:                       d.ID.ValueInt64(),
		ConfigContract:           downloadClientFloodConfigContrat,
		Implementation:           downloadClientFloodImplementation,
		Name:                     d.Name.ValueString(),
		Protocol:                 DownloadClientFloodProtocol,
		Tags:                     tags,
		Fields:                   d.toDownloadClient().readFields(ctx),
	}
}

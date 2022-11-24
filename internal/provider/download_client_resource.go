package provider

import (
	"context"
	"fmt"
	"strconv"

	"github.com/devopsarr/terraform-provider-sonarr/tools"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"golang.org/x/exp/slices"
	"golift.io/starr"
	"golift.io/starr/sonarr"
)

const downloadClientResourceName = "download_client"

// Ensure provider defined types fully satisfy framework interfaces.
var _ resource.Resource = &DownloadClientResource{}
var _ resource.ResourceWithImportState = &DownloadClientResource{}

var (
	downloadClientBoolFields        = []string{"addPaused", "useSsl", "startOnAdd", "sequentialOrder", "firstAndLast", "addStopped", "saveMagnetFiles", "readOnly"}
	downloadClientIntFields         = []string{"port", "recentTvPriority", "olderTvPriority", "initialState", "intialState"}
	downloadClientStringFields      = []string{"host", "apiKey", "urlBase", "rpcPath", "secretToken", "password", "username", "tvCategory", "tvImportedCategory", "tvDirectory", "destination", "category", "nzbFolder", "strmFolder", "torrentFolder", "magnetFileExtension", "watchFolder"}
	downloadClientStringSliceFields = []string{"fieldTags", "postImTags"}
	downloadClientIntSliceFields    = []string{"additionalTags"}
)

func NewDownloadClientResource() resource.Resource {
	return &DownloadClientResource{}
}

// DownloadClientResource defines the download client implementation.
type DownloadClientResource struct {
	client *sonarr.Sonarr
}

// DownloadClient describes the download client data model.
type DownloadClient struct {
	Tags                     types.Set    `tfsdk:"tags"`
	PostImTags               types.Set    `tfsdk:"post_im_tags"`
	FieldTags                types.Set    `tfsdk:"field_tags"`
	AdditionalTags           types.Set    `tfsdk:"additional_tags"`
	NzbFolder                types.String `tfsdk:"nzb_folder"`
	Category                 types.String `tfsdk:"category"`
	Implementation           types.String `tfsdk:"implementation"`
	Name                     types.String `tfsdk:"name"`
	Protocol                 types.String `tfsdk:"protocol"`
	MagnetFileExtension      types.String `tfsdk:"magnet_file_extension"`
	TorrentFolder            types.String `tfsdk:"torrent_folder"`
	WatchFolder              types.String `tfsdk:"watch_folder"`
	StrmFolder               types.String `tfsdk:"strm_folder"`
	Host                     types.String `tfsdk:"host"`
	ConfigContract           types.String `tfsdk:"config_contract"`
	Destination              types.String `tfsdk:"destination"`
	TvDirectory              types.String `tfsdk:"tv_directory"`
	Username                 types.String `tfsdk:"username"`
	TvImportedCategory       types.String `tfsdk:"tv_imported_category"`
	TvCategory               types.String `tfsdk:"tv_category"`
	Password                 types.String `tfsdk:"password"`
	SecretToken              types.String `tfsdk:"secret_token"`
	RPCPath                  types.String `tfsdk:"rpc_path"`
	URLBase                  types.String `tfsdk:"url_base"`
	APIKey                   types.String `tfsdk:"api_key"`
	RecentTvPriority         types.Int64  `tfsdk:"recent_tv_priority"`
	IntialState              types.Int64  `tfsdk:"intial_state"`
	InitialState             types.Int64  `tfsdk:"initial_state"`
	OlderTvPriority          types.Int64  `tfsdk:"older_tv_priority"`
	Priority                 types.Int64  `tfsdk:"priority"`
	Port                     types.Int64  `tfsdk:"port"`
	ID                       types.Int64  `tfsdk:"id"`
	AddStopped               types.Bool   `tfsdk:"add_stopped"`
	SaveMagnetFiles          types.Bool   `tfsdk:"save_magnet_files"`
	ReadOnly                 types.Bool   `tfsdk:"read_only"`
	FirstAndLast             types.Bool   `tfsdk:"first_and_last"`
	SequentialOrder          types.Bool   `tfsdk:"sequential_order"`
	StartOnAdd               types.Bool   `tfsdk:"start_on_add"`
	UseSsl                   types.Bool   `tfsdk:"use_ssl"`
	AddPaused                types.Bool   `tfsdk:"add_paused"`
	Enable                   types.Bool   `tfsdk:"enable"`
	RemoveFailedDownloads    types.Bool   `tfsdk:"remove_failed_downloads"`
	RemoveCompletedDownloads types.Bool   `tfsdk:"remove_completed_downloads"`
}

func (r *DownloadClientResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_" + downloadClientResourceName
}

func (r *DownloadClientResource) GetSchema(ctx context.Context) (tfsdk.Schema, diag.Diagnostics) {
	return tfsdk.Schema{
		MarkdownDescription: "<!-- subcategory:Download Clients -->Download Client resource.\nFor more information refer to [Download Client](https://wiki.servarr.com/sonarr/settings#download-clients).",
		Attributes: map[string]tfsdk.Attribute{
			"enable": {
				MarkdownDescription: "Enable flag.",
				Optional:            true,
				Computed:            true,
				Type:                types.BoolType,
			},
			"remove_completed_downloads": {
				MarkdownDescription: "Remove completed downloads flag.",
				Optional:            true,
				Computed:            true,
				Type:                types.BoolType,
			},
			"remove_failed_downloads": {
				MarkdownDescription: "Remove failed downloads flag.",
				Optional:            true,
				Computed:            true,
				Type:                types.BoolType,
			},
			"priority": {
				MarkdownDescription: "Priority.",
				Optional:            true,
				Computed:            true,
				Type:                types.Int64Type,
			},
			"config_contract": {
				MarkdownDescription: "DownloadClient configuration template.",
				Required:            true,
				Type:                types.StringType,
			},
			"implementation": {
				MarkdownDescription: "DownloadClient implementation name.",
				Required:            true,
				Type:                types.StringType,
			},
			"name": {
				MarkdownDescription: "Download Client name.",
				Required:            true,
				Type:                types.StringType,
			},
			"protocol": {
				MarkdownDescription: "Protocol. Valid values are 'usenet' and 'torrent'.",
				Required:            true,
				Type:                types.StringType,
				Validators: []tfsdk.AttributeValidator{
					tools.StringMatch([]string{"usenet", "torrent"}),
				},
			},
			"tags": {
				MarkdownDescription: "List of associated tags.",
				Optional:            true,
				Computed:            true,
				Type: types.SetType{
					ElemType: types.Int64Type,
				},
			},
			"id": {
				MarkdownDescription: "Download Client ID.",
				Computed:            true,
				Type:                types.Int64Type,
				PlanModifiers: tfsdk.AttributePlanModifiers{
					resource.UseStateForUnknown(),
				},
			},
			// Field values
			"add_paused": {
				MarkdownDescription: "Add paused flag.",
				Optional:            true,
				Computed:            true,
				Type:                types.BoolType,
			},
			"use_ssl": {
				MarkdownDescription: "Use SSL flag.",
				Optional:            true,
				Computed:            true,
				Type:                types.BoolType,
			},
			"start_on_add": {
				MarkdownDescription: "Start on add flag.",
				Optional:            true,
				Computed:            true,
				Type:                types.BoolType,
			},
			"sequential_order": {
				MarkdownDescription: "Sequential order flag.",
				Optional:            true,
				Computed:            true,
				Type:                types.BoolType,
			},
			"first_and_last": {
				MarkdownDescription: "First and last flag.",
				Optional:            true,
				Computed:            true,
				Type:                types.BoolType,
			},
			"add_stopped": {
				MarkdownDescription: "Add stopped flag.",
				Optional:            true,
				Computed:            true,
				Type:                types.BoolType,
			},
			"save_magnet_files": {
				MarkdownDescription: "Save magnet files flag.",
				Optional:            true,
				Computed:            true,
				Type:                types.BoolType,
			},
			"read_only": {
				MarkdownDescription: "Read only flag.",
				Optional:            true,
				Computed:            true,
				Type:                types.BoolType,
			},
			"port": {
				MarkdownDescription: "Port.",
				Optional:            true,
				Computed:            true,
				Type:                types.Int64Type,
			},
			"recent_tv_priority": {
				MarkdownDescription: "Recent TV priority. `0` Last, `1` First.",
				Optional:            true,
				Computed:            true,
				Type:                types.Int64Type,
				Validators: []tfsdk.AttributeValidator{
					tools.IntMatch([]int64{0, 1}),
				},
			},
			"older_tv_priority": {
				MarkdownDescription: "Older TV priority. `0` Last, `1` First.",
				Optional:            true,
				Computed:            true,
				Type:                types.Int64Type,
				Validators: []tfsdk.AttributeValidator{
					tools.IntMatch([]int64{0, 1}),
				},
			},
			"initial_state": {
				MarkdownDescription: "Initial state. `0` Start, `1` ForceStart, `2` Pause.",
				Optional:            true,
				Computed:            true,
				Type:                types.Int64Type,
				Validators: []tfsdk.AttributeValidator{
					tools.IntMatch([]int64{0, 1}),
				},
			},
			"intial_state": {
				MarkdownDescription: "Initial state, with Stop support. `0` Start, `1` ForceStart, `2` Pause, `3` Stop.",
				Optional:            true,
				Computed:            true,
				Type:                types.Int64Type,
			},
			"host": {
				MarkdownDescription: "host.",
				Optional:            true,
				Computed:            true,
				Type:                types.StringType,
			},
			"api_key": {
				MarkdownDescription: "API key.",
				Optional:            true,
				Computed:            true,
				Type:                types.StringType,
			},
			"rpc_path": {
				MarkdownDescription: "RPC path.",
				Optional:            true,
				Computed:            true,
				Type:                types.StringType,
			},
			"url_base": {
				MarkdownDescription: "Base URL.",
				Optional:            true,
				Computed:            true,
				Type:                types.StringType,
			},
			"secret_token": {
				MarkdownDescription: "Secret token.",
				Optional:            true,
				Computed:            true,
				Type:                types.StringType,
			},
			"username": {
				MarkdownDescription: "Username.",
				Optional:            true,
				Computed:            true,
				Type:                types.StringType,
			},
			"password": {
				MarkdownDescription: "Password.",
				Optional:            true,
				Computed:            true,
				Type:                types.StringType,
			},
			"tv_category": {
				MarkdownDescription: "TV category.",
				Optional:            true,
				Computed:            true,
				Type:                types.StringType,
			},
			"tv_imported_category": {
				MarkdownDescription: "TV imported category.",
				Optional:            true,
				Computed:            true,
				Type:                types.StringType,
			},
			"tv_directory": {
				MarkdownDescription: "TV directory.",
				Optional:            true,
				Computed:            true,
				Type:                types.StringType,
			},
			"destination": {
				MarkdownDescription: "Destination.",
				Optional:            true,
				Computed:            true,
				Type:                types.StringType,
			},
			"category": {
				MarkdownDescription: "Category.",
				Optional:            true,
				Computed:            true,
				Type:                types.StringType,
			},
			"nzb_folder": {
				MarkdownDescription: "NZB folder.",
				Optional:            true,
				Computed:            true,
				Type:                types.StringType,
			},
			"strm_folder": {
				MarkdownDescription: "STRM folder.",
				Optional:            true,
				Computed:            true,
				Type:                types.StringType,
			},
			"torrent_folder": {
				MarkdownDescription: "Torrent folder.",
				Optional:            true,
				Computed:            true,
				Type:                types.StringType,
			},
			"watch_folder": {
				MarkdownDescription: "Watch folder flag.",
				Optional:            true,
				Computed:            true,
				Type:                types.StringType,
			},
			"magnet_file_extension": {
				MarkdownDescription: "Magnet file extension.",
				Optional:            true,
				Computed:            true,
				Type:                types.StringType,
			},
			"additional_tags": {
				MarkdownDescription: "Additional tags, `0` TitleSlug, `1` Quality, `2` Language, `3` ReleaseGroup, `4` Year, `5` Indexer, `6` Network.",
				Optional:            true,
				Computed:            true,
				Type: types.SetType{
					ElemType: types.Int64Type,
				},
			},
			"field_tags": {
				MarkdownDescription: "Field tags.",
				Optional:            true,
				Computed:            true,
				Type: types.SetType{
					ElemType: types.StringType,
				},
			},
			"post_im_tags": {
				MarkdownDescription: "Post import tags.",
				Optional:            true,
				Computed:            true,
				Type: types.SetType{
					ElemType: types.StringType,
				},
			},
		},
	}, nil
}

func (r *DownloadClientResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *DownloadClientResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	// Retrieve values from plan
	var client *DownloadClient

	resp.Diagnostics.Append(req.Plan.Get(ctx, &client)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Create new DownloadClient
	request := client.read(ctx)

	response, err := r.client.AddDownloadClientContext(ctx, request)
	if err != nil {
		resp.Diagnostics.AddError(tools.ClientError, fmt.Sprintf("Unable to create %s, got error: %s", downloadClientResourceName, err))

		return
	}

	tflog.Trace(ctx, "created "+downloadClientResourceName+": "+strconv.Itoa(int(response.ID)))
	// Generate resource state struct
	// this is needed because of many empty fields are unknown in both plan and read
	var state DownloadClient

	state.write(ctx, response)
	resp.Diagnostics.Append(resp.State.Set(ctx, state)...)
}

func (r *DownloadClientResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	// Get current state
	var client DownloadClient

	resp.Diagnostics.Append(req.State.Get(ctx, &client)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Get DownloadClient current value
	response, err := r.client.GetDownloadClientContext(ctx, client.ID.ValueInt64())
	if err != nil {
		resp.Diagnostics.AddError(tools.ClientError, fmt.Sprintf("Unable to read %s, got error: %s", downloadClientResourceName, err))

		return
	}

	tflog.Trace(ctx, "read "+downloadClientResourceName+": "+strconv.Itoa(int(response.ID)))
	// Map response body to resource schema attribute
	// this is needed because of many empty fields are unknown in both plan and read
	var state DownloadClient

	state.write(ctx, response)
	resp.Diagnostics.Append(resp.State.Set(ctx, state)...)
}

func (r *DownloadClientResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	// Get plan values
	var client *DownloadClient

	resp.Diagnostics.Append(req.Plan.Get(ctx, &client)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Update DownloadClient
	request := client.read(ctx)

	response, err := r.client.UpdateDownloadClientContext(ctx, request)
	if err != nil {
		resp.Diagnostics.AddError(tools.ClientError, fmt.Sprintf("Unable to update %s, got error: %s", downloadClientResourceName, err))

		return
	}

	tflog.Trace(ctx, "updated "+downloadClientResourceName+": "+strconv.Itoa(int(response.ID)))
	// Generate resource state struct
	// this is needed because of many empty fields are unknown in both plan and read
	var state DownloadClient

	state.write(ctx, response)
	resp.Diagnostics.Append(resp.State.Set(ctx, state)...)
}

func (r *DownloadClientResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var client *DownloadClient

	resp.Diagnostics.Append(req.State.Get(ctx, &client)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Delete DownloadClient current value
	err := r.client.DeleteDownloadClientContext(ctx, client.ID.ValueInt64())
	if err != nil {
		resp.Diagnostics.AddError(tools.ClientError, fmt.Sprintf("Unable to read %s, got error: %s", downloadClientResourceName, err))

		return
	}

	tflog.Trace(ctx, "deleted "+downloadClientResourceName+": "+strconv.Itoa(int(client.ID.ValueInt64())))
	resp.State.RemoveResource(ctx)
}

func (r *DownloadClientResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	// resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
	id, err := strconv.Atoi(req.ID)
	if err != nil {
		resp.Diagnostics.AddError(
			tools.UnexpectedImportIdentifier,
			fmt.Sprintf("Expected import identifier with format: ID. Got: %q", req.ID),
		)

		return
	}

	tflog.Trace(ctx, "imported "+downloadClientResourceName+": "+strconv.Itoa(id))
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("id"), id)...)
}

func (d *DownloadClient) write(ctx context.Context, downloadClient *sonarr.DownloadClientOutput) {
	d.Enable = types.BoolValue(downloadClient.Enable)
	d.RemoveCompletedDownloads = types.BoolValue(downloadClient.RemoveCompletedDownloads)
	d.RemoveFailedDownloads = types.BoolValue(downloadClient.RemoveFailedDownloads)
	d.Priority = types.Int64Value(int64(downloadClient.Priority))
	d.ID = types.Int64Value(downloadClient.ID)
	d.ConfigContract = types.StringValue(downloadClient.ConfigContract)
	d.Implementation = types.StringValue(downloadClient.Implementation)
	d.Name = types.StringValue(downloadClient.Name)
	d.Protocol = types.StringValue(downloadClient.Protocol)
	d.Tags = types.SetValueMust(types.Int64Type, nil)
	d.AdditionalTags = types.SetValueMust(types.Int64Type, nil)
	d.FieldTags = types.SetValueMust(types.StringType, nil)
	d.PostImTags = types.SetValueMust(types.StringType, nil)
	tfsdk.ValueFrom(ctx, downloadClient.Tags, d.Tags.Type(ctx), &d.Tags)
	d.writeFields(ctx, downloadClient.Fields)
}

func (d *DownloadClient) writeFields(ctx context.Context, fields []*starr.FieldOutput) {
	for _, f := range fields {
		if f.Value == nil {
			continue
		}

		if slices.Contains(downloadClientStringFields, f.Name) {
			tools.WriteStringField(f, d)

			continue
		}

		if slices.Contains(downloadClientBoolFields, f.Name) {
			tools.WriteBoolField(f, d)

			continue
		}

		if slices.Contains(downloadClientIntFields, f.Name) {
			tools.WriteIntField(f, d)

			continue
		}

		if slices.Contains(downloadClientIntSliceFields, f.Name) {
			tools.WriteIntSliceField(ctx, f, d)

			continue
		}

		if slices.Contains(downloadClientStringSliceFields, f.Name) {
			tools.WriteStringSliceField(ctx, f, d)
		}
	}
}

func (d *DownloadClient) read(ctx context.Context) *sonarr.DownloadClientInput {
	var tags []int

	tfsdk.ValueAs(ctx, d.Tags, &tags)

	return &sonarr.DownloadClientInput{
		Enable:                   d.Enable.ValueBool(),
		RemoveCompletedDownloads: d.RemoveCompletedDownloads.ValueBool(),
		RemoveFailedDownloads:    d.RemoveFailedDownloads.ValueBool(),
		Priority:                 int(d.Priority.ValueInt64()),
		ID:                       d.ID.ValueInt64(),
		ConfigContract:           d.ConfigContract.ValueString(),
		Implementation:           d.Implementation.ValueString(),
		Name:                     d.Name.ValueString(),
		Protocol:                 d.Protocol.ValueString(),
		Tags:                     tags,
		Fields:                   d.readFields(ctx),
	}
}

func (d *DownloadClient) readFields(ctx context.Context) []*starr.FieldInput {
	var output []*starr.FieldInput

	for _, b := range downloadClientBoolFields {
		if field := tools.ReadBoolField(b, d); field != nil {
			output = append(output, field)
		}
	}

	for _, i := range downloadClientIntFields {
		if field := tools.ReadIntField(i, d); field != nil {
			output = append(output, field)
		}
	}

	for _, s := range downloadClientStringFields {
		if field := tools.ReadStringField(s, d); field != nil {
			output = append(output, field)
		}
	}

	for _, s := range downloadClientStringSliceFields {
		if field := tools.ReadStringSliceField(ctx, s, d); field != nil {
			output = append(output, field)
		}
	}

	for _, s := range downloadClientIntSliceFields {
		if field := tools.ReadIntSliceField(ctx, s, d); field != nil {
			output = append(output, field)
		}
	}

	return output
}

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
	"golift.io/starr"
	"golift.io/starr/sonarr"
)

// Ensure provider defined types fully satisfy framework interfaces.
var _ resource.Resource = &DownloadClientResource{}
var _ resource.ResourceWithImportState = &DownloadClientResource{}

func NewDownloadClientResource() resource.Resource {
	return &DownloadClientResource{}
}

// DownloadClientResource defines the download client implementation.
type DownloadClientResource struct {
	client *sonarr.Sonarr
}

// DownloadClient describes the download client data model.
type DownloadClient struct {
	Enable                   types.Bool   `tfsdk:"enable"`
	RemoveCompletedDownloads types.Bool   `tfsdk:"remove_completed_downloads"`
	RemoveFailedDownloads    types.Bool   `tfsdk:"remove_failed_downloads"`
	Priority                 types.Int64  `tfsdk:"priority"`
	ID                       types.Int64  `tfsdk:"id"`
	ConfigContract           types.String `tfsdk:"config_contract"`
	Implementation           types.String `tfsdk:"implementation"`
	Name                     types.String `tfsdk:"name"`
	Protocol                 types.String `tfsdk:"protocol"`
	Tags                     types.Set    `tfsdk:"tags"`
	// Fields values
	AddPaused           types.Bool   `tfsdk:"add_paused"`
	UseSsl              types.Bool   `tfsdk:"use_ssl"`
	StartOnAdd          types.Bool   `tfsdk:"start_on_add"`
	SequentialOrder     types.Bool   `tfsdk:"sequential_order"`
	FirstAndLast        types.Bool   `tfsdk:"first_and_last"`
	AddStopped          types.Bool   `tfsdk:"add_stopped"`
	SaveMagnetFiles     types.Bool   `tfsdk:"save_magnet_files"`
	ReadOnly            types.Bool   `tfsdk:"read_only"`
	WatchFolder         types.Bool   `tfsdk:"watch_folder"`
	Port                types.Int64  `tfsdk:"port"`
	RecentTvPriority    types.Int64  `tfsdk:"recent_tv_priority"` // from 0 to 1 "Last, First"
	OlderTvPriority     types.Int64  `tfsdk:"older_tv_priority"`  // from 0 to 1 "Last, First"
	InitialState        types.Int64  `tfsdk:"initial_state"`      // from 0 to 2 "Start, ForceStart, Pause"
	IntialState         types.Int64  `tfsdk:"intial_state"`       // from 0 to 3 "Start, ForceStart, Pause, Stop"
	Host                types.String `tfsdk:"host"`
	APIKey              types.String `tfsdk:"api_key"`
	URLBase             types.String `tfsdk:"url_base"`
	RPCPath             types.String `tfsdk:"rpc_path"`
	SecretToken         types.String `tfsdk:"secret_token"`
	Password            types.String `tfsdk:"password"`
	TvCategory          types.String `tfsdk:"tv_category"`
	TvImportedCategory  types.String `tfsdk:"tv_imported_category"`
	Username            types.String `tfsdk:"username"`
	TvDirectory         types.String `tfsdk:"tv_directory"`
	Destination         types.String `tfsdk:"destination"`
	Category            types.String `tfsdk:"category"`
	NzbFolder           types.String `tfsdk:"nzb_folder"`
	StrmFolder          types.String `tfsdk:"strm_folder"`
	TorrentFolder       types.String `tfsdk:"torrent_folder"`
	MagnetFileExtension types.String `tfsdk:"magnet_file_extension"`
	AdditionalTags      types.Set    `tfsdk:"additional_tags"` // int
	FieldTags           types.Set    `tfsdk:"field_tags"`      // strings
	PostImTags          types.Set    `tfsdk:"post_im_tags"`    // strings
}

func (r *DownloadClientResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_download_client"
}

func (r *DownloadClientResource) GetSchema(ctx context.Context) (tfsdk.Schema, diag.Diagnostics) {
	return tfsdk.Schema{
		MarkdownDescription: "Download Client resource.<br/>For more information refer to [Download Client](https://wiki.servarr.com/sonarr/settings#download-clients).",
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
					helpers.StringMatch([]string{"usenet", "torrent"}),
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
			"watch_folder": {
				MarkdownDescription: "Watch folder flag.",
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
					helpers.IntMatch([]int64{0, 1}),
				},
			},
			"older_tv_priority": {
				MarkdownDescription: "Older TV priority. `0` Last, `1` First.",
				Optional:            true,
				Computed:            true,
				Type:                types.Int64Type,
				Validators: []tfsdk.AttributeValidator{
					helpers.IntMatch([]int64{0, 1}),
				},
			},
			"initial_state": {
				MarkdownDescription: "Initial state. `0` Start, `1` ForceStart, `2` Pause.",
				Optional:            true,
				Computed:            true,
				Type:                types.Int64Type,
				Validators: []tfsdk.AttributeValidator{
					helpers.IntMatch([]int64{0, 1}),
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
			UnexpectedResourceConfigureType,
			fmt.Sprintf("Expected *sonarr.Sonarr, got: %T. Please report this issue to the provider developers.", req.ProviderData),
		)

		return
	}

	r.client = client
}

func (r *DownloadClientResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	// Retrieve values from plan
	var plan DownloadClient

	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Create new DownloadClient
	request := readDownloadClient(ctx, &plan)

	response, err := r.client.AddDownloadClientContext(ctx, request)
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to create DownloadClient, got error: %s", err))

		return
	}

	tflog.Trace(ctx, "created download_client: "+strconv.Itoa(int(response.ID)))
	// Generate resource state struct
	result := writeDownloadClient(ctx, response)
	resp.Diagnostics.Append(resp.State.Set(ctx, &result)...)
}

func (r *DownloadClientResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	// Get current state
	var state DownloadClient

	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Get DownloadClient current value
	response, err := r.client.GetDownloadClientContext(ctx, int(state.ID.Value))
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to read DownloadClients, got error: %s", err))

		return
	}

	tflog.Trace(ctx, "read download_client: "+strconv.Itoa(int(response.ID)))
	// Map response body to resource schema attribute
	result := writeDownloadClient(ctx, response)
	resp.Diagnostics.Append(resp.State.Set(ctx, &result)...)
}

func (r *DownloadClientResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	// Get plan values
	var plan DownloadClient
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Update DownloadClient
	request := readDownloadClient(ctx, &plan)

	response, err := r.client.UpdateDownloadClientContext(ctx, request)
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to update DownloadClient, got error: %s", err))

		return
	}

	tflog.Trace(ctx, "updated download_client: "+strconv.Itoa(int(response.ID)))
	// Generate resource state struct
	result := writeDownloadClient(ctx, response)
	resp.Diagnostics.Append(resp.State.Set(ctx, &result)...)
}

func (r *DownloadClientResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state DownloadClient

	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Delete DownloadClient current value
	err := r.client.DeleteDownloadClientContext(ctx, int(state.ID.Value))
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to read DownloadClients, got error: %s", err))

		return
	}

	tflog.Trace(ctx, "deleted download_client: "+strconv.Itoa(int(state.ID.Value)))
	resp.State.RemoveResource(ctx)
}

func (r *DownloadClientResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	// resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
	id, err := strconv.Atoi(req.ID)
	if err != nil {
		resp.Diagnostics.AddError(
			"Unexpected Import Identifier",
			fmt.Sprintf("Expected import identifier with format: ID. Got: %q", req.ID),
		)

		return
	}

	tflog.Trace(ctx, "imported download_client: "+strconv.Itoa(id))
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("id"), id)...)
}

func writeDownloadClient(ctx context.Context, downloadClient *sonarr.DownloadClientOutput) *DownloadClient {
	output := DownloadClient{
		Enable:                   types.Bool{Value: downloadClient.Enable},
		RemoveCompletedDownloads: types.Bool{Value: downloadClient.RemoveCompletedDownloads},
		RemoveFailedDownloads:    types.Bool{Value: downloadClient.RemoveFailedDownloads},
		Priority:                 types.Int64{Value: int64(downloadClient.Priority)},
		ID:                       types.Int64{Value: downloadClient.ID},
		ConfigContract:           types.String{Value: downloadClient.ConfigContract},
		Implementation:           types.String{Value: downloadClient.Implementation},
		Name:                     types.String{Value: downloadClient.Name},
		Protocol:                 types.String{Value: downloadClient.Protocol},
		Tags:                     types.Set{ElemType: types.Int64Type},
		AdditionalTags:           types.Set{ElemType: types.Int64Type},
		FieldTags:                types.Set{ElemType: types.StringType},
		PostImTags:               types.Set{ElemType: types.StringType},
	}
	tfsdk.ValueFrom(ctx, downloadClient.Tags, output.Tags.Type(ctx), &output.Tags)

	for _, f := range downloadClient.Fields {
		if f.Value != nil {
			switch f.Name {
			case "addPaused":
				output.AddPaused = types.Bool{Value: f.Value.(bool)}
			case "useSsl":
				output.UseSsl = types.Bool{Value: f.Value.(bool)}
			case "startOnAdd":
				output.StartOnAdd = types.Bool{Value: f.Value.(bool)}
			case "sequentialOrder":
				output.SequentialOrder = types.Bool{Value: f.Value.(bool)}
			case "firstAndLast":
				output.FirstAndLast = types.Bool{Value: f.Value.(bool)}
			case "addStopped":
				output.AddStopped = types.Bool{Value: f.Value.(bool)}
			case "saveMagnetFiles":
				output.SaveMagnetFiles = types.Bool{Value: f.Value.(bool)}
			case "readOnly":
				output.ReadOnly = types.Bool{Value: f.Value.(bool)}
			case "watchFolder":
				output.WatchFolder = types.Bool{Value: f.Value.(bool)}
			case "port":
				output.Port = types.Int64{Value: int64(f.Value.(float64))}
			case "recentTvPriority":
				output.RecentTvPriority = types.Int64{Value: int64(f.Value.(float64))}
			case "olderTvPriority":
				output.OlderTvPriority = types.Int64{Value: int64(f.Value.(float64))}
			case "initialState":
				output.InitialState = types.Int64{Value: int64(f.Value.(float64))}
			case "intialState":
				output.IntialState = types.Int64{Value: int64(f.Value.(float64))}
			case "host":
				output.Host = types.String{Value: f.Value.(string)}
			case "apiKey":
				output.APIKey = types.String{Value: f.Value.(string)}
			case "urlBase":
				output.URLBase = types.String{Value: f.Value.(string)}
			case "rpcPath":
				output.RPCPath = types.String{Value: f.Value.(string)}
			case "secretToken":
				output.SecretToken = types.String{Value: f.Value.(string)}
			case "password":
				output.Password = types.String{Value: f.Value.(string)}
			case "username":
				output.Username = types.String{Value: f.Value.(string)}
			case "tvCategory":
				output.TvCategory = types.String{Value: f.Value.(string)}
			case "tvImportedCategory":
				output.TvImportedCategory = types.String{Value: f.Value.(string)}
			case "tvDirectory":
				output.TvDirectory = types.String{Value: f.Value.(string)}
			case "destination":
				output.Destination = types.String{Value: f.Value.(string)}
			case "category":
				output.Category = types.String{Value: f.Value.(string)}
			case "nzbFolder":
				output.NzbFolder = types.String{Value: f.Value.(string)}
			case "strmFolder":
				output.StrmFolder = types.String{Value: f.Value.(string)}
			case "torrentFolder":
				output.TorrentFolder = types.String{Value: f.Value.(string)}
			case "magnetFileExtension":
				output.MagnetFileExtension = types.String{Value: f.Value.(string)}
			case "fieldTags":
				tfsdk.ValueFrom(ctx, f.Value, output.FieldTags.Type(ctx), &output.FieldTags)
			case "postImTags":
				output.PostImTags = types.Set{ElemType: types.StringType}
				tfsdk.ValueFrom(ctx, f.Value, output.PostImTags.Type(ctx), &output.PostImTags)
			case "additionalTags":
				tfsdk.ValueFrom(ctx, f.Value, output.AdditionalTags.Type(ctx), &output.AdditionalTags)
			// TODO: manage unknown values
			default:
			}
		}
	}

	return &output
}

func readDownloadClient(ctx context.Context, downloadClient *DownloadClient) *sonarr.DownloadClientInput {
	var tags []int

	tfsdk.ValueAs(ctx, downloadClient.Tags, &tags)

	return &sonarr.DownloadClientInput{
		Enable:                   downloadClient.Enable.Value,
		RemoveCompletedDownloads: downloadClient.RemoveCompletedDownloads.Value,
		RemoveFailedDownloads:    downloadClient.RemoveFailedDownloads.Value,
		Priority:                 int(downloadClient.Priority.Value),
		ID:                       downloadClient.ID.Value,
		ConfigContract:           downloadClient.ConfigContract.Value,
		Implementation:           downloadClient.Implementation.Value,
		Name:                     downloadClient.Name.Value,
		Protocol:                 downloadClient.Protocol.Value,
		Tags:                     tags,
		Fields:                   readDownloadClientFields(ctx, downloadClient),
	}
}

func readDownloadClientFields(ctx context.Context, downloadClient *DownloadClient) []*starr.FieldInput {
	var output []*starr.FieldInput
	if !downloadClient.AddPaused.IsNull() && !downloadClient.AddPaused.IsUnknown() {
		output = append(output, &starr.FieldInput{
			Name:  "addPaused",
			Value: downloadClient.AddPaused.Value,
		})
	}

	if !downloadClient.UseSsl.IsNull() && !downloadClient.UseSsl.IsUnknown() {
		output = append(output, &starr.FieldInput{
			Name:  "useSsl",
			Value: downloadClient.UseSsl.Value,
		})
	}

	if !downloadClient.StartOnAdd.IsNull() && !downloadClient.StartOnAdd.IsUnknown() {
		output = append(output, &starr.FieldInput{
			Name:  "startOnAdd",
			Value: downloadClient.StartOnAdd.Value,
		})
	}

	if !downloadClient.SequentialOrder.IsNull() && !downloadClient.SequentialOrder.IsUnknown() {
		output = append(output, &starr.FieldInput{
			Name:  "sequentialOrder",
			Value: downloadClient.SequentialOrder.Value,
		})
	}

	if !downloadClient.FirstAndLast.IsNull() && !downloadClient.FirstAndLast.IsUnknown() {
		output = append(output, &starr.FieldInput{
			Name:  "firstAndLast",
			Value: downloadClient.FirstAndLast.Value,
		})
	}

	if !downloadClient.AddStopped.IsNull() && !downloadClient.AddStopped.IsUnknown() {
		output = append(output, &starr.FieldInput{
			Name:  "addStopped",
			Value: downloadClient.AddStopped.Value,
		})
	}

	if !downloadClient.SaveMagnetFiles.IsNull() && !downloadClient.SaveMagnetFiles.IsUnknown() {
		output = append(output, &starr.FieldInput{
			Name:  "saveMagnetFiles",
			Value: downloadClient.SaveMagnetFiles.Value,
		})
	}

	if !downloadClient.ReadOnly.IsNull() && !downloadClient.ReadOnly.IsUnknown() {
		output = append(output, &starr.FieldInput{
			Name:  "readOnly",
			Value: downloadClient.ReadOnly.Value,
		})
	}

	if !downloadClient.WatchFolder.IsNull() && !downloadClient.WatchFolder.IsUnknown() {
		output = append(output, &starr.FieldInput{
			Name:  "watchFolder",
			Value: downloadClient.WatchFolder.Value,
		})
	}

	if !downloadClient.Port.IsNull() && !downloadClient.Port.IsUnknown() {
		output = append(output, &starr.FieldInput{
			Name:  "port",
			Value: downloadClient.Port.Value,
		})
	}

	if !downloadClient.RecentTvPriority.IsNull() && !downloadClient.RecentTvPriority.IsUnknown() {
		output = append(output, &starr.FieldInput{
			Name:  "recentTvPriority",
			Value: downloadClient.RecentTvPriority.Value,
		})
	}

	if !downloadClient.OlderTvPriority.IsNull() && !downloadClient.OlderTvPriority.IsUnknown() {
		output = append(output, &starr.FieldInput{
			Name:  "olderTvPriority",
			Value: downloadClient.OlderTvPriority.Value,
		})
	}

	if !downloadClient.InitialState.IsNull() && !downloadClient.InitialState.IsUnknown() {
		output = append(output, &starr.FieldInput{
			Name:  "initialState",
			Value: downloadClient.InitialState.Value,
		})
	}

	if !downloadClient.IntialState.IsNull() && !downloadClient.IntialState.IsUnknown() {
		output = append(output, &starr.FieldInput{
			Name:  "intialState",
			Value: downloadClient.IntialState.Value,
		})
	}

	if !downloadClient.APIKey.IsNull() && !downloadClient.APIKey.IsUnknown() {
		output = append(output, &starr.FieldInput{
			Name:  "apiKey",
			Value: downloadClient.APIKey.Value,
		})
	}

	if !downloadClient.Host.IsNull() && !downloadClient.Host.IsUnknown() {
		output = append(output, &starr.FieldInput{
			Name:  "host",
			Value: downloadClient.Host.Value,
		})
	}

	if !downloadClient.URLBase.IsNull() && !downloadClient.URLBase.IsUnknown() {
		output = append(output, &starr.FieldInput{
			Name:  "urlBase",
			Value: downloadClient.URLBase.Value,
		})
	}

	if !downloadClient.RPCPath.IsNull() && !downloadClient.RPCPath.IsUnknown() {
		output = append(output, &starr.FieldInput{
			Name:  "rpcPath",
			Value: downloadClient.RPCPath.Value,
		})
	}

	if !downloadClient.SecretToken.IsNull() && !downloadClient.SecretToken.IsUnknown() {
		output = append(output, &starr.FieldInput{
			Name:  "secretToken",
			Value: downloadClient.SecretToken.Value,
		})
	}

	if !downloadClient.TvCategory.IsNull() && !downloadClient.TvCategory.IsUnknown() {
		output = append(output, &starr.FieldInput{
			Name:  "tvCategory",
			Value: downloadClient.TvCategory.Value,
		})
	}

	if !downloadClient.TvImportedCategory.IsNull() && !downloadClient.TvImportedCategory.IsUnknown() {
		output = append(output, &starr.FieldInput{
			Name:  "tvImportedCategory",
			Value: downloadClient.TvImportedCategory.Value,
		})
	}

	if !downloadClient.TvDirectory.IsNull() && !downloadClient.TvDirectory.IsUnknown() {
		output = append(output, &starr.FieldInput{
			Name:  "tvDirectory",
			Value: downloadClient.TvDirectory.Value,
		})
	}

	if !downloadClient.Destination.IsNull() && !downloadClient.Destination.IsUnknown() {
		output = append(output, &starr.FieldInput{
			Name:  "destination",
			Value: downloadClient.Destination.Value,
		})
	}

	if !downloadClient.Password.IsNull() && !downloadClient.Password.IsUnknown() {
		output = append(output, &starr.FieldInput{
			Name:  "password",
			Value: downloadClient.Password.Value,
		})
	}

	if !downloadClient.Username.IsNull() && !downloadClient.Username.IsUnknown() {
		output = append(output, &starr.FieldInput{
			Name:  "username",
			Value: downloadClient.Username.Value,
		})
	}

	if !downloadClient.Category.IsNull() && !downloadClient.Category.IsUnknown() {
		output = append(output, &starr.FieldInput{
			Name:  "category",
			Value: downloadClient.Category.Value,
		})
	}

	if !downloadClient.NzbFolder.IsNull() && !downloadClient.NzbFolder.IsUnknown() {
		output = append(output, &starr.FieldInput{
			Name:  "nzbFolder",
			Value: downloadClient.NzbFolder.Value,
		})
	}

	if !downloadClient.StrmFolder.IsNull() && !downloadClient.StrmFolder.IsUnknown() {
		output = append(output, &starr.FieldInput{
			Name:  "strmFolder",
			Value: downloadClient.StrmFolder.Value,
		})
	}

	if !downloadClient.TorrentFolder.IsNull() && !downloadClient.TorrentFolder.IsUnknown() {
		output = append(output, &starr.FieldInput{
			Name:  "torrentFolder",
			Value: downloadClient.TorrentFolder.Value,
		})
	}

	if !downloadClient.MagnetFileExtension.IsNull() && !downloadClient.MagnetFileExtension.IsUnknown() {
		output = append(output, &starr.FieldInput{
			Name:  "magnetFileExtension",
			Value: downloadClient.MagnetFileExtension.Value,
		})
	}

	if len(downloadClient.AdditionalTags.Elems) != 0 {
		tags := make([]int64, len(downloadClient.AdditionalTags.Elems))
		tfsdk.ValueAs(ctx, downloadClient.AdditionalTags, &tags)

		output = append(output, &starr.FieldInput{
			Name:  "additionalTags",
			Value: tags,
		})
	}

	if len(downloadClient.FieldTags.Elems) != 0 {
		tags := make([]string, len(downloadClient.FieldTags.Elems))
		tfsdk.ValueAs(ctx, downloadClient.FieldTags, &tags)

		output = append(output, &starr.FieldInput{
			Name:  "fieldTags",
			Value: tags,
		})
	}

	if len(downloadClient.PostImTags.Elems) != 0 {
		tags := make([]string, len(downloadClient.PostImTags.Elems))
		tfsdk.ValueAs(ctx, downloadClient.PostImTags, &tags)

		output = append(output, &starr.FieldInput{
			Name:  "postImTags",
			Value: tags,
		})
	}

	return output
}

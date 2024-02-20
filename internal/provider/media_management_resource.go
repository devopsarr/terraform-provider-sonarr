package provider

import (
	"context"
	"strconv"

	"github.com/devopsarr/sonarr-go/sonarr"
	"github.com/devopsarr/terraform-provider-sonarr/internal/helpers"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

const mediaManagementResourceName = "media_management"

// Ensure provider defined types fully satisfy framework interfaces.
var (
	_ resource.Resource                = &MediaManagementResource{}
	_ resource.ResourceWithImportState = &MediaManagementResource{}
)

func NewMediaManagementResource() resource.Resource {
	return &MediaManagementResource{}
}

// MediaManagementResource defines the media management implementation.
type MediaManagementResource struct {
	client *sonarr.APIClient
	auth   context.Context
}

// MediaManagement describes the media management data model.
type MediaManagement struct {
	ChmodFolder               types.String `tfsdk:"chmod_folder"`
	RescanAfterRefresh        types.String `tfsdk:"rescan_after_refresh"`
	RecycleBinPath            types.String `tfsdk:"recycle_bin_path"`
	FileDate                  types.String `tfsdk:"file_date"`
	ExtraFileExtensions       types.String `tfsdk:"extra_file_extensions"`
	EpisodeTitleRequired      types.String `tfsdk:"episode_title_required"`
	DownloadPropersRepacks    types.String `tfsdk:"download_propers_repacks"`
	ChownGroup                types.String `tfsdk:"chown_group"`
	ID                        types.Int64  `tfsdk:"id"`
	MinimumFreeSpace          types.Int64  `tfsdk:"minimum_free_space"`
	RecycleBinDays            types.Int64  `tfsdk:"recycle_bin_days"`
	UnmonitorPreviousEpisodes types.Bool   `tfsdk:"unmonitor_previous_episodes"`
	SkipFreeSpaceCheck        types.Bool   `tfsdk:"skip_free_space_check"`
	SetPermissions            types.Bool   `tfsdk:"set_permissions"`
	ImportExtraFiles          types.Bool   `tfsdk:"import_extra_files"`
	EnableMediaInfo           types.Bool   `tfsdk:"enable_media_info"`
	DeleteEmptyFolders        types.Bool   `tfsdk:"delete_empty_folders"`
	CreateEmptyFolders        types.Bool   `tfsdk:"create_empty_folders"`
	HardlinksCopy             types.Bool   `tfsdk:"hardlinks_copy"`
}

func (r *MediaManagementResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_" + mediaManagementResourceName
}

func (r *MediaManagementResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "<!-- subcategory:Media Management -->Media Management resource.\nFor more information refer to [Naming](https://wiki.servarr.com/sonarr/settings#file-management) documentation.",
		Attributes: map[string]schema.Attribute{
			"id": schema.Int64Attribute{
				MarkdownDescription: "Media Management ID.",
				Computed:            true,
				PlanModifiers: []planmodifier.Int64{
					int64planmodifier.UseStateForUnknown(),
				},
			},
			"unmonitor_previous_episodes": schema.BoolAttribute{
				MarkdownDescription: "Unmonitor deleted files.",
				Required:            true,
			},
			"hardlinks_copy": schema.BoolAttribute{
				MarkdownDescription: "Use hardlinks instead of copy.",
				Required:            true,
			},
			"create_empty_folders": schema.BoolAttribute{
				MarkdownDescription: "Create empty series directories.",
				Required:            true,
			},
			"delete_empty_folders": schema.BoolAttribute{
				MarkdownDescription: "Delete empty series directories.",
				Required:            true,
			},
			"enable_media_info": schema.BoolAttribute{
				MarkdownDescription: "Scan files details.",
				Required:            true,
			},
			"import_extra_files": schema.BoolAttribute{
				MarkdownDescription: "Import extra files. If enabled it will leverage 'extra_file_extensions'.",
				Required:            true,
			},
			"set_permissions": schema.BoolAttribute{
				MarkdownDescription: "Set permission for imported files.",
				Required:            true,
			},
			"skip_free_space_check": schema.BoolAttribute{
				MarkdownDescription: "Skip free space check before importing.",
				Required:            true,
			},
			"minimum_free_space": schema.Int64Attribute{
				MarkdownDescription: "Minimum free space in MB to allow import.",
				Required:            true,
			},
			"recycle_bin_days": schema.Int64Attribute{
				MarkdownDescription: "Recyle bin days of retention.",
				Required:            true,
			},
			"chmod_folder": schema.StringAttribute{
				MarkdownDescription: "Permission in linux format.",
				Required:            true,
			},
			"chown_group": schema.StringAttribute{
				MarkdownDescription: "Group used for permission.",
				Required:            true,
			},
			"download_propers_repacks": schema.StringAttribute{
				MarkdownDescription: "Download proper and repack policy. valid inputs are: 'preferAndUpgrade', 'doNotUpgrade', and 'doNotPrefer'.",
				Required:            true,
				Validators: []validator.String{
					stringvalidator.OneOf("preferAndUpgrade", "doNotUpgrade", "doNotPrefer"),
				},
			},
			"episode_title_required": schema.StringAttribute{
				MarkdownDescription: "Episode title requirement policy. valid inputs are: 'always', 'bulkSeasonReleases' and 'never'.",
				Required:            true,
				Validators: []validator.String{
					stringvalidator.OneOf("always", "bulkSeasonReleases", "never"),
				},
			},
			"extra_file_extensions": schema.StringAttribute{
				MarkdownDescription: "Comma separated list of extra files to import (.nfo will be imported as .nfo-orig).",
				Required:            true,
			},
			"file_date": schema.StringAttribute{
				MarkdownDescription: "Define the file date modification. valid inputs are: 'none', 'localAirDate, and 'utcAirDate'.",
				Required:            true,
				Validators: []validator.String{
					stringvalidator.OneOf("none", "localAirDate", "utcAirDate"),
				},
			},
			"recycle_bin_path": schema.StringAttribute{
				MarkdownDescription: "Recycle bin absolute path.",
				Required:            true,
			},
			"rescan_after_refresh": schema.StringAttribute{
				MarkdownDescription: "Rescan after refresh policy. valid inputs are: 'always', 'afterManual' and 'never'.",
				Required:            true,
				Validators: []validator.String{
					stringvalidator.OneOf("always", "afterManual", "never"),
				},
			},
		},
	}
}

func (r *MediaManagementResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	if auth, client := resourceConfigure(ctx, req, resp); client != nil {
		r.client = client
		r.auth = auth
	}
}

func (r *MediaManagementResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	// Retrieve values from plan
	var management *MediaManagement

	resp.Diagnostics.Append(req.Plan.Get(ctx, &management)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Build Create resource
	request := management.read()
	request.SetId(1)

	// Create new MediaManagement
	response, _, err := r.client.MediaManagementConfigAPI.UpdateMediaManagementConfig(r.auth, strconv.Itoa(int(request.GetId()))).MediaManagementConfigResource(*request).Execute()
	if err != nil {
		resp.Diagnostics.AddError(helpers.ClientError, helpers.ParseClientError(helpers.Create, mediaManagementResourceName, err))

		return
	}

	tflog.Trace(ctx, "created media_management: "+strconv.Itoa(int(response.GetId())))
	// Generate resource state struct
	management.write(response)
	resp.Diagnostics.Append(resp.State.Set(ctx, &management)...)
}

func (r *MediaManagementResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	// Get current state
	var management *MediaManagement

	resp.Diagnostics.Append(req.State.Get(ctx, &management)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Get mediamanagement current value
	response, _, err := r.client.MediaManagementConfigAPI.GetMediaManagementConfig(r.auth).Execute()
	if err != nil {
		resp.Diagnostics.AddError(helpers.ClientError, helpers.ParseClientError(helpers.Read, mediaManagementResourceName, err))

		return
	}

	tflog.Trace(ctx, "read "+mediaManagementResourceName+": "+strconv.Itoa(int(response.GetId())))
	// Map response body to resource schema attribute
	management.write(response)
	resp.Diagnostics.Append(resp.State.Set(ctx, &management)...)
}

func (r *MediaManagementResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	// Get plan values
	var management *MediaManagement

	resp.Diagnostics.Append(req.Plan.Get(ctx, &management)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Build Update resource
	request := management.read()

	// Update MediaManagement
	response, _, err := r.client.MediaManagementConfigAPI.UpdateMediaManagementConfig(r.auth, strconv.Itoa(int(request.GetId()))).MediaManagementConfigResource(*request).Execute()
	if err != nil {
		resp.Diagnostics.AddError(helpers.ClientError, helpers.ParseClientError(helpers.Update, mediaManagementResourceName, err))

		return
	}

	tflog.Trace(ctx, "updated "+mediaManagementResourceName+": "+strconv.Itoa(int(response.GetId())))
	// Generate resource state struct
	management.write(response)
	resp.Diagnostics.Append(resp.State.Set(ctx, &management)...)
}

func (r *MediaManagementResource) Delete(ctx context.Context, _ resource.DeleteRequest, resp *resource.DeleteResponse) {
	// Mediamanagement cannot be really deleted just removing configuration
	tflog.Trace(ctx, "decoupled "+mediaManagementResourceName+": 1")
	resp.State.RemoveResource(ctx)
}

func (r *MediaManagementResource) ImportState(ctx context.Context, _ resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	tflog.Trace(ctx, "imported "+mediaManagementResourceName+": 1")
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("id"), 1)...)
}

func (m *MediaManagement) write(mediaMgt *sonarr.MediaManagementConfigResource) {
	m.UnmonitorPreviousEpisodes = types.BoolValue(mediaMgt.GetAutoUnmonitorPreviouslyDownloadedEpisodes())
	m.HardlinksCopy = types.BoolValue(mediaMgt.GetCopyUsingHardlinks())
	m.CreateEmptyFolders = types.BoolValue(mediaMgt.GetCreateEmptySeriesFolders())
	m.DeleteEmptyFolders = types.BoolValue(mediaMgt.GetDeleteEmptyFolders())
	m.EnableMediaInfo = types.BoolValue(mediaMgt.GetEnableMediaInfo())
	m.ImportExtraFiles = types.BoolValue(mediaMgt.GetImportExtraFiles())
	m.SetPermissions = types.BoolValue(mediaMgt.GetSetPermissionsLinux())
	m.SkipFreeSpaceCheck = types.BoolValue(mediaMgt.GetSkipFreeSpaceCheckWhenImporting())
	m.ID = types.Int64Value(int64(mediaMgt.GetId()))
	m.MinimumFreeSpace = types.Int64Value(int64(mediaMgt.GetMinimumFreeSpaceWhenImporting()))
	m.RecycleBinDays = types.Int64Value(int64(mediaMgt.GetRecycleBinCleanupDays()))
	m.ChmodFolder = types.StringValue(mediaMgt.GetChmodFolder())
	m.ChownGroup = types.StringValue(mediaMgt.GetChownGroup())
	m.DownloadPropersRepacks = types.StringValue(string(mediaMgt.GetDownloadPropersAndRepacks()))
	m.EpisodeTitleRequired = types.StringValue(string(mediaMgt.GetEpisodeTitleRequired()))
	m.ExtraFileExtensions = types.StringValue(mediaMgt.GetExtraFileExtensions())
	m.FileDate = types.StringValue(string(mediaMgt.GetFileDate()))
	m.RecycleBinPath = types.StringValue(mediaMgt.GetRecycleBin())
	m.RescanAfterRefresh = types.StringValue(string(mediaMgt.GetRescanAfterRefresh()))
}

func (m *MediaManagement) read() *sonarr.MediaManagementConfigResource {
	mediaMgt := sonarr.NewMediaManagementConfigResource()
	mediaMgt.SetAutoUnmonitorPreviouslyDownloadedEpisodes(m.UnmonitorPreviousEpisodes.ValueBool())
	mediaMgt.SetCopyUsingHardlinks(m.HardlinksCopy.ValueBool())
	mediaMgt.SetCreateEmptySeriesFolders(m.CreateEmptyFolders.ValueBool())
	mediaMgt.SetDeleteEmptyFolders(m.DeleteEmptyFolders.ValueBool())
	mediaMgt.SetEnableMediaInfo(m.EnableMediaInfo.ValueBool())
	mediaMgt.SetImportExtraFiles(m.ImportExtraFiles.ValueBool())
	mediaMgt.SetSetPermissionsLinux(m.SetPermissions.ValueBool())
	mediaMgt.SetSkipFreeSpaceCheckWhenImporting(m.SkipFreeSpaceCheck.ValueBool())
	mediaMgt.SetId(int32(m.ID.ValueInt64()))
	mediaMgt.SetMinimumFreeSpaceWhenImporting(int32(m.MinimumFreeSpace.ValueInt64()))
	mediaMgt.SetRecycleBinCleanupDays(int32(m.RecycleBinDays.ValueInt64()))
	mediaMgt.SetChmodFolder(m.ChmodFolder.ValueString())
	mediaMgt.SetChownGroup(m.ChownGroup.ValueString())
	mediaMgt.SetDownloadPropersAndRepacks(sonarr.ProperDownloadTypes(m.DownloadPropersRepacks.ValueString()))
	mediaMgt.SetEpisodeTitleRequired(sonarr.EpisodeTitleRequiredType(m.EpisodeTitleRequired.ValueString()))
	mediaMgt.SetExtraFileExtensions(m.ExtraFileExtensions.ValueString())
	mediaMgt.SetFileDate(sonarr.FileDateType(m.FileDate.ValueString()))
	mediaMgt.SetRecycleBin(m.RecycleBinPath.ValueString())
	mediaMgt.SetRescanAfterRefresh(sonarr.RescanAfterRefreshType(m.RescanAfterRefresh.ValueString()))

	return mediaMgt
}

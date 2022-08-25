package provider

import (
	"context"
	"fmt"
	"strconv"

	"github.com/devopsarr/terraform-provider-sonarr/internal/helpers"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"golift.io/starr/sonarr"
)

// Ensure provider defined types fully satisfy framework interfaces.
var (
	_ provider.ResourceType            = resourceMediaManagementType{}
	_ resource.Resource                = resourceMediaManagement{}
	_ resource.ResourceWithImportState = resourceMediaManagement{}
)

type resourceMediaManagementType struct{}

type resourceMediaManagement struct {
	provider sonarrProvider
}

// MediaManagement is the MediaManagement resource.
type MediaManagement struct {
	UnmonitorPreviousEpisodes types.Bool   `tfsdk:"unmonitor_previous_episodes"`
	HardlinksCopy             types.Bool   `tfsdk:"hardlinks_copy"`
	CreateEmptyFolders        types.Bool   `tfsdk:"create_empty_folders"`
	DeleteEmptyFolders        types.Bool   `tfsdk:"delete_empty_folders"`
	EnableMediaInfo           types.Bool   `tfsdk:"enable_media_info"`
	ImportExtraFiles          types.Bool   `tfsdk:"import_extra_files"`
	SetPermissions            types.Bool   `tfsdk:"set_permissions"`
	SkipFreeSpaceCheck        types.Bool   `tfsdk:"skip_free_space_check"`
	ID                        types.Int64  `tfsdk:"id"`
	MinimumFreeSpace          types.Int64  `tfsdk:"minimum_free_space"`
	RecycleBinDays            types.Int64  `tfsdk:"recycle_bin_days"`
	ChmodFolder               types.String `tfsdk:"chmod_folder"`
	ChownGroup                types.String `tfsdk:"chown_group"`
	DownloadPropersRepacks    types.String `tfsdk:"download_propers_repacks"`
	EpisodeTitleRequired      types.String `tfsdk:"episode_title_required"`
	ExtraFileExtensions       types.String `tfsdk:"extra_file_extensions"`
	FileDate                  types.String `tfsdk:"file_date"`
	RecycleBinPath            types.String `tfsdk:"recycle_bin_path"`
	RescanAfterRefresh        types.String `tfsdk:"rescan_after_refresh"`
}

func (t resourceMediaManagementType) GetSchema(ctx context.Context) (tfsdk.Schema, diag.Diagnostics) {
	return tfsdk.Schema{
		MarkdownDescription: "MediaManagement resource",
		Attributes: map[string]tfsdk.Attribute{
			"id": {
				MarkdownDescription: "ID of mediamanagement",
				Computed:            true,
				Type:                types.Int64Type,
				PlanModifiers: tfsdk.AttributePlanModifiers{
					resource.UseStateForUnknown(),
				},
			},
			"unmonitor_previous_episodes": {
				MarkdownDescription: "Unmonitor deleted files",
				Required:            true,
				Type:                types.BoolType,
			},
			"hardlinks_copy": {
				MarkdownDescription: "Use hardlinks instead of copy",
				Required:            true,
				Type:                types.BoolType,
			},
			"create_empty_folders": {
				MarkdownDescription: "Create empty series directories",
				Required:            true,
				Type:                types.BoolType,
			},
			"delete_empty_folders": {
				MarkdownDescription: "Delete empty series directories",
				Required:            true,
				Type:                types.BoolType,
			},
			"enable_media_info": {
				MarkdownDescription: "Scan files details",
				Required:            true,
				Type:                types.BoolType,
			},
			"import_extra_files": {
				MarkdownDescription: "Import extra files. If enabled it will leverage 'extra_file_extensions'",
				Required:            true,
				Type:                types.BoolType,
			},
			"set_permissions": {
				MarkdownDescription: "Set permission for imported files",
				Required:            true,
				Type:                types.BoolType,
			},
			"skip_free_space_check": {
				MarkdownDescription: "Skip free space check before importing",
				Required:            true,
				Type:                types.BoolType,
			},
			"minimum_free_space": {
				MarkdownDescription: "Minimum free space in MB to allow import",
				Required:            true,
				Type:                types.Int64Type,
			},
			"recycle_bin_days": {
				MarkdownDescription: "Recyle bin days of retention",
				Required:            true,
				Type:                types.Int64Type,
			},
			"chmod_folder": {
				MarkdownDescription: "Permission in linux format",
				Required:            true,
				Type:                types.StringType,
			},
			"chown_group": {
				MarkdownDescription: "Group used for permission",
				Required:            true,
				Type:                types.StringType,
			},
			"download_propers_repacks": {
				MarkdownDescription: "Download proper and repack policy. valid inputs are: 'preferAndUpgrade', 'doNotUpgrade', and 'doNotPrefer'",
				Required:            true,
				Type:                types.StringType,
				Validators: []tfsdk.AttributeValidator{
					helpers.StringMatch([]string{"preferAndUpgrade", "doNotUpgrade", "doNotPrefer"}),
				},
			},
			"episode_title_required": {
				MarkdownDescription: "Episode title requirement policy. valid inputs are: 'always', 'bulkSeasonReleases' and 'never'",
				Required:            true,
				Type:                types.StringType,
				Validators: []tfsdk.AttributeValidator{
					helpers.StringMatch([]string{"always", "bulkSeasonReleases", "never"}),
				},
			},
			"extra_file_extensions": {
				MarkdownDescription: "Comma separated list of extra files to import (.nfo will be imported as .nfo-orig)",
				Required:            true,
				Type:                types.StringType,
			},
			"file_date": {
				MarkdownDescription: "Define the file date modification. valid inputs are: 'none', 'localAirDate, and 'utcAirDate'",
				Required:            true,
				Type:                types.StringType,
				Validators: []tfsdk.AttributeValidator{
					helpers.StringMatch([]string{"none", "localAirDate", "utcAirDate"}),
				},
			},
			"recycle_bin_path": {
				MarkdownDescription: "Recycle bin absolute path",
				Required:            true,
				Type:                types.StringType,
			},
			"rescan_after_refresh": {
				MarkdownDescription: "Rescan after refresh policy. valid inputs are: 'always', 'afterManual' and 'never'",
				Required:            true,
				Type:                types.StringType,
				Validators: []tfsdk.AttributeValidator{
					helpers.StringMatch([]string{"always", "afterManual", "never"}),
				},
			},
		},
	}, nil
}

func (t resourceMediaManagementType) NewResource(ctx context.Context, in provider.Provider) (resource.Resource, diag.Diagnostics) {
	provider, diags := convertProviderType(in)

	return resourceMediaManagement{
		provider: provider,
	}, diags
}

func (r resourceMediaManagement) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	// Retrieve values from plan
	var plan MediaManagement
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Build Create resource
	data := readMediaManagement(&plan)
	data.ID = 1

	// Create new MediaManagement
	response, err := r.provider.client.UpdateMediaManagementContext(ctx, data)
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to create mediamanagement, got error: %s", err))

		return
	}

	tflog.Trace(ctx, "created mediamanagement: "+strconv.Itoa(int(response.ID)))

	// Generate resource state struct
	result := writeMediaManagement(response)

	diags = resp.State.Set(ctx, result)
	resp.Diagnostics.Append(diags...)

	if resp.Diagnostics.HasError() {
		return
	}
}

func (r resourceMediaManagement) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	// Get current state
	var state MediaManagement
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Get mediamanagement current value
	response, err := r.provider.client.GetMediaManagementContext(ctx)
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to read mediamanagements, got error: %s", err))

		return
	}
	// Map response body to resource schema attribute
	result := writeMediaManagement(response)

	diags = resp.State.Set(ctx, result)
	resp.Diagnostics.Append(diags...)
}

func (r resourceMediaManagement) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	// Get plan values
	var plan MediaManagement
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Build Update resource
	data := readMediaManagement(&plan)

	// Update MediaManagement
	response, err := r.provider.client.UpdateMediaManagementContext(ctx, data)
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to update mediamanagement, got error: %s", err))

		return
	}

	tflog.Trace(ctx, "update mediamanagement: "+strconv.Itoa(int(response.ID)))

	// Generate resource state struct
	result := writeMediaManagement(response)

	diags = resp.State.Set(ctx, result)
	resp.Diagnostics.Append(diags...)

	if resp.Diagnostics.HasError() {
		return
	}
}

func (r resourceMediaManagement) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	// Mediamanagement cannot be really deleted just removing configuration
	resp.State.RemoveResource(ctx)
}

func (r resourceMediaManagement) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	// resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("id"), 1)...)
}

func writeMediaManagement(mediaMgt *sonarr.MediaManagement) *MediaManagement {
	return &MediaManagement{
		UnmonitorPreviousEpisodes: types.Bool{Value: mediaMgt.AutoUnmonitorPreviouslyDownloadedEpisodes},
		HardlinksCopy:             types.Bool{Value: mediaMgt.CopyUsingHardlinks},
		CreateEmptyFolders:        types.Bool{Value: mediaMgt.CreateEmptySeriesFolders},
		DeleteEmptyFolders:        types.Bool{Value: mediaMgt.DeleteEmptyFolders},
		EnableMediaInfo:           types.Bool{Value: mediaMgt.EnableMediaInfo},
		ImportExtraFiles:          types.Bool{Value: mediaMgt.ImportExtraFiles},
		SetPermissions:            types.Bool{Value: mediaMgt.SetPermissionsLinux},
		SkipFreeSpaceCheck:        types.Bool{Value: mediaMgt.SkipFreeSpaceCheckWhenImporting},
		ID:                        types.Int64{Value: mediaMgt.ID},
		MinimumFreeSpace:          types.Int64{Value: mediaMgt.MinimumFreeSpaceWhenImporting},
		RecycleBinDays:            types.Int64{Value: mediaMgt.RecycleBinCleanupDays},
		ChmodFolder:               types.String{Value: mediaMgt.ChmodFolder},
		ChownGroup:                types.String{Value: mediaMgt.ChownGroup},
		DownloadPropersRepacks:    types.String{Value: mediaMgt.DownloadPropersAndRepacks},
		EpisodeTitleRequired:      types.String{Value: mediaMgt.EpisodeTitleRequired},
		ExtraFileExtensions:       types.String{Value: mediaMgt.ExtraFileExtensions},
		FileDate:                  types.String{Value: mediaMgt.FileDate},
		RecycleBinPath:            types.String{Value: mediaMgt.RecycleBin},
		RescanAfterRefresh:        types.String{Value: mediaMgt.RescanAfterRefresh},
	}
}

func readMediaManagement(mediaMgt *MediaManagement) *sonarr.MediaManagement {
	return &sonarr.MediaManagement{
		AutoUnmonitorPreviouslyDownloadedEpisodes: mediaMgt.UnmonitorPreviousEpisodes.Value,
		CopyUsingHardlinks:                        mediaMgt.HardlinksCopy.Value,
		CreateEmptySeriesFolders:                  mediaMgt.CreateEmptyFolders.Value,
		DeleteEmptyFolders:                        mediaMgt.DeleteEmptyFolders.Value,
		EnableMediaInfo:                           mediaMgt.EnableMediaInfo.Value,
		ImportExtraFiles:                          mediaMgt.ImportExtraFiles.Value,
		SetPermissionsLinux:                       mediaMgt.SetPermissions.Value,
		SkipFreeSpaceCheckWhenImporting:           mediaMgt.SkipFreeSpaceCheck.Value,
		ID:                                        mediaMgt.ID.Value,
		MinimumFreeSpaceWhenImporting:             mediaMgt.MinimumFreeSpace.Value,
		RecycleBinCleanupDays:                     mediaMgt.RecycleBinDays.Value,
		ChmodFolder:                               mediaMgt.ChmodFolder.Value,
		ChownGroup:                                mediaMgt.ChownGroup.Value,
		DownloadPropersAndRepacks:                 mediaMgt.DownloadPropersRepacks.Value,
		EpisodeTitleRequired:                      mediaMgt.EpisodeTitleRequired.Value,
		ExtraFileExtensions:                       mediaMgt.ExtraFileExtensions.Value,
		FileDate:                                  mediaMgt.FileDate.Value,
		RecycleBin:                                mediaMgt.RecycleBinPath.Value,
		RescanAfterRefresh:                        mediaMgt.RescanAfterRefresh.Value,
	}
}

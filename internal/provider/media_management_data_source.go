package provider

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// Ensure provider defined types fully satisfy framework interfaces.
var (
	_ provider.DataSourceType = dataMediaManagementType{}
	_ datasource.DataSource   = dataMediaManagement{}
)

type dataMediaManagementType struct{}

type dataMediaManagement struct {
	provider sonarrProvider
}

func (t dataMediaManagementType) GetSchema(ctx context.Context) (tfsdk.Schema, diag.Diagnostics) {
	return tfsdk.Schema{
		// This description is used by the documentation generator and the delay server.
		MarkdownDescription: "[Media Management](../resources/media_management).",
		Attributes: map[string]tfsdk.Attribute{
			"id": {
				MarkdownDescription: "Delay Profile ID.",
				Computed:            true,
				Type:                types.Int64Type,
			},
			"unmonitor_previous_episodes": {
				MarkdownDescription: "Unmonitor deleted files.",
				Computed:            true,
				Type:                types.BoolType,
			},
			"hardlinks_copy": {
				MarkdownDescription: "Use hardlinks instead of copy.",
				Computed:            true,
				Type:                types.BoolType,
			},
			"create_empty_folders": {
				MarkdownDescription: "Create empty series directories.",
				Computed:            true,
				Type:                types.BoolType,
			},
			"delete_empty_folders": {
				MarkdownDescription: "Delete empty series directories.",
				Computed:            true,
				Type:                types.BoolType,
			},
			"enable_media_info": {
				MarkdownDescription: "Scan files details.",
				Computed:            true,
				Type:                types.BoolType,
			},
			"import_extra_files": {
				MarkdownDescription: "Import extra files. If enabled it will leverage 'extra_file_extensions'.",
				Computed:            true,
				Type:                types.BoolType,
			},
			"set_permissions": {
				MarkdownDescription: "Set permission for imported files.",
				Computed:            true,
				Type:                types.BoolType,
			},
			"skip_free_space_check": {
				MarkdownDescription: "Skip free space check before importing.",
				Computed:            true,
				Type:                types.BoolType,
			},
			"minimum_free_space": {
				MarkdownDescription: "Minimum free space in MB to allow import.",
				Computed:            true,
				Type:                types.Int64Type,
			},
			"recycle_bin_days": {
				MarkdownDescription: "Recyle bin days of retention.",
				Computed:            true,
				Type:                types.Int64Type,
			},
			"chmod_folder": {
				MarkdownDescription: "Permission in linux format.",
				Computed:            true,
				Type:                types.StringType,
			},
			"chown_group": {
				MarkdownDescription: "Group used for permission.",
				Computed:            true,
				Type:                types.StringType,
			},
			"download_propers_repacks": {
				MarkdownDescription: "Download proper and repack policy. valid inputs are: 'preferAndUpgrade', 'doNotUpgrade', and 'doNotPrefer'.",
				Computed:            true,
				Type:                types.StringType,
			},
			"episode_title_required": {
				MarkdownDescription: "Episode title requirement policy. valid inputs are: 'always', 'bulkSeasonReleases' and 'never'.",
				Computed:            true,
				Type:                types.StringType,
			},
			"extra_file_extensions": {
				MarkdownDescription: "Comma separated list of extra files to import (.nfo will be imported as .nfo-orig).",
				Computed:            true,
				Type:                types.StringType,
			},
			"file_date": {
				MarkdownDescription: "Define the file date modification. valid inputs are: 'none', 'localAirDate, and 'utcAirDate'.",
				Computed:            true,
				Type:                types.StringType,
			},
			"recycle_bin_path": {
				MarkdownDescription: "Recycle bin absolute path.",
				Computed:            true,
				Type:                types.StringType,
			},
			"rescan_after_refresh": {
				MarkdownDescription: "Rescan after refresh policy. valid inputs are: 'always', 'afterManual' and 'never'.",
				Computed:            true,
				Type:                types.StringType,
			},
		},
	}, nil
}

func (t dataMediaManagementType) NewDataSource(ctx context.Context, in provider.Provider) (datasource.DataSource, diag.Diagnostics) {
	provider, diags := convertProviderType(in)

	return dataMediaManagement{
		provider: provider,
	}, diags
}

func (d dataMediaManagement) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	// Get indexer config current value
	response, err := d.provider.client.GetMediaManagementContext(ctx)
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to read indexer cofig, got error: %s", err))

		return
	}

	result := writeMediaManagement(response)
	diags := resp.State.Set(ctx, &result)
	resp.Diagnostics.Append(diags...)
}

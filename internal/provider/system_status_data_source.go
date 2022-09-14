package provider

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"golift.io/starr/sonarr"
)

// Ensure provider defined types fully satisfy framework interfaces.
var _ datasource.DataSource = &SystemStatusDataSource{}

func NewSystemStatusDataSource() datasource.DataSource {
	return &SystemStatusDataSource{}
}

// SystemStatusDataSource defines the system status implementation.
type SystemStatusDataSource struct {
	client *sonarr.Sonarr
}

// SystemStatus describes the system status data model.
type SystemStatus struct {
	IsDebug           types.Bool `tfsdk:"is_debug"`
	IsProduction      types.Bool `tfsdk:"is_production"`
	IsAdmin           types.Bool `tfsdk:"is_admin"`
	IsUserInteractive types.Bool `tfsdk:"is_user_interactive"`
	IsMonoRuntime     types.Bool `tfsdk:"is_mono_runtime"`
	IsMono            types.Bool `tfsdk:"is_mono"`
	IsLinux           types.Bool `tfsdk:"is_linux"`
	IsOsx             types.Bool `tfsdk:"is_osx"`
	IsWindows         types.Bool `tfsdk:"is_windows"`
	// TODO: remove ID once framework support tests without ID https://www.terraform.io/plugin/framework/acctests#implement-id-attribute
	ID                     types.Int64  `tfsdk:"id"`
	Version                types.String `tfsdk:"version"`
	StartupPath            types.String `tfsdk:"startup_path"`
	AppData                types.String `tfsdk:"app_data"`
	OsName                 types.String `tfsdk:"os_name"`
	OsVersion              types.String `tfsdk:"os_version"`
	Mode                   types.String `tfsdk:"mode"`
	Branch                 types.String `tfsdk:"branch"`
	Authentication         types.String `tfsdk:"authentication"`
	SqliteVersion          types.String `tfsdk:"sqlite_version"`
	URLBase                types.String `tfsdk:"url_base"`
	RuntimeVersion         types.String `tfsdk:"runtime_version"`
	RuntimeName            types.String `tfsdk:"runtime_name"`
	PackageVersion         types.String `tfsdk:"package_version"`
	PackageAuthor          types.String `tfsdk:"package_author"`
	PackageUpdateMechanism types.String `tfsdk:"package_update_mechanism"`
	BuildTime              types.String `tfsdk:"build_time"`
	StartTime              types.String `tfsdk:"start_time"`
}

func (d *SystemStatusDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_system_status"
}

func (d *SystemStatusDataSource) GetSchema(ctx context.Context) (tfsdk.Schema, diag.Diagnostics) {
	return tfsdk.Schema{
		// This description is used by the documentation generator and the delay server.
		MarkdownDescription: "System Status resource. User must have rights to read config.xml.<br/>For more information refer to [System Status](https://wiki.servarr.com/sonarr/system#status) documentation.",
		Attributes: map[string]tfsdk.Attribute{
			// TODO: remove ID once framework support tests without ID https://www.terraform.io/plugin/framework/acctests#implement-id-attribute
			"id": {
				MarkdownDescription: "Delay Profile ID.",
				Computed:            true,
				Type:                types.Int64Type,
			},
			"is_debug": {
				MarkdownDescription: "Is debug flag.",
				Computed:            true,
				Type:                types.BoolType,
			},
			"is_production": {
				MarkdownDescription: "Is production flag.",
				Computed:            true,
				Type:                types.BoolType,
			},
			"is_admin": {
				MarkdownDescription: "Is admin flag.",
				Computed:            true,
				Type:                types.BoolType,
			},
			"is_user_interactive": {
				MarkdownDescription: "Is user interactive flag.",
				Computed:            true,
				Type:                types.BoolType,
			},
			"is_mono_runtime": {
				MarkdownDescription: "Is mono runtime flag.",
				Computed:            true,
				Type:                types.BoolType,
			},
			"is_mono": {
				MarkdownDescription: "Is mono flag.",
				Computed:            true,
				Type:                types.BoolType,
			},
			"is_linux": {
				MarkdownDescription: "Is linux flag.",
				Computed:            true,
				Type:                types.BoolType,
			},
			"is_osx": {
				MarkdownDescription: "Is osx flag.",
				Computed:            true,
				Type:                types.BoolType,
			},
			"is_windows": {
				MarkdownDescription: "Is windows flag.",
				Computed:            true,
				Type:                types.BoolType,
			},
			"version": {
				MarkdownDescription: "Version.",
				Computed:            true,
				Type:                types.StringType,
			},
			"startup_path": {
				MarkdownDescription: "Startup path.",
				Computed:            true,
				Type:                types.StringType,
			},
			"app_data": {
				MarkdownDescription: "App data folder.",
				Computed:            true,
				Type:                types.StringType,
			},
			"os_name": {
				MarkdownDescription: "OS name.",
				Computed:            true,
				Type:                types.StringType,
			},
			"os_version": {
				MarkdownDescription: "OS version.",
				Computed:            true,
				Type:                types.StringType,
			},
			"mode": {
				MarkdownDescription: "Mode.",
				Computed:            true,
				Type:                types.StringType,
			},
			"branch": {
				MarkdownDescription: "Branch.",
				Computed:            true,
				Type:                types.StringType,
			},
			"authentication": {
				MarkdownDescription: "Authentication.",
				Computed:            true,
				Type:                types.StringType,
			},
			"sqlite_version": {
				MarkdownDescription: "SQLite version.",
				Computed:            true,
				Type:                types.StringType,
			},
			"url_base": {
				MarkdownDescription: "Base URL.",
				Computed:            true,
				Type:                types.StringType,
			},
			"runtime_version": {
				MarkdownDescription: "Runtime version.",
				Computed:            true,
				Type:                types.StringType,
			},
			"runtime_name": {
				MarkdownDescription: "Runtime name.",
				Computed:            true,
				Type:                types.StringType,
			},
			"package_version": {
				MarkdownDescription: "Package version.",
				Computed:            true,
				Type:                types.StringType,
			},
			"package_author": {
				MarkdownDescription: "Package author.",
				Computed:            true,
				Type:                types.StringType,
			},
			"package_update_mechanism": {
				MarkdownDescription: "Package update mechanism.",
				Computed:            true,
				Type:                types.StringType,
			},
			"build_time": {
				MarkdownDescription: "Build time.",
				Computed:            true,
				Type:                types.StringType,
			},
			"start_time": {
				MarkdownDescription: "Start time.",
				Computed:            true,
				Type:                types.StringType,
			},
		},
	}, nil
}

func (d *SystemStatusDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	// Prevent panic if the provider has not been configured.
	if req.ProviderData == nil {
		return
	}

	client, ok := req.ProviderData.(*sonarr.Sonarr)
	if !ok {
		resp.Diagnostics.AddError(
			UnexpectedDataSourceConfigureType,
			fmt.Sprintf("Expected *sonarr.Sonarr, got: %T. Please report this issue to the provider developers.", req.ProviderData),
		)

		return
	}

	d.client = client
}

func (d *SystemStatusDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	// Get naming current value
	response, err := d.client.GetSystemStatusContext(ctx)
	if err != nil {
		resp.Diagnostics.AddError(ClientError, fmt.Sprintf("Unable to read system status, got error: %s", err))

		return
	}

	tflog.Trace(ctx, "read system status")

	result := writeSystemStatus(response)
	resp.Diagnostics.Append(resp.State.Set(ctx, &result)...)
}

func writeSystemStatus(status *sonarr.SystemStatus) *SystemStatus {
	return &SystemStatus{
		IsDebug:                types.Bool{Value: status.IsDebug},
		IsProduction:           types.Bool{Value: status.IsProduction},
		IsAdmin:                types.Bool{Value: status.IsProduction},
		IsUserInteractive:      types.Bool{Value: status.IsUserInteractive},
		IsMonoRuntime:          types.Bool{Value: status.IsMonoRuntime},
		IsMono:                 types.Bool{Value: status.IsMono},
		IsLinux:                types.Bool{Value: status.IsLinux},
		IsOsx:                  types.Bool{Value: status.IsOsx},
		IsWindows:              types.Bool{Value: status.IsWindows},
		ID:                     types.Int64{Value: int64(1)},
		Version:                types.String{Value: status.Version},
		StartupPath:            types.String{Value: status.StartupPath},
		AppData:                types.String{Value: status.AppData},
		OsName:                 types.String{Value: status.OsName},
		OsVersion:              types.String{Value: status.OsVersion},
		Mode:                   types.String{Value: status.Mode},
		Branch:                 types.String{Value: status.Branch},
		Authentication:         types.String{Value: status.Authentication},
		SqliteVersion:          types.String{Value: status.SqliteVersion},
		URLBase:                types.String{Value: status.URLBase},
		RuntimeVersion:         types.String{Value: status.RuntimeVersion},
		RuntimeName:            types.String{Value: status.RuntimeName},
		PackageVersion:         types.String{Value: status.PackageVersion},
		PackageAuthor:          types.String{Value: status.PackageAuthor},
		PackageUpdateMechanism: types.String{Value: status.PackageUpdateMechanism},
		BuildTime:              types.String{Value: status.BuildTime.String()},
		StartTime:              types.String{Value: status.StartTime.String()},
	}
}

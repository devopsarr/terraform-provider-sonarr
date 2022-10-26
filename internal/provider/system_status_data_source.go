package provider

import (
	"context"
	"fmt"

	"github.com/devopsarr/terraform-provider-sonarr/internal/helpers"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"golift.io/starr/sonarr"
)

const systemStatusDataSourceName = "system_status"

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
	SqliteVersion          types.String `tfsdk:"sqlite_version"`
	URLBase                types.String `tfsdk:"url_base"`
	AppData                types.String `tfsdk:"app_data"`
	OsName                 types.String `tfsdk:"os_name"`
	BuildTime              types.String `tfsdk:"build_time"`
	PackageUpdateMechanism types.String `tfsdk:"package_update_mechanism"`
	PackageAuthor          types.String `tfsdk:"package_author"`
	PackageVersion         types.String `tfsdk:"package_version"`
	OsVersion              types.String `tfsdk:"os_version"`
	RuntimeVersion         types.String `tfsdk:"runtime_version"`
	Version                types.String `tfsdk:"version"`
	StartupPath            types.String `tfsdk:"startup_path"`
	Authentication         types.String `tfsdk:"authentication"`
	StartTime              types.String `tfsdk:"start_time"`
	RuntimeName            types.String `tfsdk:"runtime_name"`
	Mode                   types.String `tfsdk:"mode"`
	Branch                 types.String `tfsdk:"branch"`
	ID                     types.Int64  `tfsdk:"id"`
	IsAdmin                types.Bool   `tfsdk:"is_admin"`
	IsDebug                types.Bool   `tfsdk:"is_debug"`
	IsProduction           types.Bool   `tfsdk:"is_production"`
	IsWindows              types.Bool   `tfsdk:"is_windows"`
	IsOsx                  types.Bool   `tfsdk:"is_osx"`
	IsLinux                types.Bool   `tfsdk:"is_linux"`
	IsMono                 types.Bool   `tfsdk:"is_mono"`
	IsMonoRuntime          types.Bool   `tfsdk:"is_mono_runtime"`
	IsUserInteractive      types.Bool   `tfsdk:"is_user_interactive"`
}

func (d *SystemStatusDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_" + systemStatusDataSourceName
}

func (d *SystemStatusDataSource) GetSchema(ctx context.Context) (tfsdk.Schema, diag.Diagnostics) {
	return tfsdk.Schema{
		// This description is used by the documentation generator and the delay server.
		MarkdownDescription: "[subcategory:Status]: #\nSystem Status resource. User must have rights to read `config.xml`.\nFor more information refer to [System Status](https://wiki.servarr.com/sonarr/system#status) documentation.",
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
			helpers.UnexpectedDataSourceConfigureType,
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
		resp.Diagnostics.AddError(helpers.ClientError, fmt.Sprintf("Unable to read %s, got error: %s", systemStatusDataSourceName, err))

		return
	}

	tflog.Trace(ctx, "read "+systemStatusDataSourceName)

	status := SystemStatus{}
	status.write(response)
	resp.Diagnostics.Append(resp.State.Set(ctx, status)...)
}

func (s *SystemStatus) write(status *sonarr.SystemStatus) {
	s.IsDebug = types.Bool{Value: status.IsDebug}
	s.IsProduction = types.Bool{Value: status.IsProduction}
	s.IsAdmin = types.Bool{Value: status.IsProduction}
	s.IsUserInteractive = types.Bool{Value: status.IsUserInteractive}
	s.IsMonoRuntime = types.Bool{Value: status.IsMonoRuntime}
	s.IsMono = types.Bool{Value: status.IsMono}
	s.IsLinux = types.Bool{Value: status.IsLinux}
	s.IsOsx = types.Bool{Value: status.IsOsx}
	s.IsWindows = types.Bool{Value: status.IsWindows}
	s.ID = types.Int64{Value: int64(1)}
	s.Version = types.String{Value: status.Version}
	s.StartupPath = types.String{Value: status.StartupPath}
	s.AppData = types.String{Value: status.AppData}
	s.OsName = types.String{Value: status.OsName}
	s.OsVersion = types.String{Value: status.OsVersion}
	s.Mode = types.String{Value: status.Mode}
	s.Branch = types.String{Value: status.Branch}
	s.Authentication = types.String{Value: status.Authentication}
	s.SqliteVersion = types.String{Value: status.SqliteVersion}
	s.URLBase = types.String{Value: status.URLBase}
	s.RuntimeVersion = types.String{Value: status.RuntimeVersion}
	s.RuntimeName = types.String{Value: status.RuntimeName}
	s.PackageVersion = types.String{Value: status.PackageVersion}
	s.PackageAuthor = types.String{Value: status.PackageAuthor}
	s.PackageUpdateMechanism = types.String{Value: status.PackageUpdateMechanism}
	s.BuildTime = types.String{Value: status.BuildTime.String()}
	s.StartTime = types.String{Value: status.StartTime.String()}
}

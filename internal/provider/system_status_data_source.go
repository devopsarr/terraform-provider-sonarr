package provider

import (
	"context"
	"fmt"

	"github.com/devopsarr/terraform-provider-sonarr/tools"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
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

func (d *SystemStatusDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		// This description is used by the documentation generator and the delay server.
		MarkdownDescription: "<!-- subcategory:Status -->System Status resource. User must have rights to read `config.xml`.\nFor more information refer to [System Status](https://wiki.servarr.com/sonarr/system#status) documentation.",
		Attributes: map[string]schema.Attribute{
			// TODO: remove ID once framework support tests without ID https://www.terraform.io/plugin/framework/acctests#implement-id-attribute
			"id": schema.Int64Attribute{
				MarkdownDescription: "Delay Profile ID.",
				Computed:            true,
			},
			"is_debug": schema.BoolAttribute{
				MarkdownDescription: "Is debug flag.",
				Computed:            true,
			},
			"is_production": schema.BoolAttribute{
				MarkdownDescription: "Is production flag.",
				Computed:            true,
			},
			"is_admin": schema.BoolAttribute{
				MarkdownDescription: "Is admin flag.",
				Computed:            true,
			},
			"is_user_interactive": schema.BoolAttribute{
				MarkdownDescription: "Is user interactive flag.",
				Computed:            true,
			},
			"is_mono_runtime": schema.BoolAttribute{
				MarkdownDescription: "Is mono runtime flag.",
				Computed:            true,
			},
			"is_mono": schema.BoolAttribute{
				MarkdownDescription: "Is mono flag.",
				Computed:            true,
			},
			"is_linux": schema.BoolAttribute{
				MarkdownDescription: "Is linux flag.",
				Computed:            true,
			},
			"is_osx": schema.BoolAttribute{
				MarkdownDescription: "Is osx flag.",
				Computed:            true,
			},
			"is_windows": schema.BoolAttribute{
				MarkdownDescription: "Is windows flag.",
				Computed:            true,
			},
			"version": schema.StringAttribute{
				MarkdownDescription: "Version.",
				Computed:            true,
			},
			"startup_path": schema.StringAttribute{
				MarkdownDescription: "Startup path.",
				Computed:            true,
			},
			"app_data": schema.StringAttribute{
				MarkdownDescription: "App data folder.",
				Computed:            true,
			},
			"os_name": schema.StringAttribute{
				MarkdownDescription: "OS name.",
				Computed:            true,
			},
			"os_version": schema.StringAttribute{
				MarkdownDescription: "OS version.",
				Computed:            true,
			},
			"mode": schema.StringAttribute{
				MarkdownDescription: "Mode.",
				Computed:            true,
			},
			"branch": schema.StringAttribute{
				MarkdownDescription: "Branch.",
				Computed:            true,
			},
			"authentication": schema.StringAttribute{
				MarkdownDescription: "Authentication.",
				Computed:            true,
			},
			"sqlite_version": schema.StringAttribute{
				MarkdownDescription: "SQLite version.",
				Computed:            true,
			},
			"url_base": schema.StringAttribute{
				MarkdownDescription: "Base URL.",
				Computed:            true,
			},
			"runtime_version": schema.StringAttribute{
				MarkdownDescription: "Runtime version.",
				Computed:            true,
			},
			"runtime_name": schema.StringAttribute{
				MarkdownDescription: "Runtime name.",
				Computed:            true,
			},
			"package_version": schema.StringAttribute{
				MarkdownDescription: "Package version.",
				Computed:            true,
			},
			"package_author": schema.StringAttribute{
				MarkdownDescription: "Package author.",
				Computed:            true,
			},
			"package_update_mechanism": schema.StringAttribute{
				MarkdownDescription: "Package update mechanism.",
				Computed:            true,
			},
			"build_time": schema.StringAttribute{
				MarkdownDescription: "Build time.",
				Computed:            true,
			},
			"start_time": schema.StringAttribute{
				MarkdownDescription: "Start time.",
				Computed:            true,
			},
		},
	}
}

func (d *SystemStatusDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	// Prevent panic if the provider has not been configured.
	if req.ProviderData == nil {
		return
	}

	client, ok := req.ProviderData.(*sonarr.Sonarr)
	if !ok {
		resp.Diagnostics.AddError(
			tools.UnexpectedDataSourceConfigureType,
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
		resp.Diagnostics.AddError(tools.ClientError, fmt.Sprintf("Unable to read %s, got error: %s", systemStatusDataSourceName, err))

		return
	}

	tflog.Trace(ctx, "read "+systemStatusDataSourceName)

	status := SystemStatus{}
	status.write(response)
	resp.Diagnostics.Append(resp.State.Set(ctx, status)...)
}

func (s *SystemStatus) write(status *sonarr.SystemStatus) {
	s.IsDebug = types.BoolValue(status.IsDebug)
	s.IsProduction = types.BoolValue(status.IsProduction)
	s.IsAdmin = types.BoolValue(status.IsProduction)
	s.IsUserInteractive = types.BoolValue(status.IsUserInteractive)
	s.IsMonoRuntime = types.BoolValue(status.IsMonoRuntime)
	s.IsMono = types.BoolValue(status.IsMono)
	s.IsLinux = types.BoolValue(status.IsLinux)
	s.IsOsx = types.BoolValue(status.IsOsx)
	s.IsWindows = types.BoolValue(status.IsWindows)
	s.ID = types.Int64Value(int64(1))
	s.Version = types.StringValue(status.Version)
	s.StartupPath = types.StringValue(status.StartupPath)
	s.AppData = types.StringValue(status.AppData)
	s.OsName = types.StringValue(status.OsName)
	s.OsVersion = types.StringValue(status.OsVersion)
	s.Mode = types.StringValue(status.Mode)
	s.Branch = types.StringValue(status.Branch)
	s.Authentication = types.StringValue(status.Authentication)
	s.SqliteVersion = types.StringValue(status.SqliteVersion)
	s.URLBase = types.StringValue(status.URLBase)
	s.RuntimeVersion = types.StringValue(status.RuntimeVersion)
	s.RuntimeName = types.StringValue(status.RuntimeName)
	s.PackageVersion = types.StringValue(status.PackageVersion)
	s.PackageAuthor = types.StringValue(status.PackageAuthor)
	s.PackageUpdateMechanism = types.StringValue(status.PackageUpdateMechanism)
	s.BuildTime = types.StringValue(status.BuildTime.String())
	s.StartTime = types.StringValue(status.StartTime.String())
}

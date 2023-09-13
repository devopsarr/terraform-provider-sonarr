package provider

import (
	"context"

	"github.com/devopsarr/sonarr-go/sonarr"
	"github.com/devopsarr/terraform-provider-sonarr/internal/helpers"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

const hostDataSourceName = "host"

// Ensure provider defined types fully satisfy framework interfaces.
var _ datasource.DataSource = &HostDataSource{}

func NewHostDataSource() datasource.DataSource {
	return &HostDataSource{}
}

// HostDataSource defines the host implementation.
type HostDataSource struct {
	client *sonarr.APIClient
}

func (d *HostDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_" + hostDataSourceName
}

func (d *HostDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		// This description is used by the documentation generator and the delay server.
		MarkdownDescription: "<!-- subcategory:System -->[Host](../resources/host).",
		Attributes: map[string]schema.Attribute{
			"launch_browser": schema.BoolAttribute{
				MarkdownDescription: "Launch browser flag.",
				Computed:            true,
			},
			"port": schema.Int64Attribute{
				MarkdownDescription: "TCP port.",
				Computed:            true,
			},
			"id": schema.Int64Attribute{
				MarkdownDescription: "Host ID.",
				Computed:            true,
			},
			"url_base": schema.StringAttribute{
				MarkdownDescription: "URL base.",
				Computed:            true,
			},
			"bind_address": schema.StringAttribute{
				MarkdownDescription: "Bind address.",
				Computed:            true,
			},
			"application_url": schema.StringAttribute{
				MarkdownDescription: "Application URL.",
				Computed:            true,
			},
			"instance_name": schema.StringAttribute{
				MarkdownDescription: "Instance name.",
				Computed:            true,
			},
			"update": schema.SingleNestedAttribute{
				MarkdownDescription: "Update configuration.",
				Computed:            true,
				Attributes: map[string]schema.Attribute{
					"mechanism": schema.StringAttribute{
						MarkdownDescription: "Update mechanism.",
						Computed:            true,
					},
					"script_path": schema.StringAttribute{
						MarkdownDescription: "Script path.",
						Computed:            true,
					},
					"branch": schema.StringAttribute{
						MarkdownDescription: "Branch reference.",
						Computed:            true,
					},
					"update_automatically": schema.BoolAttribute{
						MarkdownDescription: "Update automatically flag.",
						Computed:            true,
					},
				},
			},
			"logging": schema.SingleNestedAttribute{
				MarkdownDescription: "Logging configuration.",
				Computed:            true,
				Attributes: map[string]schema.Attribute{
					"log_level": schema.StringAttribute{
						MarkdownDescription: "Log level.",
						Computed:            true,
					},
					"console_log_level": schema.StringAttribute{
						MarkdownDescription: "Console log level.",
						Computed:            true,
					},
					"analytics_enabled": schema.BoolAttribute{
						MarkdownDescription: "Enable analytics flag.",
						Computed:            true,
					},
				},
			},
			"backup": schema.SingleNestedAttribute{
				MarkdownDescription: "Backup configuration.",
				Computed:            true,
				Attributes: map[string]schema.Attribute{
					"folder": schema.StringAttribute{
						MarkdownDescription: "Backup folder.",
						Computed:            true,
					},
					"interval": schema.Int64Attribute{
						MarkdownDescription: "Backup interval.",
						Computed:            true,
					},
					"retention": schema.Int64Attribute{
						MarkdownDescription: "Backup retention.",
						Computed:            true,
					},
				},
			},
			"authentication": schema.SingleNestedAttribute{
				MarkdownDescription: "Authentication configuration.",
				Computed:            true,
				Attributes: map[string]schema.Attribute{
					"method": schema.StringAttribute{
						MarkdownDescription: "Authentication method.",
						Computed:            true,
					},
					"username": schema.StringAttribute{
						MarkdownDescription: "Username.",
						Computed:            true,
					},
					"password": schema.StringAttribute{
						MarkdownDescription: "Password.",
						Computed:            true,
						Sensitive:           true,
					},
					"encrypted_password": schema.StringAttribute{
						MarkdownDescription: "Needed for validation.",
						Computed:            true,
						Sensitive:           true,
					},
					"required": schema.StringAttribute{
						MarkdownDescription: "Required for everyone or disabled for local addresses.",
						Computed:            true,
					},
				},
			},
			"ssl": schema.SingleNestedAttribute{
				MarkdownDescription: "Backup configuration.",
				Computed:            true,
				Attributes: map[string]schema.Attribute{
					"certificate_validation": schema.StringAttribute{
						MarkdownDescription: "Certificate validation.",
						Computed:            true,
					},
					"cert_path": schema.StringAttribute{
						MarkdownDescription: "Certificate path.",
						Computed:            true,
					},
					"cert_password": schema.StringAttribute{
						MarkdownDescription: "Certificate Password.",
						Computed:            true,
						Sensitive:           true,
					},
					"port": schema.Int64Attribute{
						MarkdownDescription: "SSL port.",
						Computed:            true,
					},
					"enabled": schema.BoolAttribute{
						MarkdownDescription: "Enabled.",
						Computed:            true,
					},
				},
			},
			"proxy": schema.SingleNestedAttribute{
				MarkdownDescription: "Proxy configuration.",
				Computed:            true,
				Attributes: map[string]schema.Attribute{
					"bypass_filter": schema.StringAttribute{
						MarkdownDescription: "Bypass filder.",
						Computed:            true,
					},
					"hostname": schema.StringAttribute{
						MarkdownDescription: "Proxy hostname.",
						Computed:            true,
					},
					"username": schema.StringAttribute{
						MarkdownDescription: "Proxy username.",
						Computed:            true,
					},
					"password": schema.StringAttribute{
						MarkdownDescription: "Proxy password.",
						Computed:            true,
						Sensitive:           true,
					},
					"type": schema.StringAttribute{
						MarkdownDescription: "Proxy type.",
						Computed:            true,
					},
					"port": schema.Int64Attribute{
						MarkdownDescription: "Proxy port.",
						Computed:            true,
					},
					"bypass_local_addresses": schema.BoolAttribute{
						MarkdownDescription: "Bypass for local addresses flag.",
						Computed:            true,
					},
					"enabled": schema.BoolAttribute{
						MarkdownDescription: "Enabled.",
						Computed:            true,
					},
				},
			},
		},
	}
}

func (d *HostDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	if client := helpers.DataSourceConfigure(ctx, req, resp); client != nil {
		d.client = client
	}
}

func (d *HostDataSource) Read(ctx context.Context, _ datasource.ReadRequest, resp *datasource.ReadResponse) {
	var tempDiag diag.Diagnostics
	// assign default empty password value to empty string since it cannot be read
	auth := AuthConfig{
		Password: types.StringValue(""),
	}
	state := Host{}
	state.AuthConfig, tempDiag = types.ObjectValueFrom(ctx, auth.getType().(attr.TypeWithAttributeTypes).AttributeTypes(), auth)
	resp.Diagnostics.Append(tempDiag...)

	// Get host current value
	response, _, err := d.client.HostConfigApi.GetHostConfig(ctx).Execute()
	if err != nil {
		resp.Diagnostics.AddError(helpers.ClientError, helpers.ParseClientError(helpers.Read, hostDataSourceName, err))

		return
	}

	tflog.Trace(ctx, "read "+hostDataSourceName)

	state.write(ctx, response, &resp.Diagnostics)
	resp.Diagnostics.Append(resp.State.Set(ctx, state)...)
}

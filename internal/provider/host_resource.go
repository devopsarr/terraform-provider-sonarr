package provider

import (
	"context"
	"strconv"

	"github.com/devopsarr/sonarr-go/sonarr"
	"github.com/devopsarr/terraform-provider-sonarr/internal/helpers"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/booldefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringdefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

const hostResourceName = "host"

// Ensure provider defined types fully satisfy framework interfaces.
var (
	_ resource.Resource                = &HostResource{}
	_ resource.ResourceWithImportState = &HostResource{}
)

func NewHostResource() resource.Resource {
	return &HostResource{}
}

// HostResource defines the host implementation.
type HostResource struct {
	client *sonarr.APIClient
}

// Host describes the host data model.
type Host struct {
	ProxyConfig    types.Object `tfsdk:"proxy"`
	SSLConfig      types.Object `tfsdk:"ssl"`
	AuthConfig     types.Object `tfsdk:"authentication"`
	BackupConfig   types.Object `tfsdk:"backup"`
	UpdateConfig   types.Object `tfsdk:"update"`
	LoggingConfig  types.Object `tfsdk:"logging"`
	InstanceName   types.String `tfsdk:"instance_name"`
	ApplicationURL types.String `tfsdk:"application_url"`
	BindAddress    types.String `tfsdk:"bind_address"`
	URLBase        types.String `tfsdk:"url_base"`
	ID             types.Int64  `tfsdk:"id"`
	Port           types.Int64  `tfsdk:"port"`
	LaunchBrowser  types.Bool   `tfsdk:"launch_browser"`
}

// ProxyConfig is part of Host.
type ProxyConfig struct {
	Username             types.String `tfsdk:"username"`
	BypassFilter         types.String `tfsdk:"bypass_filter"`
	Password             types.String `tfsdk:"password"`
	Hostname             types.String `tfsdk:"hostname"`
	Type                 types.String `tfsdk:"type"`
	Port                 types.Int64  `tfsdk:"port"`
	Enabled              types.Bool   `tfsdk:"enabled"`
	BypassLocalAddresses types.Bool   `tfsdk:"bypass_local_addresses"`
}

func (p ProxyConfig) getType() attr.Type {
	return types.ObjectType{}.WithAttributeTypes(
		map[string]attr.Type{
			"username":               types.StringType,
			"bypass_filter":          types.StringType,
			"password":               types.StringType,
			"hostname":               types.StringType,
			"type":                   types.StringType,
			"port":                   types.Int64Type,
			"enabled":                types.BoolType,
			"bypass_local_addresses": types.BoolType,
		})
}

// SSLConfig is part of Host.
type SSLConfig struct {
	CertificateValidation types.String `tfsdk:"certificate_validation"`
	CertPath              types.String `tfsdk:"cert_path"`
	CertPassword          types.String `tfsdk:"cert_password"`
	Port                  types.Int64  `tfsdk:"port"`
	Enabled               types.Bool   `tfsdk:"enabled"`
}

func (s SSLConfig) getType() attr.Type {
	return types.ObjectType{}.WithAttributeTypes(
		map[string]attr.Type{
			"certificate_validation": types.StringType,
			"cert_path":              types.StringType,
			"cert_password":          types.StringType,
			"port":                   types.Int64Type,
			"enabled":                types.BoolType,
		})
}

// AuthConfig is part of Host.
type AuthConfig struct {
	Method            types.String `tfsdk:"method"`
	Username          types.String `tfsdk:"username"`
	Password          types.String `tfsdk:"password"`
	EncryptedPassword types.String `tfsdk:"encrypted_password"`
	Required          types.String `tfsdk:"required"`
}

func (a AuthConfig) getType() attr.Type {
	return types.ObjectType{}.WithAttributeTypes(
		map[string]attr.Type{
			"method":             types.StringType,
			"username":           types.StringType,
			"password":           types.StringType,
			"encrypted_password": types.StringType,
			"required":           types.StringType,
		})
}

// BackupConfig is part of Host.
type BackupConfig struct {
	Folder    types.String `tfsdk:"folder"`
	Interval  types.Int64  `tfsdk:"interval"`
	Retention types.Int64  `tfsdk:"retention"`
}

func (b BackupConfig) getType() attr.Type {
	return types.ObjectType{}.WithAttributeTypes(
		map[string]attr.Type{
			"folder":    types.StringType,
			"interval":  types.Int64Type,
			"retention": types.Int64Type,
		})
}

// UpdateConfig is part of Host.
type UpdateConfig struct {
	Branch              types.String `tfsdk:"branch"`
	Mechanism           types.String `tfsdk:"mechanism"`
	ScriptPath          types.String `tfsdk:"script_path"`
	UpdateAutomatically types.Bool   `tfsdk:"update_automatically"`
}

func (u UpdateConfig) getType() attr.Type {
	return types.ObjectType{}.WithAttributeTypes(
		map[string]attr.Type{
			"branch":               types.StringType,
			"mechanism":            types.StringType,
			"script_path":          types.StringType,
			"update_automatically": types.BoolType,
		})
}

// LoggingConfig is part of Host.
type LoggingConfig struct {
	LogLevel         types.String `tfsdk:"log_level"`
	ConsoleLogLevel  types.String `tfsdk:"console_log_level"`
	AnalyticsEnabled types.Bool   `tfsdk:"analytics_enabled"`
}

func (l LoggingConfig) getType() attr.Type {
	return types.ObjectType{}.WithAttributeTypes(
		map[string]attr.Type{
			"log_level":         types.StringType,
			"console_log_level": types.StringType,
			"analytics_enabled": types.BoolType,
		})
}

func (r *HostResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_" + hostResourceName
}

func (r *HostResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "<!-- subcategory:System -->Host resource.\nFor more information refer to [Host](https://wiki.servarr.com/sonarr/settings#general) documentation.",
		Attributes: map[string]schema.Attribute{
			"launch_browser": schema.BoolAttribute{
				MarkdownDescription: "Launch browser flag.",
				Optional:            true,
				Computed:            true,
				Default:             booldefault.StaticBool(true),
			},
			"port": schema.Int64Attribute{
				MarkdownDescription: "TCP port.",
				Required:            true,
			},
			"id": schema.Int64Attribute{
				MarkdownDescription: "Host ID.",
				Computed:            true,
				PlanModifiers: []planmodifier.Int64{
					int64planmodifier.UseStateForUnknown(),
				},
			},
			"url_base": schema.StringAttribute{
				MarkdownDescription: "URL base.",
				Required:            true,
			},
			"bind_address": schema.StringAttribute{
				MarkdownDescription: "Bind address.",
				Required:            true,
			},
			"application_url": schema.StringAttribute{
				MarkdownDescription: "Application URL.",
				Required:            true,
			},
			"instance_name": schema.StringAttribute{
				MarkdownDescription: "Instance name.",
				Required:            true,
			},
			"update": schema.SingleNestedAttribute{
				MarkdownDescription: "Update configuration.",
				Required:            true,
				Attributes: map[string]schema.Attribute{
					"mechanism": schema.StringAttribute{
						MarkdownDescription: "Update mechanism.",
						Required:            true,
					},
					"script_path": schema.StringAttribute{
						MarkdownDescription: "Script path.",
						Optional:            true,
						Computed:            true,
					},
					"branch": schema.StringAttribute{
						MarkdownDescription: "Branch reference.",
						Required:            true,
					},
					"update_automatically": schema.BoolAttribute{
						MarkdownDescription: "Update automatically flag.",
						Optional:            true,
						Computed:            true,
					},
				},
			},
			"logging": schema.SingleNestedAttribute{
				MarkdownDescription: "Logging configuration.",
				Required:            true,
				Attributes: map[string]schema.Attribute{
					"log_level": schema.StringAttribute{
						MarkdownDescription: "Log level.",
						Required:            true,
					},
					"console_log_level": schema.StringAttribute{
						MarkdownDescription: "Console log level.",
						Optional:            true,
						Computed:            true,
					},
					"analytics_enabled": schema.BoolAttribute{
						MarkdownDescription: "Enable analytics flag.",
						Optional:            true,
						Computed:            true,
					},
				},
			},
			"backup": schema.SingleNestedAttribute{
				MarkdownDescription: "Backup configuration.",
				Required:            true,
				Attributes: map[string]schema.Attribute{
					"folder": schema.StringAttribute{
						MarkdownDescription: "Backup folder.",
						Required:            true,
					},
					"interval": schema.Int64Attribute{
						MarkdownDescription: "Backup interval.",
						Required:            true,
					},
					"retention": schema.Int64Attribute{
						MarkdownDescription: "Backup retention.",
						Required:            true,
					},
				},
			},
			"authentication": schema.SingleNestedAttribute{
				MarkdownDescription: "Authentication configuration.",
				Required:            true,
				Attributes: map[string]schema.Attribute{
					"method": schema.StringAttribute{
						MarkdownDescription: "Authentication method.",
						Required:            true,
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
						Default:             stringdefault.StaticString(""),
						PlanModifiers: []planmodifier.String{
							stringplanmodifier.UseStateForUnknown(),
						},
					},
					"encrypted_password": schema.StringAttribute{
						MarkdownDescription: "Needed for validation.",
						Computed:            true,
						Sensitive:           true,
					},
					"required": schema.StringAttribute{
						MarkdownDescription: "Required for everyone or disabled for local addresses.",
						Optional:            true,
						Computed:            true,
						Validators: []validator.String{
							stringvalidator.OneOf("disabledForLocalAddresses", "enabled"),
						},
						Default: stringdefault.StaticString("disabledForLocalAddresses"),
					},
				},
			},
			"ssl": schema.SingleNestedAttribute{
				MarkdownDescription: "Backup configuration.",
				Required:            true,
				Attributes: map[string]schema.Attribute{
					"certificate_validation": schema.StringAttribute{
						MarkdownDescription: "Certificate validation.",
						Required:            true,
					},
					"cert_path": schema.StringAttribute{
						MarkdownDescription: "Certificate path.",
						Optional:            true,
						Computed:            true,
					},
					"cert_password": schema.StringAttribute{
						MarkdownDescription: "Certificate Password.",
						Optional:            true,
						Computed:            true,
						Sensitive:           true,
					},
					"port": schema.Int64Attribute{
						MarkdownDescription: "SSL port.",
						Optional:            true,
						Computed:            true,
					},
					"enabled": schema.BoolAttribute{
						MarkdownDescription: "Enabled.",
						Required:            true,
					},
				},
			},
			"proxy": schema.SingleNestedAttribute{
				MarkdownDescription: "Proxy configuration.",
				Required:            true,
				Attributes: map[string]schema.Attribute{
					"bypass_filter": schema.StringAttribute{
						MarkdownDescription: "Bypass filder.",
						Optional:            true,
						Computed:            true,
					},
					"hostname": schema.StringAttribute{
						MarkdownDescription: "Proxy hostname.",
						Optional:            true,
						Computed:            true,
					},
					"username": schema.StringAttribute{
						MarkdownDescription: "Proxy username.",
						Optional:            true,
						Computed:            true,
					},
					"password": schema.StringAttribute{
						MarkdownDescription: "Proxy password.",
						Optional:            true,
						Computed:            true,
						Sensitive:           true,
					},
					"type": schema.StringAttribute{
						MarkdownDescription: "Proxy type.",
						Optional:            true,
						Computed:            true,
					},
					"port": schema.Int64Attribute{
						MarkdownDescription: "Proxy port.",
						Optional:            true,
						Computed:            true,
					},
					"bypass_local_addresses": schema.BoolAttribute{
						MarkdownDescription: "Bypass for local addresses flag.",
						Optional:            true,
						Computed:            true,
					},
					"enabled": schema.BoolAttribute{
						MarkdownDescription: "Enabled.",
						Required:            true,
					},
				},
			},
		},
	}
}

func (r *HostResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	if client := helpers.ResourceConfigure(ctx, req, resp); client != nil {
		r.client = client
	}
}

func (r *HostResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	// Retrieve values from plan
	var host *Host

	resp.Diagnostics.Append(req.Plan.Get(ctx, &host)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Build Create resource
	request := host.read(ctx, &resp.Diagnostics)
	request.SetId(1)

	// Create new Host
	response, _, err := r.client.HostConfigApi.UpdateHostConfig(ctx, strconv.Itoa(int(request.GetId()))).HostConfigResource(*request).Execute()
	if err != nil {
		resp.Diagnostics.AddError(helpers.ClientError, helpers.ParseClientError(helpers.Create, hostResourceName, err))

		return
	}

	tflog.Trace(ctx, "created "+hostResourceName+": "+strconv.Itoa(int(response.GetId())))
	// Generate resource state struct
	host.write(ctx, response, &resp.Diagnostics)
	resp.Diagnostics.Append(resp.State.Set(ctx, &host)...)
}

func (r *HostResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	// Get current state
	var host *Host

	resp.Diagnostics.Append(req.State.Get(ctx, &host)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Get host current value
	response, _, err := r.client.HostConfigApi.GetHostConfig(ctx).Execute()
	if err != nil {
		resp.Diagnostics.AddError(helpers.ClientError, helpers.ParseClientError(helpers.Read, hostResourceName, err))

		return
	}

	tflog.Trace(ctx, "read "+hostResourceName+": "+strconv.Itoa(int(response.GetId())))
	// Map response body to resource schema attribute
	host.write(ctx, response, &resp.Diagnostics)
	resp.Diagnostics.Append(resp.State.Set(ctx, &host)...)
}

func (r *HostResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	// Get plan values
	var host *Host

	resp.Diagnostics.Append(req.Plan.Get(ctx, &host)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Build Update resource
	request := host.read(ctx, &resp.Diagnostics)

	// Update Host
	response, _, err := r.client.HostConfigApi.UpdateHostConfig(ctx, strconv.Itoa(int(request.GetId()))).HostConfigResource(*request).Execute()
	if err != nil {
		resp.Diagnostics.AddError(helpers.ClientError, helpers.ParseClientError(helpers.Update, hostResourceName, err))

		return
	}

	tflog.Trace(ctx, "updated "+hostResourceName+": "+strconv.Itoa(int(response.GetId())))
	// Generate resource state struct
	host.write(ctx, response, &resp.Diagnostics)
	resp.Diagnostics.Append(resp.State.Set(ctx, &host)...)
}

func (r *HostResource) Delete(ctx context.Context, _ resource.DeleteRequest, resp *resource.DeleteResponse) {
	// Host cannot be really deleted just removing configuration
	tflog.Trace(ctx, "decoupled "+hostResourceName+": 1")
	resp.State.RemoveResource(ctx)
}

func (r *HostResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	tflog.Trace(ctx, "imported "+hostResourceName+": 1")
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("id"), 1)...)
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("authentication").AtName("password"), req.ID)...)
}

func (h *Host) write(ctx context.Context, host *sonarr.HostConfigResource, diags *diag.Diagnostics) {
	var tempDiag diag.Diagnostics

	h.InstanceName = types.StringValue(host.GetInstanceName())
	h.ApplicationURL = types.StringValue(host.GetApplicationUrl())
	h.BindAddress = types.StringValue(host.GetBindAddress())
	h.URLBase = types.StringValue(host.GetUrlBase())
	h.ID = types.Int64Value(int64(host.GetId()))
	h.Port = types.Int64Value(int64(host.GetPort()))
	h.LaunchBrowser = types.BoolValue(host.GetLaunchBrowser())

	proxy := ProxyConfig{}
	ssl := SSLConfig{}
	auth := AuthConfig{}
	backup := BackupConfig{}
	update := UpdateConfig{}
	log := LoggingConfig{}

	// Get the state/plan password to propagate the same
	diags.Append(h.AuthConfig.As(ctx, &auth, basetypes.ObjectAsOptions{})...)

	proxy.write(host)
	ssl.write(host)
	auth.write(host)
	backup.write(host)
	update.write(host)
	log.write(host)

	h.ProxyConfig, tempDiag = types.ObjectValueFrom(ctx, proxy.getType().(attr.TypeWithAttributeTypes).AttributeTypes(), proxy)
	diags.Append(tempDiag...)
	h.SSLConfig, tempDiag = types.ObjectValueFrom(ctx, ssl.getType().(attr.TypeWithAttributeTypes).AttributeTypes(), ssl)
	diags.Append(tempDiag...)
	h.AuthConfig, tempDiag = types.ObjectValueFrom(ctx, auth.getType().(attr.TypeWithAttributeTypes).AttributeTypes(), auth)
	diags.Append(tempDiag...)
	h.BackupConfig, tempDiag = types.ObjectValueFrom(ctx, backup.getType().(attr.TypeWithAttributeTypes).AttributeTypes(), backup)
	diags.Append(tempDiag...)
	h.UpdateConfig, tempDiag = types.ObjectValueFrom(ctx, update.getType().(attr.TypeWithAttributeTypes).AttributeTypes(), update)
	diags.Append(tempDiag...)
	h.LoggingConfig, tempDiag = types.ObjectValueFrom(ctx, log.getType().(attr.TypeWithAttributeTypes).AttributeTypes(), log)
	diags.Append(tempDiag...)
}

func (l *LoggingConfig) write(host *sonarr.HostConfigResource) {
	l.AnalyticsEnabled = types.BoolValue(host.GetAnalyticsEnabled())
	l.ConsoleLogLevel = types.StringValue(host.GetConsoleLogLevel())
	l.LogLevel = types.StringValue(host.GetLogLevel())
}

func (u *UpdateConfig) write(host *sonarr.HostConfigResource) {
	u.Branch = types.StringValue(host.GetBranch())
	u.Mechanism = types.StringValue(string(host.GetUpdateMechanism()))
	u.ScriptPath = types.StringValue(host.GetUpdateScriptPath())
	u.UpdateAutomatically = types.BoolValue(host.GetUpdateAutomatically())
}

func (b *BackupConfig) write(host *sonarr.HostConfigResource) {
	b.Folder = types.StringValue(host.GetBackupFolder())
	b.Interval = types.Int64Value(int64(host.GetBackupInterval()))
	b.Retention = types.Int64Value(int64(host.GetBackupRetention()))
}

func (a *AuthConfig) write(host *sonarr.HostConfigResource) {
	a.Method = types.StringValue(string(host.GetAuthenticationMethod()))
	a.Username = types.StringValue(host.GetUsername())
	a.EncryptedPassword = types.StringValue(host.GetPassword())
	a.Required = types.StringValue(string(host.GetAuthenticationRequired()))
}

func (s *SSLConfig) write(host *sonarr.HostConfigResource) {
	s.CertificateValidation = types.StringValue(string(host.GetCertificateValidation()))
	s.CertPath = types.StringValue(host.GetSslCertPath())
	s.CertPassword = types.StringValue(host.GetSslCertPassword())
	s.Port = types.Int64Value(int64(host.GetSslPort()))
	s.Enabled = types.BoolValue(host.GetEnableSsl())
}

func (p *ProxyConfig) write(host *sonarr.HostConfigResource) {
	p.Username = types.StringValue(host.GetProxyUsername())
	p.Password = types.StringValue(host.GetProxyPassword())
	p.BypassFilter = types.StringValue(host.GetProxyBypassFilter())
	p.Hostname = types.StringValue(host.GetProxyHostname())
	p.Type = types.StringValue(string(host.GetProxyType()))
	p.Port = types.Int64Value(int64(host.GetProxyPort()))
	p.BypassLocalAddresses = types.BoolValue(host.GetProxyBypassLocalAddresses())
	p.Enabled = types.BoolValue(host.GetProxyEnabled())
}

func (h *Host) read(ctx context.Context, diags *diag.Diagnostics) *sonarr.HostConfigResource {
	host := sonarr.NewHostConfigResource()
	host.SetInstanceName(h.InstanceName.ValueString())
	host.SetApplicationUrl(h.ApplicationURL.ValueString())
	host.SetBindAddress(h.BindAddress.ValueString())
	host.SetUrlBase(h.URLBase.ValueString())
	host.SetId(int32(h.ID.ValueInt64()))
	host.SetPort(int32(h.Port.ValueInt64()))
	host.SetLaunchBrowser(h.LaunchBrowser.ValueBool())

	proxy := ProxyConfig{}
	ssl := SSLConfig{}
	auth := AuthConfig{}
	backup := BackupConfig{}
	update := UpdateConfig{}
	log := LoggingConfig{}

	diags.Append(h.ProxyConfig.As(ctx, &proxy, basetypes.ObjectAsOptions{})...)
	proxy.read(host)

	diags.Append(h.SSLConfig.As(ctx, &ssl, basetypes.ObjectAsOptions{})...)
	ssl.read(host)

	diags.Append(h.AuthConfig.As(ctx, &auth, basetypes.ObjectAsOptions{})...)
	auth.read(host)

	diags.Append(h.BackupConfig.As(ctx, &backup, basetypes.ObjectAsOptions{})...)
	backup.read(host)

	diags.Append(h.UpdateConfig.As(ctx, &update, basetypes.ObjectAsOptions{})...)
	update.read(host)

	diags.Append(h.LoggingConfig.As(ctx, &log, basetypes.ObjectAsOptions{})...)
	log.read(host)

	return host
}

func (l *LoggingConfig) read(host *sonarr.HostConfigResource) {
	host.SetAnalyticsEnabled(l.AnalyticsEnabled.ValueBool())
	host.SetConsoleLogLevel(l.LogLevel.ValueString())
	host.SetLogLevel(l.LogLevel.ValueString())
}

func (u *UpdateConfig) read(host *sonarr.HostConfigResource) {
	host.SetBranch(u.Branch.ValueString())
	host.SetUpdateMechanism(sonarr.UpdateMechanism(u.Mechanism.ValueString()))
	host.SetUpdateScriptPath(u.ScriptPath.ValueString())
	host.SetUpdateAutomatically(u.UpdateAutomatically.ValueBool())
}

func (b *BackupConfig) read(host *sonarr.HostConfigResource) {
	host.SetBackupFolder(b.Folder.ValueString())
	host.SetBackupInterval(int32(b.Interval.ValueInt64()))
	host.SetBackupRetention(int32(b.Retention.ValueInt64()))
}

func (a *AuthConfig) read(host *sonarr.HostConfigResource) {
	host.SetAuthenticationMethod(sonarr.AuthenticationType(a.Method.ValueString()))
	host.SetUsername(a.Username.ValueString())
	host.SetAuthenticationRequired(sonarr.AuthenticationRequiredType(a.Required.ValueString()))
	host.SetPassword(a.Password.ValueString())
}

func (s *SSLConfig) read(host *sonarr.HostConfigResource) {
	host.SetCertificateValidation(sonarr.CertificateValidationType(s.CertificateValidation.ValueString()))
	host.SetSslCertPath(s.CertPath.ValueString())
	host.SetSslCertPassword(s.CertPassword.ValueString())
	host.SetSslPort(int32(s.Port.ValueInt64()))
	host.SetEnableSsl(s.Enabled.ValueBool())
}

func (p *ProxyConfig) read(host *sonarr.HostConfigResource) {
	host.SetProxyUsername(p.Username.ValueString())
	host.SetProxyPassword(p.Password.ValueString())
	host.SetProxyBypassFilter(p.BypassFilter.ValueString())
	host.SetProxyHostname(p.Hostname.ValueString())
	host.SetProxyPort(int32(p.Port.ValueInt64()))
	host.SetProxyEnabled(p.Enabled.ValueBool())
	host.SetProxyBypassLocalAddresses(p.BypassLocalAddresses.ValueBool())
}

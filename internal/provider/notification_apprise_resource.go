package provider

import (
	"context"
	"strconv"

	"github.com/devopsarr/sonarr-go/sonarr"
	"github.com/devopsarr/terraform-provider-sonarr/internal/helpers"
	"github.com/hashicorp/terraform-plugin-framework-validators/int64validator"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

const (
	notificationAppriseResourceName   = "notification_apprise"
	notificationAppriseImplementation = "Apprise"
	notificationAppriseConfigContract = "AppriseSettings"
)

// Ensure provider defined types fully satisfy framework interfaces.
var (
	_ resource.Resource                = &NotificationAppriseResource{}
	_ resource.ResourceWithImportState = &NotificationAppriseResource{}
)

func NewNotificationAppriseResource() resource.Resource {
	return &NotificationAppriseResource{}
}

// NotificationAppriseResource defines the notification implementation.
type NotificationAppriseResource struct {
	client *sonarr.APIClient
	auth   context.Context
}

// NotificationApprise describes the notification data model.
type NotificationApprise struct {
	Tags                          types.Set    `tfsdk:"tags"`
	FieldTags                     types.Set    `tfsdk:"field_tags"`
	Name                          types.String `tfsdk:"name"`
	StatelessURLs                 types.String `tfsdk:"stateless_urls"`
	ServerURL                     types.String `tfsdk:"server_url"`
	AuthUsername                  types.String `tfsdk:"auth_username"`
	AuthPassword                  types.String `tfsdk:"auth_password"`
	ConfigurationKey              types.String `tfsdk:"configuration_key"`
	NotificationType              types.Int64  `tfsdk:"notification_type"`
	ID                            types.Int64  `tfsdk:"id"`
	OnGrab                        types.Bool   `tfsdk:"on_grab"`
	OnEpisodeFileDeleteForUpgrade types.Bool   `tfsdk:"on_episode_file_delete_for_upgrade"`
	OnEpisodeFileDelete           types.Bool   `tfsdk:"on_episode_file_delete"`
	IncludeHealthWarnings         types.Bool   `tfsdk:"include_health_warnings"`
	OnApplicationUpdate           types.Bool   `tfsdk:"on_application_update"`
	OnHealthIssue                 types.Bool   `tfsdk:"on_health_issue"`
	OnHealthRestored              types.Bool   `tfsdk:"on_health_restored"`
	OnManualInteractionRequired   types.Bool   `tfsdk:"on_manual_interaction_required"`
	OnSeriesAdd                   types.Bool   `tfsdk:"on_series_add"`
	OnSeriesDelete                types.Bool   `tfsdk:"on_series_delete"`
	OnUpgrade                     types.Bool   `tfsdk:"on_upgrade"`
	OnDownload                    types.Bool   `tfsdk:"on_download"`
}

func (n NotificationApprise) toNotification() *Notification {
	return &Notification{
		Tags:                          n.Tags,
		FieldTags:                     n.FieldTags,
		StatelessURLs:                 n.StatelessURLs,
		ServerURL:                     n.ServerURL,
		AuthUsername:                  n.AuthUsername,
		AuthPassword:                  n.AuthPassword,
		ConfigurationKey:              n.ConfigurationKey,
		NotificationType:              n.NotificationType,
		Name:                          n.Name,
		ID:                            n.ID,
		OnGrab:                        n.OnGrab,
		OnEpisodeFileDeleteForUpgrade: n.OnEpisodeFileDeleteForUpgrade,
		OnEpisodeFileDelete:           n.OnEpisodeFileDelete,
		IncludeHealthWarnings:         n.IncludeHealthWarnings,
		OnApplicationUpdate:           n.OnApplicationUpdate,
		OnHealthIssue:                 n.OnHealthIssue,
		OnHealthRestored:              n.OnHealthRestored,
		OnManualInteractionRequired:   n.OnManualInteractionRequired,
		OnSeriesAdd:                   n.OnSeriesAdd,
		OnSeriesDelete:                n.OnSeriesDelete,
		OnUpgrade:                     n.OnUpgrade,
		OnDownload:                    n.OnDownload,
		ConfigContract:                types.StringValue(notificationAppriseConfigContract),
		Implementation:                types.StringValue(notificationAppriseImplementation),
	}
}

func (n *NotificationApprise) fromNotification(notification *Notification) {
	n.Tags = notification.Tags
	n.FieldTags = notification.FieldTags
	n.StatelessURLs = notification.StatelessURLs
	n.ServerURL = notification.ServerURL
	n.AuthUsername = notification.AuthUsername
	n.AuthPassword = notification.AuthPassword
	n.ConfigurationKey = notification.ConfigurationKey
	n.NotificationType = notification.NotificationType
	n.Name = notification.Name
	n.ID = notification.ID
	n.OnGrab = notification.OnGrab
	n.OnEpisodeFileDeleteForUpgrade = notification.OnEpisodeFileDeleteForUpgrade
	n.OnEpisodeFileDelete = notification.OnEpisodeFileDelete
	n.IncludeHealthWarnings = notification.IncludeHealthWarnings
	n.OnApplicationUpdate = notification.OnApplicationUpdate
	n.OnHealthIssue = notification.OnHealthIssue
	n.OnHealthRestored = notification.OnHealthRestored
	n.OnManualInteractionRequired = notification.OnManualInteractionRequired
	n.OnSeriesAdd = notification.OnSeriesAdd
	n.OnSeriesDelete = notification.OnSeriesDelete
	n.OnUpgrade = notification.OnUpgrade
	n.OnDownload = notification.OnDownload
}

func (r *NotificationAppriseResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_" + notificationAppriseResourceName
}

func (r *NotificationAppriseResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "<!-- subcategory:Notifications -->Notification Apprise resource.\nFor more information refer to [Notification](https://wiki.servarr.com/sonarr/settings#connect) and [Apprise](https://wiki.servarr.com/sonarr/supported#apprise).",
		Attributes: map[string]schema.Attribute{
			"on_grab": schema.BoolAttribute{
				MarkdownDescription: "On grab flag.",
				Optional:            true,
				Computed:            true,
			},
			"on_download": schema.BoolAttribute{
				MarkdownDescription: "On download flag.",
				Optional:            true,
				Computed:            true,
			},
			"on_upgrade": schema.BoolAttribute{
				MarkdownDescription: "On upgrade flag.",
				Optional:            true,
				Computed:            true,
			},
			"on_series_add": schema.BoolAttribute{
				MarkdownDescription: "On series add flag.",
				Optional:            true,
				Computed:            true,
			},
			"on_series_delete": schema.BoolAttribute{
				MarkdownDescription: "On series delete flag.",
				Optional:            true,
				Computed:            true,
			},
			"on_episode_file_delete": schema.BoolAttribute{
				MarkdownDescription: "On episode file delete flag.",
				Optional:            true,
				Computed:            true,
			},
			"on_episode_file_delete_for_upgrade": schema.BoolAttribute{
				MarkdownDescription: "On episode file delete for upgrade flag.",
				Optional:            true,
				Computed:            true,
			},
			"on_health_issue": schema.BoolAttribute{
				MarkdownDescription: "On health issue flag.",
				Optional:            true,
				Computed:            true,
			},
			"on_health_restored": schema.BoolAttribute{
				MarkdownDescription: "On health restored flag.",
				Optional:            true,
				Computed:            true,
			},
			"on_manual_interaction_required": schema.BoolAttribute{
				MarkdownDescription: "On manual interaction required flag.",
				Optional:            true,
				Computed:            true,
			},
			"on_application_update": schema.BoolAttribute{
				MarkdownDescription: "On application update flag.",
				Optional:            true,
				Computed:            true,
			},
			"include_health_warnings": schema.BoolAttribute{
				MarkdownDescription: "Include health warnings.",
				Optional:            true,
				Computed:            true,
			},
			"name": schema.StringAttribute{
				MarkdownDescription: "NotificationApprise name.",
				Required:            true,
			},
			"tags": schema.SetAttribute{
				MarkdownDescription: "List of associated tags.",
				Optional:            true,
				Computed:            true,
				ElementType:         types.Int64Type,
			},
			"id": schema.Int64Attribute{
				MarkdownDescription: "Notification ID.",
				Computed:            true,
				PlanModifiers: []planmodifier.Int64{
					int64planmodifier.UseStateForUnknown(),
				},
			},
			// Field values
			"notification_type": schema.Int64Attribute{
				MarkdownDescription: "Notification type. `0` Info, `1` Success, `2` Warning, `3` Failure.",
				Optional:            true,
				Computed:            true,
				Validators: []validator.Int64{
					int64validator.OneOf(0, 1, 2, 3),
				},
			},
			"server_url": schema.StringAttribute{
				MarkdownDescription: "Server URL.",
				Required:            true,
			},
			"stateless_urls": schema.StringAttribute{
				MarkdownDescription: "Stateless URLs.",
				Optional:            true,
				Computed:            true,
			},
			"configuration_key": schema.StringAttribute{
				MarkdownDescription: "Configuration key.",
				Optional:            true,
				Computed:            true,
				Sensitive:           true,
			},
			"auth_username": schema.StringAttribute{
				MarkdownDescription: "Username.",
				Optional:            true,
				Computed:            true,
			},
			"auth_password": schema.StringAttribute{
				MarkdownDescription: "Password.",
				Optional:            true,
				Computed:            true,
				Sensitive:           true,
			},
			"field_tags": schema.SetAttribute{
				MarkdownDescription: "Tags and emojis.",
				Optional:            true,
				Computed:            true,
				ElementType:         types.StringType,
			},
		},
	}
}

func (r *NotificationAppriseResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	if auth, client := resourceConfigure(ctx, req, resp); client != nil {
		r.client = client
		r.auth = auth
	}
}

func (r *NotificationAppriseResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	// Retrieve values from plan
	var notification *NotificationApprise

	resp.Diagnostics.Append(req.Plan.Get(ctx, &notification)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Create new NotificationApprise
	request := notification.read(ctx, &resp.Diagnostics)

	response, _, err := r.client.NotificationAPI.CreateNotification(r.auth).NotificationResource(*request).Execute()
	if err != nil {
		resp.Diagnostics.AddError(helpers.ClientError, helpers.ParseClientError(helpers.Create, notificationAppriseResourceName, err))

		return
	}

	tflog.Trace(ctx, "created "+notificationAppriseResourceName+": "+strconv.Itoa(int(response.GetId())))
	// Generate resource state struct
	notification.write(ctx, response, &resp.Diagnostics)
	resp.Diagnostics.Append(resp.State.Set(ctx, &notification)...)
}

func (r *NotificationAppriseResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	// Get current state
	var notification *NotificationApprise

	resp.Diagnostics.Append(req.State.Get(ctx, &notification)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Get NotificationApprise current value
	response, _, err := r.client.NotificationAPI.GetNotificationById(r.auth, int32(notification.ID.ValueInt64())).Execute()
	if err != nil {
		resp.Diagnostics.AddError(helpers.ClientError, helpers.ParseClientError(helpers.Read, notificationAppriseResourceName, err))

		return
	}

	tflog.Trace(ctx, "read "+notificationAppriseResourceName+": "+strconv.Itoa(int(response.GetId())))
	// Map response body to resource schema attribute
	notification.write(ctx, response, &resp.Diagnostics)
	resp.Diagnostics.Append(resp.State.Set(ctx, &notification)...)
}

func (r *NotificationAppriseResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	// Get plan values
	var notification *NotificationApprise

	resp.Diagnostics.Append(req.Plan.Get(ctx, &notification)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Update NotificationApprise
	request := notification.read(ctx, &resp.Diagnostics)

	response, _, err := r.client.NotificationAPI.UpdateNotification(r.auth, strconv.Itoa(int(request.GetId()))).NotificationResource(*request).Execute()
	if err != nil {
		resp.Diagnostics.AddError(helpers.ClientError, helpers.ParseClientError(helpers.Update, notificationAppriseResourceName, err))

		return
	}

	tflog.Trace(ctx, "updated "+notificationAppriseResourceName+": "+strconv.Itoa(int(response.GetId())))
	// Generate resource state struct
	notification.write(ctx, response, &resp.Diagnostics)
	resp.Diagnostics.Append(resp.State.Set(ctx, &notification)...)
}

func (r *NotificationAppriseResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var ID int64

	resp.Diagnostics.Append(req.State.GetAttribute(ctx, path.Root("id"), &ID)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Delete NotificationApprise current value
	_, err := r.client.NotificationAPI.DeleteNotification(r.auth, int32(ID)).Execute()
	if err != nil {
		resp.Diagnostics.AddError(helpers.ClientError, helpers.ParseClientError(helpers.Delete, notificationAppriseResourceName, err))

		return
	}

	tflog.Trace(ctx, "deleted "+notificationAppriseResourceName+": "+strconv.Itoa(int(ID)))
	resp.State.RemoveResource(ctx)
}

func (r *NotificationAppriseResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	helpers.ImportStatePassthroughIntID(ctx, path.Root("id"), req, resp)
	tflog.Trace(ctx, "imported "+notificationAppriseResourceName+": "+req.ID)
}

func (n *NotificationApprise) write(ctx context.Context, notification *sonarr.NotificationResource, diags *diag.Diagnostics) {
	genericNotification := n.toNotification()
	genericNotification.write(ctx, notification, diags)
	n.fromNotification(genericNotification)
}

func (n *NotificationApprise) read(ctx context.Context, diags *diag.Diagnostics) *sonarr.NotificationResource {
	return n.toNotification().read(ctx, diags)
}

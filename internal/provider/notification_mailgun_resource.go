package provider

import (
	"context"
	"strconv"

	"github.com/devopsarr/sonarr-go/sonarr"
	"github.com/devopsarr/terraform-provider-sonarr/internal/helpers"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

const (
	notificationMailgunResourceName   = "notification_mailgun"
	notificationMailgunImplementation = "Mailgun"
	notificationMailgunConfigContract = "MailgunSettings"
)

// Ensure provider defined types fully satisfy framework interfaces.
var (
	_ resource.Resource                = &NotificationMailgunResource{}
	_ resource.ResourceWithImportState = &NotificationMailgunResource{}
)

func NewNotificationMailgunResource() resource.Resource {
	return &NotificationMailgunResource{}
}

// NotificationMailgunResource defines the notification implementation.
type NotificationMailgunResource struct {
	client *sonarr.APIClient
	auth   context.Context
}

// NotificationMailgun describes the notification data model.
type NotificationMailgun struct {
	Tags                          types.Set    `tfsdk:"tags"`
	Recipients                    types.Set    `tfsdk:"recipients"`
	From                          types.String `tfsdk:"from"`
	SenderDomain                  types.String `tfsdk:"sender_domain"`
	Name                          types.String `tfsdk:"name"`
	APIKey                        types.String `tfsdk:"api_key"`
	ID                            types.Int64  `tfsdk:"id"`
	UseEuEndpoint                 types.Bool   `tfsdk:"use_eu_endpoint"`
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

func (n NotificationMailgun) toNotification() *Notification {
	return &Notification{
		Tags:                          n.Tags,
		Recipients:                    n.Recipients,
		SenderDomain:                  n.SenderDomain,
		APIKey:                        n.APIKey,
		UseEuEndpoint:                 n.UseEuEndpoint,
		Name:                          n.Name,
		From:                          n.From,
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
		ConfigContract:                types.StringValue(notificationMailgunConfigContract),
		Implementation:                types.StringValue(notificationMailgunImplementation),
	}
}

func (n *NotificationMailgun) fromNotification(notification *Notification) {
	n.Tags = notification.Tags
	n.Recipients = notification.Recipients
	n.SenderDomain = notification.SenderDomain
	n.APIKey = notification.APIKey
	n.UseEuEndpoint = notification.UseEuEndpoint
	n.Name = notification.Name
	n.From = notification.From
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

func (r *NotificationMailgunResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_" + notificationMailgunResourceName
}

func (r *NotificationMailgunResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "<!-- subcategory:Notifications -->\nNotification Mailgun resource.\nFor more information refer to [Notification](https://wiki.servarr.com/sonarr/settings#connect) and [Mailgun](https://wiki.servarr.com/sonarr/supported#mailgun).",
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
				MarkdownDescription: "NotificationMailgun name.",
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
			"use_eu_endpoint": schema.BoolAttribute{
				MarkdownDescription: "Use EU endpoint flag.",
				Optional:            true,
				Computed:            true,
			},
			"api_key": schema.StringAttribute{
				MarkdownDescription: "API key.",
				Optional:            true,
				Computed:            true,
				Sensitive:           true,
			},
			"from": schema.StringAttribute{
				MarkdownDescription: "From.",
				Required:            true,
			},
			"sender_domain": schema.StringAttribute{
				MarkdownDescription: "Sender domain.",
				Optional:            true,
				Computed:            true,
			},
			"recipients": schema.SetAttribute{
				MarkdownDescription: "Recipients.",
				Required:            true,
				ElementType:         types.StringType,
			},
		},
	}
}

func (r *NotificationMailgunResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	if auth, client := resourceConfigure(ctx, req, resp); client != nil {
		r.client = client
		r.auth = auth
	}
}

func (r *NotificationMailgunResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	// Retrieve values from plan
	var notification *NotificationMailgun

	resp.Diagnostics.Append(req.Plan.Get(ctx, &notification)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Create new NotificationMailgun
	request := notification.read(ctx, &resp.Diagnostics)

	response, _, err := r.client.NotificationAPI.CreateNotification(r.auth).NotificationResource(*request).Execute()
	if err != nil {
		resp.Diagnostics.AddError(helpers.ClientError, helpers.ParseClientError(helpers.Create, notificationMailgunResourceName, err))

		return
	}

	tflog.Trace(ctx, "created "+notificationMailgunResourceName+": "+strconv.Itoa(int(response.GetId())))
	// Generate resource state struct
	notification.write(ctx, response, &resp.Diagnostics)
	resp.Diagnostics.Append(resp.State.Set(ctx, &notification)...)
}

func (r *NotificationMailgunResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	// Get current state
	var notification *NotificationMailgun

	resp.Diagnostics.Append(req.State.Get(ctx, &notification)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Get NotificationMailgun current value
	response, _, err := r.client.NotificationAPI.GetNotificationById(r.auth, int32(notification.ID.ValueInt64())).Execute()
	if err != nil {
		resp.Diagnostics.AddError(helpers.ClientError, helpers.ParseClientError(helpers.Read, notificationMailgunResourceName, err))

		return
	}

	tflog.Trace(ctx, "read "+notificationMailgunResourceName+": "+strconv.Itoa(int(response.GetId())))
	// Map response body to resource schema attribute
	notification.write(ctx, response, &resp.Diagnostics)
	resp.Diagnostics.Append(resp.State.Set(ctx, &notification)...)
}

func (r *NotificationMailgunResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	// Get plan values
	var notification *NotificationMailgun

	resp.Diagnostics.Append(req.Plan.Get(ctx, &notification)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Update NotificationMailgun
	request := notification.read(ctx, &resp.Diagnostics)

	response, _, err := r.client.NotificationAPI.UpdateNotification(r.auth, strconv.Itoa(int(request.GetId()))).NotificationResource(*request).Execute()
	if err != nil {
		resp.Diagnostics.AddError(helpers.ClientError, helpers.ParseClientError(helpers.Update, notificationMailgunResourceName, err))

		return
	}

	tflog.Trace(ctx, "updated "+notificationMailgunResourceName+": "+strconv.Itoa(int(response.GetId())))
	// Generate resource state struct
	notification.write(ctx, response, &resp.Diagnostics)
	resp.Diagnostics.Append(resp.State.Set(ctx, &notification)...)
}

func (r *NotificationMailgunResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var ID int64

	resp.Diagnostics.Append(req.State.GetAttribute(ctx, path.Root("id"), &ID)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Delete NotificationMailgun current value
	_, err := r.client.NotificationAPI.DeleteNotification(r.auth, int32(ID)).Execute()
	if err != nil {
		resp.Diagnostics.AddError(helpers.ClientError, helpers.ParseClientError(helpers.Delete, notificationMailgunResourceName, err))

		return
	}

	tflog.Trace(ctx, "deleted "+notificationMailgunResourceName+": "+strconv.Itoa(int(ID)))
	resp.State.RemoveResource(ctx)
}

func (r *NotificationMailgunResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	helpers.ImportStatePassthroughIntID(ctx, path.Root("id"), req, resp)
	tflog.Trace(ctx, "imported "+notificationMailgunResourceName+": "+req.ID)
}

func (n *NotificationMailgun) write(ctx context.Context, notification *sonarr.NotificationResource, diags *diag.Diagnostics) {
	genericNotification := n.toNotification()
	genericNotification.write(ctx, notification, diags)
	n.fromNotification(genericNotification)
}

func (n *NotificationMailgun) read(ctx context.Context, diags *diag.Diagnostics) *sonarr.NotificationResource {
	return n.toNotification().read(ctx, diags)
}

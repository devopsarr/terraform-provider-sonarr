package provider

import (
	"context"
	"fmt"
	"strconv"

	"github.com/devopsarr/terraform-provider-sonarr/tools"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"golift.io/starr/sonarr"
)

const (
	notificationSlackResourceName   = "notification_slack"
	notificationSlackImplementation = "Slack"
	notificationSlackConfigContract = "SlackSettings"
)

// Ensure provider defined types fully satisfy framework interfaces.
var (
	_ resource.Resource                = &NotificationSlackResource{}
	_ resource.ResourceWithImportState = &NotificationSlackResource{}
)

func NewNotificationSlackResource() resource.Resource {
	return &NotificationSlackResource{}
}

// NotificationSlackResource defines the notification implementation.
type NotificationSlackResource struct {
	client *sonarr.Sonarr
}

// NotificationSlack describes the notification data model.
type NotificationSlack struct {
	Tags                          types.Set    `tfsdk:"tags"`
	WebHookURL                    types.String `tfsdk:"web_hook_url"`
	Name                          types.String `tfsdk:"name"`
	Username                      types.String `tfsdk:"username"`
	Icon                          types.String `tfsdk:"icon"`
	Channel                       types.String `tfsdk:"channel"`
	ID                            types.Int64  `tfsdk:"id"`
	OnGrab                        types.Bool   `tfsdk:"on_grab"`
	OnEpisodeFileDeleteForUpgrade types.Bool   `tfsdk:"on_episode_file_delete_for_upgrade"`
	OnEpisodeFileDelete           types.Bool   `tfsdk:"on_episode_file_delete"`
	IncludeHealthWarnings         types.Bool   `tfsdk:"include_health_warnings"`
	OnApplicationUpdate           types.Bool   `tfsdk:"on_application_update"`
	OnHealthIssue                 types.Bool   `tfsdk:"on_health_issue"`
	OnSeriesDelete                types.Bool   `tfsdk:"on_series_delete"`
	OnRename                      types.Bool   `tfsdk:"on_rename"`
	OnUpgrade                     types.Bool   `tfsdk:"on_upgrade"`
	OnDownload                    types.Bool   `tfsdk:"on_download"`
}

func (n NotificationSlack) toNotification() *Notification {
	return &Notification{
		Tags:                          n.Tags,
		WebHookURL:                    n.WebHookURL,
		Icon:                          n.Icon,
		Username:                      n.Username,
		Channel:                       n.Channel,
		Name:                          n.Name,
		ID:                            n.ID,
		OnGrab:                        n.OnGrab,
		OnEpisodeFileDeleteForUpgrade: n.OnEpisodeFileDeleteForUpgrade,
		OnEpisodeFileDelete:           n.OnEpisodeFileDelete,
		IncludeHealthWarnings:         n.IncludeHealthWarnings,
		OnApplicationUpdate:           n.OnApplicationUpdate,
		OnHealthIssue:                 n.OnHealthIssue,
		OnSeriesDelete:                n.OnSeriesDelete,
		OnRename:                      n.OnRename,
		OnUpgrade:                     n.OnUpgrade,
		OnDownload:                    n.OnDownload,
	}
}

func (n *NotificationSlack) fromNotification(notification *Notification) {
	n.Tags = notification.Tags
	n.WebHookURL = notification.WebHookURL
	n.Icon = notification.Icon
	n.Username = notification.Username
	n.Channel = notification.Channel
	n.Name = notification.Name
	n.ID = notification.ID
	n.OnGrab = notification.OnGrab
	n.OnEpisodeFileDeleteForUpgrade = notification.OnEpisodeFileDeleteForUpgrade
	n.OnEpisodeFileDelete = notification.OnEpisodeFileDelete
	n.IncludeHealthWarnings = notification.IncludeHealthWarnings
	n.OnApplicationUpdate = notification.OnApplicationUpdate
	n.OnHealthIssue = notification.OnHealthIssue
	n.OnSeriesDelete = notification.OnSeriesDelete
	n.OnRename = notification.OnRename
	n.OnUpgrade = notification.OnUpgrade
	n.OnDownload = notification.OnDownload
}

func (r *NotificationSlackResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_" + notificationSlackResourceName
}

func (r *NotificationSlackResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "<!-- subcategory:Notifications -->Notification Slack resource.\nFor more information refer to [Notification](https://wiki.servarr.com/sonarr/settings#connect) and [Slack](https://wiki.servarr.com/sonarr/supported#slack).",
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
			"on_rename": schema.BoolAttribute{
				MarkdownDescription: "On rename flag.",
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
				MarkdownDescription: "NotificationSlack name.",
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
			"web_hook_url": schema.StringAttribute{
				MarkdownDescription: "URL.",
				Required:            true,
			},
			"username": schema.StringAttribute{
				MarkdownDescription: "Username.",
				Required:            true,
			},
			"icon": schema.StringAttribute{
				MarkdownDescription: "Icon.",
				Optional:            true,
				Computed:            true,
			},
			"channel": schema.StringAttribute{
				MarkdownDescription: "Channel.",
				Optional:            true,
				Computed:            true,
			},
		},
	}
}

func (r *NotificationSlackResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	// Prevent panic if the provider has not been configured.
	if req.ProviderData == nil {
		return
	}

	client, ok := req.ProviderData.(*sonarr.Sonarr)
	if !ok {
		resp.Diagnostics.AddError(
			tools.UnexpectedResourceConfigureType,
			fmt.Sprintf("Expected *sonarr.Sonarr, got: %T. Please report this issue to the provider developers.", req.ProviderData),
		)

		return
	}

	r.client = client
}

func (r *NotificationSlackResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	// Retrieve values from plan
	var notification *NotificationSlack

	resp.Diagnostics.Append(req.Plan.Get(ctx, &notification)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Create new NotificationSlack
	request := notification.read(ctx)

	response, err := r.client.AddNotificationContext(ctx, request)
	if err != nil {
		resp.Diagnostics.AddError(tools.ClientError, fmt.Sprintf("Unable to create %s, got error: %s", notificationSlackResourceName, err))

		return
	}

	tflog.Trace(ctx, "created "+notificationSlackResourceName+": "+strconv.Itoa(int(response.ID)))
	// Generate resource state struct
	notification.write(ctx, response)
	resp.Diagnostics.Append(resp.State.Set(ctx, &notification)...)
}

func (r *NotificationSlackResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	// Get current state
	var notification *NotificationSlack

	resp.Diagnostics.Append(req.State.Get(ctx, &notification)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Get NotificationSlack current value
	response, err := r.client.GetNotificationContext(ctx, int(notification.ID.ValueInt64()))
	if err != nil {
		resp.Diagnostics.AddError(tools.ClientError, fmt.Sprintf("Unable to read %s, got error: %s", notificationSlackResourceName, err))

		return
	}

	tflog.Trace(ctx, "read "+notificationSlackResourceName+": "+strconv.Itoa(int(response.ID)))
	// Map response body to resource schema attribute
	notification.write(ctx, response)
	resp.Diagnostics.Append(resp.State.Set(ctx, &notification)...)
}

func (r *NotificationSlackResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	// Get plan values
	var notification *NotificationSlack

	resp.Diagnostics.Append(req.Plan.Get(ctx, &notification)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Update NotificationSlack
	request := notification.read(ctx)

	response, err := r.client.UpdateNotificationContext(ctx, request)
	if err != nil {
		resp.Diagnostics.AddError(tools.ClientError, fmt.Sprintf("Unable to update %s, got error: %s", notificationSlackResourceName, err))

		return
	}

	tflog.Trace(ctx, "updated "+notificationSlackResourceName+": "+strconv.Itoa(int(response.ID)))
	// Generate resource state struct
	notification.write(ctx, response)
	resp.Diagnostics.Append(resp.State.Set(ctx, &notification)...)
}

func (r *NotificationSlackResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var notification *NotificationSlack

	resp.Diagnostics.Append(req.State.Get(ctx, &notification)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Delete NotificationSlack current value
	err := r.client.DeleteNotificationContext(ctx, notification.ID.ValueInt64())
	if err != nil {
		resp.Diagnostics.AddError(tools.ClientError, fmt.Sprintf("Unable to read %s, got error: %s", notificationSlackResourceName, err))

		return
	}

	tflog.Trace(ctx, "deleted "+notificationSlackResourceName+": "+strconv.Itoa(int(notification.ID.ValueInt64())))
	resp.State.RemoveResource(ctx)
}

func (r *NotificationSlackResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	// resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
	id, err := strconv.Atoi(req.ID)
	if err != nil {
		resp.Diagnostics.AddError(
			tools.UnexpectedImportIdentifier,
			fmt.Sprintf("Expected import identifier with format: ID. Got: %q", req.ID),
		)

		return
	}

	tflog.Trace(ctx, "imported "+notificationSlackResourceName+": "+strconv.Itoa(id))
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("id"), id)...)
}

func (n *NotificationSlack) write(ctx context.Context, notification *sonarr.NotificationOutput) {
	genericNotification := Notification{
		OnGrab:                        types.BoolValue(notification.OnGrab),
		OnDownload:                    types.BoolValue(notification.OnDownload),
		OnUpgrade:                     types.BoolValue(notification.OnUpgrade),
		OnRename:                      types.BoolValue(notification.OnRename),
		OnSeriesDelete:                types.BoolValue(notification.OnSeriesDelete),
		OnEpisodeFileDelete:           types.BoolValue(notification.OnEpisodeFileDelete),
		OnEpisodeFileDeleteForUpgrade: types.BoolValue(notification.OnEpisodeFileDeleteForUpgrade),
		OnHealthIssue:                 types.BoolValue(notification.OnHealthIssue),
		OnApplicationUpdate:           types.BoolValue(notification.OnApplicationUpdate),
		IncludeHealthWarnings:         types.BoolValue(notification.IncludeHealthWarnings),
		ID:                            types.Int64Value(notification.ID),
		Name:                          types.StringValue(notification.Name),
	}
	genericNotification.Tags, _ = types.SetValueFrom(ctx, types.Int64Type, notification.Tags)
	genericNotification.writeFields(ctx, notification.Fields)
	n.fromNotification(&genericNotification)
}

func (n *NotificationSlack) read(ctx context.Context) *sonarr.NotificationInput {
	var tags []int

	tfsdk.ValueAs(ctx, n.Tags, &tags)

	return &sonarr.NotificationInput{
		OnGrab:                        n.OnGrab.ValueBool(),
		OnDownload:                    n.OnDownload.ValueBool(),
		OnUpgrade:                     n.OnUpgrade.ValueBool(),
		OnRename:                      n.OnRename.ValueBool(),
		OnSeriesDelete:                n.OnSeriesDelete.ValueBool(),
		OnEpisodeFileDelete:           n.OnEpisodeFileDelete.ValueBool(),
		OnEpisodeFileDeleteForUpgrade: n.OnEpisodeFileDeleteForUpgrade.ValueBool(),
		OnHealthIssue:                 n.OnHealthIssue.ValueBool(),
		OnApplicationUpdate:           n.OnApplicationUpdate.ValueBool(),
		IncludeHealthWarnings:         n.IncludeHealthWarnings.ValueBool(),
		ConfigContract:                notificationSlackConfigContract,
		Implementation:                notificationSlackImplementation,
		ID:                            n.ID.ValueInt64(),
		Name:                          n.Name.ValueString(),
		Tags:                          tags,
		Fields:                        n.toNotification().readFields(ctx),
	}
}

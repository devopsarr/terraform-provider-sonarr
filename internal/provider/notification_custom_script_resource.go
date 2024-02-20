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
	notificationCustomScriptResourceName   = "notification_custom_script"
	notificationCustomScriptImplementation = "CustomScript"
	notificationCustomScriptConfigContract = "CustomScriptSettings"
)

// Ensure provider defined types fully satisfy framework interfaces.
var (
	_ resource.Resource                = &NotificationCustomScriptResource{}
	_ resource.ResourceWithImportState = &NotificationCustomScriptResource{}
)

func NewNotificationCustomScriptResource() resource.Resource {
	return &NotificationCustomScriptResource{}
}

// NotificationCustomScriptResource defines the notification implementation.
type NotificationCustomScriptResource struct {
	client *sonarr.APIClient
	auth   context.Context
}

// NotificationCustomScript describes the notification data model.
type NotificationCustomScript struct {
	Tags                          types.Set    `tfsdk:"tags"`
	Arguments                     types.String `tfsdk:"arguments"`
	Path                          types.String `tfsdk:"path"`
	Name                          types.String `tfsdk:"name"`
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
	OnRename                      types.Bool   `tfsdk:"on_rename"`
	OnUpgrade                     types.Bool   `tfsdk:"on_upgrade"`
	OnDownload                    types.Bool   `tfsdk:"on_download"`
}

func (n NotificationCustomScript) toNotification() *Notification {
	return &Notification{
		Tags:                          n.Tags,
		Path:                          n.Path,
		Arguments:                     n.Arguments,
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
		OnRename:                      n.OnRename,
		OnUpgrade:                     n.OnUpgrade,
		OnDownload:                    n.OnDownload,
		ConfigContract:                types.StringValue(notificationCustomScriptConfigContract),
		Implementation:                types.StringValue(notificationCustomScriptImplementation),
	}
}

func (n *NotificationCustomScript) fromNotification(notification *Notification) {
	n.Tags = notification.Tags
	n.Path = notification.Path
	n.Arguments = notification.Arguments
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
	n.OnRename = notification.OnRename
	n.OnUpgrade = notification.OnUpgrade
	n.OnDownload = notification.OnDownload
}

func (r *NotificationCustomScriptResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_" + notificationCustomScriptResourceName
}

func (r *NotificationCustomScriptResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "<!-- subcategory:Notifications -->Notification Custom Script resource.\nFor more information refer to [Notification](https://wiki.servarr.com/sonarr/settings#connect) and [Custom Script](https://wiki.servarr.com/sonarr/supported#customscript).",
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
				MarkdownDescription: "NotificationCustomScript name.",
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
			"arguments": schema.StringAttribute{
				MarkdownDescription: "Arguments.",
				Optional:            true,
				Computed:            true,
			},
			"path": schema.StringAttribute{
				MarkdownDescription: "Path.",
				Required:            true,
			},
		},
	}
}

func (r *NotificationCustomScriptResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	if auth, client := resourceConfigure(ctx, req, resp); client != nil {
		r.client = client
		r.auth = auth
	}
}

func (r *NotificationCustomScriptResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	// Retrieve values from plan
	var notification *NotificationCustomScript

	resp.Diagnostics.Append(req.Plan.Get(ctx, &notification)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Create new NotificationCustomScript
	request := notification.read(ctx, &resp.Diagnostics)

	response, _, err := r.client.NotificationAPI.CreateNotification(r.auth).NotificationResource(*request).Execute()
	if err != nil {
		resp.Diagnostics.AddError(helpers.ClientError, helpers.ParseClientError(helpers.Create, notificationCustomScriptResourceName, err))

		return
	}

	tflog.Trace(ctx, "created "+notificationCustomScriptResourceName+": "+strconv.Itoa(int(response.GetId())))
	// Generate resource state struct
	notification.write(ctx, response, &resp.Diagnostics)
	resp.Diagnostics.Append(resp.State.Set(ctx, &notification)...)
}

func (r *NotificationCustomScriptResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	// Get current state
	var notification *NotificationCustomScript

	resp.Diagnostics.Append(req.State.Get(ctx, &notification)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Get NotificationCustomScript current value
	response, _, err := r.client.NotificationAPI.GetNotificationById(r.auth, int32(notification.ID.ValueInt64())).Execute()
	if err != nil {
		resp.Diagnostics.AddError(helpers.ClientError, helpers.ParseClientError(helpers.Read, notificationCustomScriptResourceName, err))

		return
	}

	tflog.Trace(ctx, "read "+notificationCustomScriptResourceName+": "+strconv.Itoa(int(response.GetId())))
	// Map response body to resource schema attribute
	notification.write(ctx, response, &resp.Diagnostics)
	resp.Diagnostics.Append(resp.State.Set(ctx, &notification)...)
}

func (r *NotificationCustomScriptResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	// Get plan values
	var notification *NotificationCustomScript

	resp.Diagnostics.Append(req.Plan.Get(ctx, &notification)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Update NotificationCustomScript
	request := notification.read(ctx, &resp.Diagnostics)

	response, _, err := r.client.NotificationAPI.UpdateNotification(r.auth, strconv.Itoa(int(request.GetId()))).NotificationResource(*request).Execute()
	if err != nil {
		resp.Diagnostics.AddError(helpers.ClientError, helpers.ParseClientError(helpers.Update, notificationCustomScriptResourceName, err))

		return
	}

	tflog.Trace(ctx, "updated "+notificationCustomScriptResourceName+": "+strconv.Itoa(int(response.GetId())))
	// Generate resource state struct
	notification.write(ctx, response, &resp.Diagnostics)
	resp.Diagnostics.Append(resp.State.Set(ctx, &notification)...)
}

func (r *NotificationCustomScriptResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var ID int64

	resp.Diagnostics.Append(req.State.GetAttribute(ctx, path.Root("id"), &ID)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Delete NotificationCustomScript current value
	_, err := r.client.NotificationAPI.DeleteNotification(r.auth, int32(ID)).Execute()
	if err != nil {
		resp.Diagnostics.AddError(helpers.ClientError, helpers.ParseClientError(helpers.Delete, notificationCustomScriptResourceName, err))

		return
	}

	tflog.Trace(ctx, "deleted "+notificationCustomScriptResourceName+": "+strconv.Itoa(int(ID)))
	resp.State.RemoveResource(ctx)
}

func (r *NotificationCustomScriptResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	helpers.ImportStatePassthroughIntID(ctx, path.Root("id"), req, resp)
	tflog.Trace(ctx, "imported "+notificationCustomScriptResourceName+": "+req.ID)
}

func (n *NotificationCustomScript) write(ctx context.Context, notification *sonarr.NotificationResource, diags *diag.Diagnostics) {
	genericNotification := n.toNotification()
	genericNotification.write(ctx, notification, diags)
	n.fromNotification(genericNotification)
}

func (n *NotificationCustomScript) read(ctx context.Context, diags *diag.Diagnostics) *sonarr.NotificationResource {
	return n.toNotification().read(ctx, diags)
}

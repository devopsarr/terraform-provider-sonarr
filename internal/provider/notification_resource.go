package provider

import (
	"context"
	"fmt"
	"strconv"

	"github.com/devopsarr/terraform-provider-sonarr/internal/helpers"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"golift.io/starr"
	"golift.io/starr/sonarr"
)

// Ensure provider defined types fully satisfy framework interfaces.
var _ resource.Resource = &NotificationResource{}
var _ resource.ResourceWithImportState = &NotificationResource{}

func NewNotificationResource() resource.Resource {
	return &NotificationResource{}
}

// NotificationResource defines the notification implementation.
type NotificationResource struct {
	client *sonarr.Sonarr
}

// Notification describes the notification data model.
type Notification struct {
	OnGrab                        types.Bool   `tfsdk:"on_grab"`
	OnDownload                    types.Bool   `tfsdk:"on_download"`
	OnUpgrade                     types.Bool   `tfsdk:"on_upgrade"`
	OnRename                      types.Bool   `tfsdk:"on_rename"`
	OnSeriesDelete                types.Bool   `tfsdk:"on_series_delete"`
	OnEpisodeFileDelete           types.Bool   `tfsdk:"on_episode_file_delete"`
	OnEpisodeFileDeleteForUpgrade types.Bool   `tfsdk:"on_episode_file_delete_for_upgrade"`
	OnHealthIssue                 types.Bool   `tfsdk:"on_health_issue"`
	OnApplicationUpdate           types.Bool   `tfsdk:"on_application_update"`
	IncludeHealthWarnings         types.Bool   `tfsdk:"include_health_warnings"`
	ID                            types.Int64  `tfsdk:"id"`
	Name                          types.String `tfsdk:"name"`
	Implementation                types.String `tfsdk:"implementation"`
	ConfigContract                types.String `tfsdk:"config_contract"`
	Tags                          types.Set    `tfsdk:"tags"`
	// Fields values
	AlwaysUpdate      types.Bool   `tfsdk:"always_update"`
	CleanLibrary      types.Bool   `tfsdk:"clean_library"`
	DirectMessage     types.Bool   `tfsdk:"direct_message"`
	Notify            types.Bool   `tfsdk:"notify"`
	RequireEncryption types.Bool   `tfsdk:"require_encryption"`
	SendSilently      types.Bool   `tfsdk:"send_silently"`
	UpdateLibrary     types.Bool   `tfsdk:"update_library"`
	UseEuEndpoint     types.Bool   `tfsdk:"use_eu_endpoint"`
	UseSSL            types.Bool   `tfsdk:"use_ssl"`
	GrabFields        types.Int64  `tfsdk:"grab_fields"`   // 0-Overview 1-Rating 2-Genres 3-Quality 4-Group 5-Size 6-Links 7-Release 8-Poster 9-Fanart
	ImportFields      types.Int64  `tfsdk:"import_fields"` // 0-Overview 1-Rating 2-Genres 3-Quality 4-Codecs 5-Group 6-Size 7-Languages 8-Subtitles 9-Links 10-Release 11-Poster 12-Fanart
	Method            types.Int64  `tfsdk:"method"`        // 1-POST 2-PUT
	Port              types.Int64  `tfsdk:"port"`
	Priority          types.Int64  `tfsdk:"priority"` // 0-Min 2-Low 2-Normal 3-High | -2 Silent -1 Quiet 0 Normal 1 High 2 Emergency | -2 VeryLow -1 Low 0 Normal 1 High 2 Emergency
	AccessToken       types.String `tfsdk:"access_token"`
	AccessTokenSecret types.String `tfsdk:"access_token_secret"`
	APIKey            types.String `tfsdk:"api_key"`
	AppToken          types.String `tfsdk:"app_token"`
	Arguments         types.String `tfsdk:"arguments"`
	Author            types.String `tfsdk:"author"`
	AuthToken         types.String `tfsdk:"auth_token"`
	AuthUser          types.String `tfsdk:"auth_user"`
	Avatar            types.String `tfsdk:"avatar"`
	BCC               types.String `tfsdk:"bcc"`
	BotToken          types.String `tfsdk:"bot_token"`
	CC                types.String `tfsdk:"cc"`
	Channel           types.String `tfsdk:"channel"`
	ChatID            types.String `tfsdk:"chat_id"`
	ConsumerKey       types.String `tfsdk:"consumer_key"`
	ConsumerSecret    types.String `tfsdk:"consumer_secret"`
	DeviceNames       types.String `tfsdk:"device_names"`
	DisplayTime       types.String `tfsdk:"display_time"` // or types.Int64?
	Expire            types.String `tfsdk:"expire"`       // or types.Int64?
	Expires           types.String `tfsdk:"expires"`
	From              types.String `tfsdk:"from"`
	Host              types.String `tfsdk:"host"`
	Icon              types.String `tfsdk:"icon"`
	Mention           types.String `tfsdk:"mention"`
	Password          types.String `tfsdk:"password"`
	Path              types.String `tfsdk:"path"`
	RefreshToken      types.String `tfsdk:"refresh_token"`
	Retry             types.String `tfsdk:"retry"` // or types.Int64?
	SenderDomain      types.String `tfsdk:"sender_domain"`
	SenderID          types.String `tfsdk:"sender_id"`
	Server            types.String `tfsdk:"server"`
	SignIn            types.String `tfsdk:"sign_in"`
	Sound             types.String `tfsdk:"sound"`
	To                types.String `tfsdk:"to"`
	Token             types.String `tfsdk:"token"`
	URL               types.String `tfsdk:"url"`
	UserKey           types.String `tfsdk:"user_key"`
	Username          types.String `tfsdk:"username"`
	WebHookURL        types.String `tfsdk:"web_hook_url"`
	ChannelTags       types.Set    `tfsdk:"channel_tags"` // string
	DeviceIds         types.Set    `tfsdk:"device_ids"`   // string
	Devices           types.Set    `tfsdk:"devices"`      // string
	Recipients        types.Set    `tfsdk:"recipients"`   // string
	// AuthorizeNotification ?
}

func (r *NotificationResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_notification"
}

func (r *NotificationResource) GetSchema(ctx context.Context) (tfsdk.Schema, diag.Diagnostics) {
	return tfsdk.Schema{
		MarkdownDescription: "[subcategory:Notifications]: #\nNotification resource.\nFor more information refer to [Notification](https://wiki.servarr.com/sonarr/settings#connect).",
		Attributes: map[string]tfsdk.Attribute{
			"on_grab": {
				MarkdownDescription: "On grab flag.",
				Required:            true,
				Type:                types.BoolType,
			},
			"on_download": {
				MarkdownDescription: "On download flag.",
				Required:            true,
				Type:                types.BoolType,
			},
			"on_upgrade": {
				MarkdownDescription: "On upgrade flag.",
				Required:            true,
				Type:                types.BoolType,
			},
			"on_rename": {
				MarkdownDescription: "On rename flag.",
				Required:            true,
				Type:                types.BoolType,
			},
			"on_series_delete": {
				MarkdownDescription: "On series delete flag.",
				Required:            true,
				Type:                types.BoolType,
			},
			"on_episode_file_delete": {
				MarkdownDescription: "On episode file delete flag.",
				Required:            true,
				Type:                types.BoolType,
			},
			"on_episode_file_delete_for_upgrade": {
				MarkdownDescription: "On episode file delete for upgrade flag.",
				Required:            true,
				Type:                types.BoolType,
			},
			"on_health_issue": {
				MarkdownDescription: "On health issue flag.",
				Required:            true,
				Type:                types.BoolType,
			},
			"on_application_update": {
				MarkdownDescription: "On application update flag.",
				Required:            true,
				Type:                types.BoolType,
			},
			"include_health_warnings": {
				MarkdownDescription: "Include health warnings.",
				Required:            true,
				Type:                types.BoolType,
			},
			"config_contract": {
				MarkdownDescription: "Notification configuration template.",
				Required:            true,
				Type:                types.StringType,
			},
			"implementation": {
				MarkdownDescription: "Notification implementation name.",
				Required:            true,
				Type:                types.StringType,
			},
			"name": {
				MarkdownDescription: "Notification name.",
				Required:            true,
				Type:                types.StringType,
			},
			"tags": {
				MarkdownDescription: "List of associated tags.",
				Optional:            true,
				Computed:            true,
				Type: types.SetType{
					ElemType: types.Int64Type,
				},
			},
			"id": {
				MarkdownDescription: "Download Client ID.",
				Computed:            true,
				Type:                types.Int64Type,
				PlanModifiers: tfsdk.AttributePlanModifiers{
					resource.UseStateForUnknown(),
				},
			},
			// Field values
			"always_update": {
				MarkdownDescription: "Always update flag.",
				Optional:            true,
				Computed:            true,
				Type:                types.BoolType,
			},
			"clean_library": {
				MarkdownDescription: "Clean library flag.",
				Optional:            true,
				Computed:            true,
				Type:                types.BoolType,
			},
			"direct_message": {
				MarkdownDescription: "Direct message flag.",
				Optional:            true,
				Computed:            true,
				Type:                types.BoolType,
			},
			"notify": {
				MarkdownDescription: "Notify flag.",
				Optional:            true,
				Computed:            true,
				Type:                types.BoolType,
			},
			"require_encryption": {
				MarkdownDescription: "Require encryption flag.",
				Optional:            true,
				Computed:            true,
				Type:                types.BoolType,
			},
			"send_silently": {
				MarkdownDescription: "Add silently flag.",
				Optional:            true,
				Computed:            true,
				Type:                types.BoolType,
			},
			"update_library": {
				MarkdownDescription: "Update library flag.",
				Optional:            true,
				Computed:            true,
				Type:                types.BoolType,
			},
			"use_eu_endpoint": {
				MarkdownDescription: "Use EU endpoint flag.",
				Optional:            true,
				Computed:            true,
				Type:                types.BoolType,
			},
			"use_ssl": {
				MarkdownDescription: "Use SSL flag.",
				Optional:            true,
				Computed:            true,
				Type:                types.BoolType,
			},
			"port": {
				MarkdownDescription: "Port.",
				Optional:            true,
				Computed:            true,
				Type:                types.Int64Type,
			},
			"grab_fields": {
				MarkdownDescription: "Grab fields. `0` Overview, `1` Rating, `2` Genres, `3` Quality, `4` Group, `5` Size, `6` Links, `7` Release, `8` Poster, `9` Fanart.",
				Optional:            true,
				Computed:            true,
				Type:                types.Int64Type,
				Validators: []tfsdk.AttributeValidator{
					helpers.IntMatch([]int64{0, 1, 2, 3, 4, 5, 6, 7, 8, 9}),
				},
			},
			"import_fields": {
				MarkdownDescription: "Import fields. `0` Overview, `1` Rating, `2` Genres, `3` Quality, `4` Codecs, `5` Group, `6` Size, `7` Languages, `8` Subtitles, `9` Links, `10` Release, `11` Poster, `12` Fanart.",
				Optional:            true,
				Computed:            true,
				Type:                types.Int64Type,
				Validators: []tfsdk.AttributeValidator{
					helpers.IntMatch([]int64{0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12}),
				},
			},
			"method": {
				MarkdownDescription: "Method. `1` POST, `2` PUT.",
				Optional:            true,
				Computed:            true,
				Type:                types.Int64Type,
				Validators: []tfsdk.AttributeValidator{
					helpers.IntMatch([]int64{1, 2}),
				},
			},
			"priority": {
				MarkdownDescription: "Priority.", // TODO: add values in description
				Optional:            true,
				Computed:            true,
				Type:                types.Int64Type,
				Validators: []tfsdk.AttributeValidator{
					helpers.IntMatch([]int64{-2, -1, 0, 1, 2, 3, 4, 5, 7}),
				},
			},
			"access_token": {
				MarkdownDescription: "Access token.",
				Optional:            true,
				Computed:            true,
				Type:                types.StringType,
			},
			"access_token_secret": {
				MarkdownDescription: "Access token secret.",
				Optional:            true,
				Computed:            true,
				Type:                types.StringType,
			},
			"api_key": {
				MarkdownDescription: "API key.",
				Optional:            true,
				Computed:            true,
				Type:                types.StringType,
			},
			"app_token": {
				MarkdownDescription: "App token.",
				Optional:            true,
				Computed:            true,
				Type:                types.StringType,
			},
			"arguments": {
				MarkdownDescription: "Arguments.",
				Optional:            true,
				Computed:            true,
				Type:                types.StringType,
			},
			"author": {
				MarkdownDescription: "Author.",
				Optional:            true,
				Computed:            true,
				Type:                types.StringType,
			},
			"auth_token": {
				MarkdownDescription: "Auth token.",
				Optional:            true,
				Computed:            true,
				Type:                types.StringType,
			},
			"auth_user": {
				MarkdownDescription: "Auth user.",
				Optional:            true,
				Computed:            true,
				Type:                types.StringType,
			},
			"avatar": {
				MarkdownDescription: "Avatar.",
				Optional:            true,
				Computed:            true,
				Type:                types.StringType,
			},
			"bcc": {
				MarkdownDescription: "BCC.",
				Optional:            true,
				Computed:            true,
				Type:                types.StringType,
			},
			"bot_token": {
				MarkdownDescription: "Bot token.",
				Optional:            true,
				Computed:            true,
				Type:                types.StringType,
			},
			"cc": {
				MarkdownDescription: "CC.",
				Optional:            true,
				Computed:            true,
				Type:                types.StringType,
			},
			"channel": {
				MarkdownDescription: "Channel.",
				Optional:            true,
				Computed:            true,
				Type:                types.StringType,
			},
			"chat_id": {
				MarkdownDescription: "Chat ID.",
				Optional:            true,
				Computed:            true,
				Type:                types.StringType,
			},
			"consumer_key": {
				MarkdownDescription: "Consumer key.",
				Optional:            true,
				Computed:            true,
				Type:                types.StringType,
			},
			"consumer_secret": {
				MarkdownDescription: "Consumer secret.",
				Optional:            true,
				Computed:            true,
				Type:                types.StringType,
			},
			"device_names": {
				MarkdownDescription: "Device names.",
				Optional:            true,
				Computed:            true,
				Type:                types.StringType,
			},
			"display_time": {
				MarkdownDescription: "Display time.",
				Optional:            true,
				Computed:            true,
				Type:                types.StringType,
			},
			"expire": {
				MarkdownDescription: "Expire.",
				Optional:            true,
				Computed:            true,
				Type:                types.StringType,
			},
			"expires": {
				MarkdownDescription: "Expires.",
				Optional:            true,
				Computed:            true,
				Type:                types.StringType,
			},
			"from": {
				MarkdownDescription: "From.",
				Optional:            true,
				Computed:            true,
				Type:                types.StringType,
			},
			"host": {
				MarkdownDescription: "Host.",
				Optional:            true,
				Computed:            true,
				Type:                types.StringType,
			},
			"icon": {
				MarkdownDescription: "Icon.",
				Optional:            true,
				Computed:            true,
				Type:                types.StringType,
			},
			"mention": {
				MarkdownDescription: "Mention.",
				Optional:            true,
				Computed:            true,
				Type:                types.StringType,
			},
			"password": {
				MarkdownDescription: "password.",
				Optional:            true,
				Computed:            true,
				Type:                types.StringType,
			},
			"path": {
				MarkdownDescription: "Path.",
				Optional:            true,
				Computed:            true,
				Type:                types.StringType,
			},
			"refresh_token": {
				MarkdownDescription: "Refresh token.",
				Optional:            true,
				Computed:            true,
				Type:                types.StringType,
			},
			"retry": {
				MarkdownDescription: "Retry.",
				Optional:            true,
				Computed:            true,
				Type:                types.StringType,
			},
			"sender_domain": {
				MarkdownDescription: "Sender domain.",
				Optional:            true,
				Computed:            true,
				Type:                types.StringType,
			},
			"sender_id": {
				MarkdownDescription: "Sender ID.",
				Optional:            true,
				Computed:            true,
				Type:                types.StringType,
			},
			"server": {
				MarkdownDescription: "server.",
				Optional:            true,
				Computed:            true,
				Type:                types.StringType,
			},
			"sign_in": {
				MarkdownDescription: "Sign in.",
				Optional:            true,
				Computed:            true,
				Type:                types.StringType,
			},
			"sound": {
				MarkdownDescription: "Sound.",
				Optional:            true,
				Computed:            true,
				Type:                types.StringType,
			},
			"to": {
				MarkdownDescription: "To.",
				Optional:            true,
				Computed:            true,
				Type:                types.StringType,
			},
			"token": {
				MarkdownDescription: "Token.",
				Optional:            true,
				Computed:            true,
				Type:                types.StringType,
			},
			"url": {
				MarkdownDescription: "URL.",
				Optional:            true,
				Computed:            true,
				Type:                types.StringType,
			},
			"user_key": {
				MarkdownDescription: "User key.",
				Optional:            true,
				Computed:            true,
				Type:                types.StringType,
			},
			"username": {
				MarkdownDescription: "Username.",
				Optional:            true,
				Computed:            true,
				Type:                types.StringType,
			},
			"web_hook_url": {
				MarkdownDescription: "Web hook url.",
				Optional:            true,
				Computed:            true,
				Type:                types.StringType,
			},
			"channel_tags": {
				MarkdownDescription: "Channel tags.",
				Optional:            true,
				Computed:            true,
				Type: types.SetType{
					ElemType: types.StringType,
				},
			},
			"device_ids": {
				MarkdownDescription: "Device IDs.",
				Optional:            true,
				Computed:            true,
				Type: types.SetType{
					ElemType: types.StringType,
				},
			},
			"devices": {
				MarkdownDescription: "Devices.",
				Optional:            true,
				Computed:            true,
				Type: types.SetType{
					ElemType: types.StringType,
				},
			},
			"recipients": {
				MarkdownDescription: "Recipients.",
				Optional:            true,
				Computed:            true,
				Type: types.SetType{
					ElemType: types.StringType,
				},
			},
		},
	}, nil
}

func (r *NotificationResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	// Prevent panic if the provider has not been configured.
	if req.ProviderData == nil {
		return
	}

	client, ok := req.ProviderData.(*sonarr.Sonarr)
	if !ok {
		resp.Diagnostics.AddError(
			UnexpectedResourceConfigureType,
			fmt.Sprintf("Expected *sonarr.Sonarr, got: %T. Please report this issue to the provider developers.", req.ProviderData),
		)

		return
	}

	r.client = client
}

func (r *NotificationResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	// Retrieve values from plan
	var plan Notification

	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Create new Notification
	request := readNotification(ctx, &plan)

	response, err := r.client.AddNotificationContext(ctx, request)
	if err != nil {
		resp.Diagnostics.AddError(ClientError, fmt.Sprintf("Unable to create Notification, got error: %s", err))

		return
	}

	tflog.Trace(ctx, "created notification: "+strconv.Itoa(int(response.ID)))
	// Generate resource state struct
	result := writeNotification(ctx, response)
	resp.Diagnostics.Append(resp.State.Set(ctx, &result)...)
}

func (r *NotificationResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	// Get current state
	var state Notification

	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Get Notification current value
	response, err := r.client.GetNotificationContext(ctx, int(state.ID.Value))
	if err != nil {
		resp.Diagnostics.AddError(ClientError, fmt.Sprintf("Unable to read Notifications, got error: %s", err))

		return
	}

	tflog.Trace(ctx, "read notification: "+strconv.Itoa(int(response.ID)))
	// Map response body to resource schema attribute
	result := writeNotification(ctx, response)
	resp.Diagnostics.Append(resp.State.Set(ctx, &result)...)
}

func (r *NotificationResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	// Get plan values
	var plan Notification

	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Update Notification
	request := readNotification(ctx, &plan)

	response, err := r.client.UpdateNotificationContext(ctx, request)
	if err != nil {
		resp.Diagnostics.AddError(ClientError, fmt.Sprintf("Unable to update Notification, got error: %s", err))

		return
	}

	tflog.Trace(ctx, "updated notification: "+strconv.Itoa(int(response.ID)))
	// Generate resource state struct
	result := writeNotification(ctx, response)
	resp.Diagnostics.Append(resp.State.Set(ctx, &result)...)
}

func (r *NotificationResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state Notification

	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Delete Notification current value
	err := r.client.DeleteNotificationContext(ctx, int(state.ID.Value))
	if err != nil {
		resp.Diagnostics.AddError(ClientError, fmt.Sprintf("Unable to read Notifications, got error: %s", err))

		return
	}

	tflog.Trace(ctx, "deleted notification: "+strconv.Itoa(int(state.ID.Value)))
	resp.State.RemoveResource(ctx)
}

func (r *NotificationResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	// resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
	id, err := strconv.Atoi(req.ID)
	if err != nil {
		resp.Diagnostics.AddError(
			UnexpectedImportIdentifier,
			fmt.Sprintf("Expected import identifier with format: ID. Got: %q", req.ID),
		)

		return
	}

	tflog.Trace(ctx, "imported notification: "+strconv.Itoa(id))
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("id"), id)...)
}

func writeNotification(ctx context.Context, notification *sonarr.NotificationOutput) *Notification {
	output := Notification{
		OnGrab:                        types.Bool{Value: notification.OnGrab},
		OnDownload:                    types.Bool{Value: notification.OnDownload},
		OnUpgrade:                     types.Bool{Value: notification.OnUpgrade},
		OnRename:                      types.Bool{Value: notification.OnRename},
		OnSeriesDelete:                types.Bool{Value: notification.OnSeriesDelete},
		OnEpisodeFileDelete:           types.Bool{Value: notification.OnEpisodeFileDelete},
		OnEpisodeFileDeleteForUpgrade: types.Bool{Value: notification.OnEpisodeFileDeleteForUpgrade},
		OnHealthIssue:                 types.Bool{Value: notification.OnHealthIssue},
		OnApplicationUpdate:           types.Bool{Value: notification.OnApplicationUpdate},
		IncludeHealthWarnings:         types.Bool{Value: notification.IncludeHealthWarnings},
		ID:                            types.Int64{Value: notification.ID},
		Name:                          types.String{Value: notification.Name},
		Implementation:                types.String{Value: notification.Implementation},
		ConfigContract:                types.String{Value: notification.ConfigContract},
		Tags:                          types.Set{ElemType: types.Int64Type},
		ChannelTags:                   types.Set{ElemType: types.StringType},
		DeviceIds:                     types.Set{ElemType: types.StringType},
		Devices:                       types.Set{ElemType: types.StringType},
		Recipients:                    types.Set{ElemType: types.StringType},
	}
	tfsdk.ValueFrom(ctx, notification.Tags, output.Tags.Type(ctx), &output.Tags)

	for _, f := range notification.Fields {
		if f.Value != nil {
			switch f.Name {
			case "alwaysUpdate":
				output.AlwaysUpdate = types.Bool{Value: f.Value.(bool)}
			case "cleanLibrary":
				output.CleanLibrary = types.Bool{Value: f.Value.(bool)}
			case "directMessage":
				output.DirectMessage = types.Bool{Value: f.Value.(bool)}
			case "notify":
				output.Notify = types.Bool{Value: f.Value.(bool)}
			case "requireEncryption":
				output.RequireEncryption = types.Bool{Value: f.Value.(bool)}
			case "sendSilently":
				output.SendSilently = types.Bool{Value: f.Value.(bool)}
			case "updateLibrary":
				output.UpdateLibrary = types.Bool{Value: f.Value.(bool)}
			case "useEuEndpoint":
				output.UseEuEndpoint = types.Bool{Value: f.Value.(bool)}
			case "useSSL":
				output.UseSSL = types.Bool{Value: f.Value.(bool)}
			case "grabFields":
				output.GrabFields = types.Int64{Value: int64(f.Value.(float64))}
			case "importFields":
				output.ImportFields = types.Int64{Value: int64(f.Value.(float64))}
			case "method":
				output.Method = types.Int64{Value: int64(f.Value.(float64))}
			case "port":
				output.Port = types.Int64{Value: int64(f.Value.(float64))}
			case "priority":
				output.Priority = types.Int64{Value: int64(f.Value.(float64))}
			case "accessToken":
				output.AccessToken = types.String{Value: f.Value.(string)}
			case "accessTokenSecret":
				output.AccessTokenSecret = types.String{Value: f.Value.(string)}
			case "apiKey":
				output.APIKey = types.String{Value: f.Value.(string)}
			case "appToken":
				output.AppToken = types.String{Value: f.Value.(string)}
			case "arguments":
				output.Arguments = types.String{Value: f.Value.(string)}
			case "author":
				output.Author = types.String{Value: f.Value.(string)}
			case "authToken":
				output.AuthToken = types.String{Value: f.Value.(string)}
			case "authUser":
				output.AuthUser = types.String{Value: f.Value.(string)}
			case "avatar":
				output.Avatar = types.String{Value: f.Value.(string)}
			case "bcc":
				output.BCC = types.String{Value: f.Value.(string)}
			case "botToken":
				output.BotToken = types.String{Value: f.Value.(string)}
			case "cc":
				output.CC = types.String{Value: f.Value.(string)}
			case "channel":
				output.Channel = types.String{Value: f.Value.(string)}
			case "chatId":
				output.ChatID = types.String{Value: f.Value.(string)}
			case "consumerKey":
				output.ConsumerKey = types.String{Value: f.Value.(string)}
			case "consumerSecret":
				output.ConsumerSecret = types.String{Value: f.Value.(string)}
			case "deviceNames":
				output.DeviceNames = types.String{Value: f.Value.(string)}
			case "displayTime":
				output.DisplayTime = types.String{Value: f.Value.(string)}
			case "expire":
				output.Expire = types.String{Value: f.Value.(string)}
			case "expires":
				output.Expires = types.String{Value: f.Value.(string)}
			case "from":
				output.From = types.String{Value: f.Value.(string)}
			case "host":
				output.Host = types.String{Value: f.Value.(string)}
			case "icon":
				output.Icon = types.String{Value: f.Value.(string)}
			case "mention":
				output.Mention = types.String{Value: f.Value.(string)}
			case "password":
				output.Password = types.String{Value: f.Value.(string)}
			case "path":
				output.Path = types.String{Value: f.Value.(string)}
			case "refreshToken":
				output.RefreshToken = types.String{Value: f.Value.(string)}
			case "retry":
				output.Retry = types.String{Value: f.Value.(string)}
			case "senderDomain":
				output.SenderDomain = types.String{Value: f.Value.(string)}
			case "senderId":
				output.SenderID = types.String{Value: f.Value.(string)}
			case "server":
				output.Server = types.String{Value: f.Value.(string)}
			case "signIn":
				output.SignIn = types.String{Value: f.Value.(string)}
			case "sound":
				output.Sound = types.String{Value: f.Value.(string)}
			case "to":
				output.To = types.String{Value: f.Value.(string)}
			case "token":
				output.Token = types.String{Value: f.Value.(string)}
			case "url":
				output.URL = types.String{Value: f.Value.(string)}
			case "userKey":
				output.UserKey = types.String{Value: f.Value.(string)}
			case "username":
				output.Username = types.String{Value: f.Value.(string)}
			case "webHookUrl":
				output.WebHookURL = types.String{Value: f.Value.(string)}
			case "channelTags":
				tfsdk.ValueFrom(ctx, f.Value, output.ChannelTags.Type(ctx), &output.ChannelTags)
			case "deviceIds":
				tfsdk.ValueFrom(ctx, f.Value, output.DeviceIds.Type(ctx), &output.DeviceIds)
			case "devices":
				tfsdk.ValueFrom(ctx, f.Value, output.Devices.Type(ctx), &output.Devices)
			case "recipients":
				tfsdk.ValueFrom(ctx, f.Value, output.Recipients.Type(ctx), &output.Recipients)
				// TODO: manage unknown values
			default:
			}
		}
	}

	return &output
}

func readNotification(ctx context.Context, notification *Notification) *sonarr.NotificationInput {
	var tags []int

	tfsdk.ValueAs(ctx, notification.Tags, &tags)

	return &sonarr.NotificationInput{
		OnGrab:                        notification.OnGrab.Value,
		OnDownload:                    notification.OnDownload.Value,
		OnUpgrade:                     notification.OnUpgrade.Value,
		OnRename:                      notification.OnRename.Value,
		OnSeriesDelete:                notification.OnSeriesDelete.Value,
		OnEpisodeFileDelete:           notification.OnEpisodeFileDelete.Value,
		OnEpisodeFileDeleteForUpgrade: notification.OnEpisodeFileDeleteForUpgrade.Value,
		OnHealthIssue:                 notification.OnHealthIssue.Value,
		OnApplicationUpdate:           notification.OnApplicationUpdate.Value,
		IncludeHealthWarnings:         notification.IncludeHealthWarnings.Value,
		ID:                            notification.ID.Value,
		Name:                          notification.Name.Value,
		Implementation:                notification.Implementation.Value,
		ConfigContract:                notification.ConfigContract.Value,
		Tags:                          tags,
		Fields:                        readNotificationFields(ctx, notification),
	}
}

func readNotificationFields(ctx context.Context, notification *Notification) []*starr.FieldInput {
	var output []*starr.FieldInput
	if !notification.AlwaysUpdate.IsNull() && !notification.AlwaysUpdate.IsUnknown() {
		output = append(output, &starr.FieldInput{
			Name:  "alwaysUpdate",
			Value: notification.AlwaysUpdate.Value,
		})
	}

	if !notification.CleanLibrary.IsNull() && !notification.CleanLibrary.IsUnknown() {
		output = append(output, &starr.FieldInput{
			Name:  "cleanLibrary",
			Value: notification.CleanLibrary.Value,
		})
	}

	if !notification.DirectMessage.IsNull() && !notification.DirectMessage.IsUnknown() {
		output = append(output, &starr.FieldInput{
			Name:  "directMessage",
			Value: notification.DirectMessage.Value,
		})
	}

	if !notification.Notify.IsNull() && !notification.Notify.IsUnknown() {
		output = append(output, &starr.FieldInput{
			Name:  "notify",
			Value: notification.Notify.Value,
		})
	}

	if !notification.RequireEncryption.IsNull() && !notification.RequireEncryption.IsUnknown() {
		output = append(output, &starr.FieldInput{
			Name:  "requireEncryption",
			Value: notification.RequireEncryption.Value,
		})
	}

	if !notification.SendSilently.IsNull() && !notification.SendSilently.IsUnknown() {
		output = append(output, &starr.FieldInput{
			Name:  "sendSilently",
			Value: notification.SendSilently.Value,
		})
	}

	if !notification.UpdateLibrary.IsNull() && !notification.UpdateLibrary.IsUnknown() {
		output = append(output, &starr.FieldInput{
			Name:  "updateLibrary",
			Value: notification.UpdateLibrary.Value,
		})
	}

	if !notification.UseEuEndpoint.IsNull() && !notification.UseEuEndpoint.IsUnknown() {
		output = append(output, &starr.FieldInput{
			Name:  "useEuEndpoint",
			Value: notification.UseEuEndpoint.Value,
		})
	}

	if !notification.UseSSL.IsNull() && !notification.UseSSL.IsUnknown() {
		output = append(output, &starr.FieldInput{
			Name:  "useSSL",
			Value: notification.UseSSL.Value,
		})
	}

	if !notification.GrabFields.IsNull() && !notification.GrabFields.IsUnknown() {
		output = append(output, &starr.FieldInput{
			Name:  "grabFields",
			Value: notification.GrabFields.Value,
		})
	}

	if !notification.ImportFields.IsNull() && !notification.ImportFields.IsUnknown() {
		output = append(output, &starr.FieldInput{
			Name:  "importFields",
			Value: notification.ImportFields.Value,
		})
	}

	if !notification.Method.IsNull() && !notification.Method.IsUnknown() {
		output = append(output, &starr.FieldInput{
			Name:  "method",
			Value: notification.Method.Value,
		})
	}

	if !notification.Port.IsNull() && !notification.Port.IsUnknown() {
		output = append(output, &starr.FieldInput{
			Name:  "port",
			Value: notification.Port.Value,
		})
	}

	if !notification.Priority.IsNull() && !notification.Priority.IsUnknown() {
		output = append(output, &starr.FieldInput{
			Name:  "priority",
			Value: notification.Priority.Value,
		})
	}

	if !notification.AccessToken.IsNull() && !notification.AccessToken.IsUnknown() {
		output = append(output, &starr.FieldInput{
			Name:  "accessToken",
			Value: notification.AccessToken.Value,
		})
	}

	if !notification.AccessTokenSecret.IsNull() && !notification.AccessTokenSecret.IsUnknown() {
		output = append(output, &starr.FieldInput{
			Name:  "accessTokenSecret",
			Value: notification.AccessTokenSecret.Value,
		})
	}

	if !notification.APIKey.IsNull() && !notification.APIKey.IsUnknown() {
		output = append(output, &starr.FieldInput{
			Name:  "apiKey",
			Value: notification.APIKey.Value,
		})
	}

	if !notification.AppToken.IsNull() && !notification.AppToken.IsUnknown() {
		output = append(output, &starr.FieldInput{
			Name:  "appToken",
			Value: notification.AppToken.Value,
		})
	}

	if !notification.Arguments.IsNull() && !notification.Arguments.IsUnknown() {
		output = append(output, &starr.FieldInput{
			Name:  "arguments",
			Value: notification.Arguments.Value,
		})
	}

	if !notification.Author.IsNull() && !notification.Author.IsUnknown() {
		output = append(output, &starr.FieldInput{
			Name:  "author",
			Value: notification.Author.Value,
		})
	}

	if !notification.AuthToken.IsNull() && !notification.AuthToken.IsUnknown() {
		output = append(output, &starr.FieldInput{
			Name:  "authToken",
			Value: notification.AuthToken.Value,
		})
	}

	if !notification.AuthUser.IsNull() && !notification.AuthUser.IsUnknown() {
		output = append(output, &starr.FieldInput{
			Name:  "authUser",
			Value: notification.AuthUser.Value,
		})
	}

	if !notification.Avatar.IsNull() && !notification.Avatar.IsUnknown() {
		output = append(output, &starr.FieldInput{
			Name:  "avatar",
			Value: notification.Avatar.Value,
		})
	}

	if !notification.BCC.IsNull() && !notification.BCC.IsUnknown() {
		output = append(output, &starr.FieldInput{
			Name:  "bcc",
			Value: notification.BCC.Value,
		})
	}

	if !notification.BotToken.IsNull() && !notification.BotToken.IsUnknown() {
		output = append(output, &starr.FieldInput{
			Name:  "botToken",
			Value: notification.BotToken.Value,
		})
	}

	if !notification.CC.IsNull() && !notification.CC.IsUnknown() {
		output = append(output, &starr.FieldInput{
			Name:  "cc",
			Value: notification.CC.Value,
		})
	}

	if !notification.Channel.IsNull() && !notification.Channel.IsUnknown() {
		output = append(output, &starr.FieldInput{
			Name:  "channel",
			Value: notification.Channel.Value,
		})
	}

	if !notification.ChatID.IsNull() && !notification.ChatID.IsUnknown() {
		output = append(output, &starr.FieldInput{
			Name:  "chatId",
			Value: notification.ChatID.Value,
		})
	}

	if !notification.ConsumerKey.IsNull() && !notification.ConsumerKey.IsUnknown() {
		output = append(output, &starr.FieldInput{
			Name:  "consumerKey",
			Value: notification.ConsumerKey.Value,
		})
	}

	if !notification.ConsumerSecret.IsNull() && !notification.ConsumerSecret.IsUnknown() {
		output = append(output, &starr.FieldInput{
			Name:  "consumerSecret",
			Value: notification.ConsumerSecret.Value,
		})
	}

	if !notification.DeviceNames.IsNull() && !notification.DeviceNames.IsUnknown() {
		output = append(output, &starr.FieldInput{
			Name:  "deviceNames",
			Value: notification.DeviceNames.Value,
		})
	}

	if !notification.DisplayTime.IsNull() && !notification.DisplayTime.IsUnknown() {
		output = append(output, &starr.FieldInput{
			Name:  "displayTime",
			Value: notification.DisplayTime.Value,
		})
	}

	if !notification.Expire.IsNull() && !notification.Expire.IsUnknown() {
		output = append(output, &starr.FieldInput{
			Name:  "expire",
			Value: notification.Expire.Value,
		})
	}

	if !notification.Expires.IsNull() && !notification.Expires.IsUnknown() {
		output = append(output, &starr.FieldInput{
			Name:  "expires",
			Value: notification.Expires.Value,
		})
	}

	if !notification.From.IsNull() && !notification.From.IsUnknown() {
		output = append(output, &starr.FieldInput{
			Name:  "from",
			Value: notification.From.Value,
		})
	}

	if !notification.Host.IsNull() && !notification.Host.IsUnknown() {
		output = append(output, &starr.FieldInput{
			Name:  "host",
			Value: notification.Host.Value,
		})
	}

	if !notification.Icon.IsNull() && !notification.Icon.IsUnknown() {
		output = append(output, &starr.FieldInput{
			Name:  "icon",
			Value: notification.Icon.Value,
		})
	}

	if !notification.Mention.IsNull() && !notification.Mention.IsUnknown() {
		output = append(output, &starr.FieldInput{
			Name:  "mention",
			Value: notification.Mention.Value,
		})
	}

	if !notification.Password.IsNull() && !notification.Password.IsUnknown() {
		output = append(output, &starr.FieldInput{
			Name:  "password",
			Value: notification.Password.Value,
		})
	}

	if !notification.Path.IsNull() && !notification.Path.IsUnknown() {
		output = append(output, &starr.FieldInput{
			Name:  "path",
			Value: notification.Path.Value,
		})
	}

	if !notification.RefreshToken.IsNull() && !notification.RefreshToken.IsUnknown() {
		output = append(output, &starr.FieldInput{
			Name:  "refreshToken",
			Value: notification.RefreshToken.Value,
		})
	}

	if !notification.Retry.IsNull() && !notification.Retry.IsUnknown() {
		output = append(output, &starr.FieldInput{
			Name:  "retry",
			Value: notification.Retry.Value,
		})
	}

	if !notification.SenderDomain.IsNull() && !notification.SenderDomain.IsUnknown() {
		output = append(output, &starr.FieldInput{
			Name:  "senderDomain",
			Value: notification.SenderDomain.Value,
		})
	}

	if !notification.SenderID.IsNull() && !notification.SenderID.IsUnknown() {
		output = append(output, &starr.FieldInput{
			Name:  "senderId",
			Value: notification.SenderID.Value,
		})
	}

	if !notification.Server.IsNull() && !notification.Server.IsUnknown() {
		output = append(output, &starr.FieldInput{
			Name:  "server",
			Value: notification.Server.Value,
		})
	}

	if !notification.SignIn.IsNull() && !notification.SignIn.IsUnknown() {
		output = append(output, &starr.FieldInput{
			Name:  "signIn",
			Value: notification.SignIn.Value,
		})
	}

	if !notification.Sound.IsNull() && !notification.Sound.IsUnknown() {
		output = append(output, &starr.FieldInput{
			Name:  "sound",
			Value: notification.Sound.Value,
		})
	}

	if !notification.To.IsNull() && !notification.To.IsUnknown() {
		output = append(output, &starr.FieldInput{
			Name:  "to",
			Value: notification.To.Value,
		})
	}

	if !notification.Token.IsNull() && !notification.Token.IsUnknown() {
		output = append(output, &starr.FieldInput{
			Name:  "token",
			Value: notification.Token.Value,
		})
	}

	if !notification.URL.IsNull() && !notification.URL.IsUnknown() {
		output = append(output, &starr.FieldInput{
			Name:  "url",
			Value: notification.URL.Value,
		})
	}

	if !notification.UserKey.IsNull() && !notification.UserKey.IsUnknown() {
		output = append(output, &starr.FieldInput{
			Name:  "userKey",
			Value: notification.UserKey.Value,
		})
	}

	if !notification.Username.IsNull() && !notification.Username.IsUnknown() {
		output = append(output, &starr.FieldInput{
			Name:  "username",
			Value: notification.Username.Value,
		})
	}

	if !notification.WebHookURL.IsNull() && !notification.WebHookURL.IsUnknown() {
		output = append(output, &starr.FieldInput{
			Name:  "webHookUrl",
			Value: notification.WebHookURL.Value,
		})
	}

	if len(notification.ChannelTags.Elems) != 0 {
		tags := make([]types.String, len(notification.ChannelTags.Elems))
		tfsdk.ValueAs(ctx, notification.ChannelTags, &tags)

		output = append(output, &starr.FieldInput{
			Name:  "channelTags",
			Value: tags,
		})
	}

	if len(notification.DeviceIds.Elems) != 0 {
		tags := make([]types.String, len(notification.DeviceIds.Elems))
		tfsdk.ValueAs(ctx, notification.DeviceIds, &tags)

		output = append(output, &starr.FieldInput{
			Name:  "deviceIds",
			Value: tags,
		})
	}

	if len(notification.Devices.Elems) != 0 {
		tags := make([]types.String, len(notification.Devices.Elems))
		tfsdk.ValueAs(ctx, notification.Devices, &tags)

		output = append(output, &starr.FieldInput{
			Name:  "devices",
			Value: tags,
		})
	}

	if len(notification.Recipients.Elems) != 0 {
		tags := make([]types.String, len(notification.Recipients.Elems))
		tfsdk.ValueAs(ctx, notification.Recipients, &tags)

		output = append(output, &starr.FieldInput{
			Name:  "recipients",
			Value: tags,
		})
	}

	return output
}

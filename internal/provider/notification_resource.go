package provider

import (
	"context"
	"strconv"

	"github.com/devopsarr/sonarr-go/sonarr"
	"github.com/devopsarr/terraform-provider-sonarr/internal/helpers"
	"github.com/hashicorp/terraform-plugin-framework-validators/int64validator"
	"github.com/hashicorp/terraform-plugin-framework/attr"
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

const notificationResourceName = "notification"

// Ensure provider defined types fully satisfy framework interfaces.
var (
	_ resource.Resource                = &NotificationResource{}
	_ resource.ResourceWithImportState = &NotificationResource{}
)

var notificationFields = helpers.Fields{
	Bools:                  []string{"alwaysUpdate", "cleanLibrary", "directMessage", "notify", "sendSilently", "updateLibrary", "useEuEndpoint", "useSsl"},
	Strings:                []string{"accessToken", "accessTokenSecret", "apiKey", "appToken", "arguments", "author", "authToken", "authUser", "avatar", "botToken", "channel", "chatId", "consumerKey", "consumerSecret", "deviceNames", "expires", "from", "host", "icon", "mention", "password", "path", "refreshToken", "senderDomain", "senderId", "server", "signIn", "sound", "token", "url", "userKey", "username", "userName", "webHookUrl", "clickUrl", "serverUrl", "authUsername", "authPassword", "statelessUrls", "configurationKey", "senderNumber", "receiverId", "key", "event"},
	Ints:                   []string{"method", "port", "priority", "retry", "expire", "displayTime", "notificationType", "useEncryption"},
	StringSlices:           []string{"channelTags", "deviceIds", "devices", "recipients", "to", "cc", "bcc", "topics", "fieldTags"},
	StringSlicesExceptions: []string{"tags"},
	IntSlices:              []string{"grabFields", "importFields"},
}

func NewNotificationResource() resource.Resource {
	return &NotificationResource{}
}

// NotificationResource defines the notification implementation.
type NotificationResource struct {
	client *sonarr.APIClient
	auth   context.Context
}

// Notification describes the notification data model.
type Notification struct {
	Tags                          types.Set    `tfsdk:"tags"`
	FieldTags                     types.Set    `tfsdk:"field_tags"`
	Topics                        types.Set    `tfsdk:"topics"`
	Recipients                    types.Set    `tfsdk:"recipients"`
	Devices                       types.Set    `tfsdk:"devices"`
	DeviceIDs                     types.Set    `tfsdk:"device_ids"`
	ChannelTags                   types.Set    `tfsdk:"channel_tags"`
	ImportFields                  types.Set    `tfsdk:"import_fields"`
	GrabFields                    types.Set    `tfsdk:"grab_fields"`
	To                            types.Set    `tfsdk:"to"`
	Cc                            types.Set    `tfsdk:"cc"`
	Bcc                           types.Set    `tfsdk:"bcc"`
	Path                          types.String `tfsdk:"path"`
	RefreshToken                  types.String `tfsdk:"refresh_token"`
	WebHookURL                    types.String `tfsdk:"web_hook_url"`
	Username                      types.String `tfsdk:"username"`
	UserKey                       types.String `tfsdk:"user_key"`
	Mention                       types.String `tfsdk:"mention"`
	ClickURL                      types.String `tfsdk:"click_url"`
	ServerURL                     types.String `tfsdk:"server_url"`
	StatelessURLs                 types.String `tfsdk:"stateless_urls"`
	Name                          types.String `tfsdk:"name"`
	Avatar                        types.String `tfsdk:"avatar"`
	ConfigContract                types.String `tfsdk:"config_contract"`
	URL                           types.String `tfsdk:"url"`
	Token                         types.String `tfsdk:"token"`
	Sound                         types.String `tfsdk:"sound"`
	SignIn                        types.String `tfsdk:"sign_in"`
	Server                        types.String `tfsdk:"server"`
	SenderID                      types.String `tfsdk:"sender_id"`
	SenderNumber                  types.String `tfsdk:"sender_number"`
	ReceiverID                    types.String `tfsdk:"receiver_id"`
	BotToken                      types.String `tfsdk:"bot_token"`
	SenderDomain                  types.String `tfsdk:"sender_domain"`
	Icon                          types.String `tfsdk:"icon"`
	Host                          types.String `tfsdk:"host"`
	From                          types.String `tfsdk:"from"`
	Expires                       types.String `tfsdk:"expires"`
	AccessToken                   types.String `tfsdk:"access_token"`
	AccessTokenSecret             types.String `tfsdk:"access_token_secret"`
	APIKey                        types.String `tfsdk:"api_key"`
	AppToken                      types.String `tfsdk:"app_token"`
	Arguments                     types.String `tfsdk:"arguments"`
	Author                        types.String `tfsdk:"author"`
	AuthToken                     types.String `tfsdk:"auth_token"`
	AuthUser                      types.String `tfsdk:"auth_user"`
	Implementation                types.String `tfsdk:"implementation"`
	Password                      types.String `tfsdk:"password"`
	Channel                       types.String `tfsdk:"channel"`
	ChatID                        types.String `tfsdk:"chat_id"`
	ConsumerKey                   types.String `tfsdk:"consumer_key"`
	ConsumerSecret                types.String `tfsdk:"consumer_secret"`
	DeviceNames                   types.String `tfsdk:"device_names"`
	AuthUsername                  types.String `tfsdk:"auth_username"`
	AuthPassword                  types.String `tfsdk:"auth_password"`
	ConfigurationKey              types.String `tfsdk:"configuration_key"`
	Key                           types.String `tfsdk:"key"`
	Event                         types.String `tfsdk:"event"`
	NotificationType              types.Int64  `tfsdk:"notification_type"`
	Expire                        types.Int64  `tfsdk:"expire"`
	DisplayTime                   types.Int64  `tfsdk:"display_time"`
	Priority                      types.Int64  `tfsdk:"priority"`
	Port                          types.Int64  `tfsdk:"port"`
	Method                        types.Int64  `tfsdk:"method"`
	Retry                         types.Int64  `tfsdk:"retry"`
	UseEncryption                 types.Int64  `tfsdk:"use_encryption"`
	ID                            types.Int64  `tfsdk:"id"`
	UpdateLibrary                 types.Bool   `tfsdk:"update_library"`
	OnGrab                        types.Bool   `tfsdk:"on_grab"`
	UseEuEndpoint                 types.Bool   `tfsdk:"use_eu_endpoint"`
	Notify                        types.Bool   `tfsdk:"notify"`
	UseSSL                        types.Bool   `tfsdk:"use_ssl"`
	OnEpisodeFileDeleteForUpgrade types.Bool   `tfsdk:"on_episode_file_delete_for_upgrade"`
	SendSilently                  types.Bool   `tfsdk:"send_silently"`
	DirectMessage                 types.Bool   `tfsdk:"direct_message"`
	CleanLibrary                  types.Bool   `tfsdk:"clean_library"`
	AlwaysUpdate                  types.Bool   `tfsdk:"always_update"`
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
	OnImportComplete              types.Bool   `tfsdk:"on_import_complete"`
}

func (n Notification) getType() attr.Type {
	return types.ObjectType{}.WithAttributeTypes(
		map[string]attr.Type{
			"tags":                               types.SetType{}.WithElementType(types.Int64Type),
			"import_fields":                      types.SetType{}.WithElementType(types.Int64Type),
			"grab_fields":                        types.SetType{}.WithElementType(types.Int64Type),
			"field_tags":                         types.SetType{}.WithElementType(types.StringType),
			"recipients":                         types.SetType{}.WithElementType(types.StringType),
			"devices":                            types.SetType{}.WithElementType(types.StringType),
			"device_ids":                         types.SetType{}.WithElementType(types.StringType),
			"to":                                 types.SetType{}.WithElementType(types.StringType),
			"cc":                                 types.SetType{}.WithElementType(types.StringType),
			"bcc":                                types.SetType{}.WithElementType(types.StringType),
			"channel_tags":                       types.SetType{}.WithElementType(types.StringType),
			"topics":                             types.SetType{}.WithElementType(types.StringType),
			"path":                               types.StringType,
			"refresh_token":                      types.StringType,
			"web_hook_url":                       types.StringType,
			"username":                           types.StringType,
			"user_key":                           types.StringType,
			"mention":                            types.StringType,
			"click_url":                          types.StringType,
			"server_url":                         types.StringType,
			"stateless_urls":                     types.StringType,
			"name":                               types.StringType,
			"avatar":                             types.StringType,
			"config_contract":                    types.StringType,
			"url":                                types.StringType,
			"token":                              types.StringType,
			"sound":                              types.StringType,
			"sign_in":                            types.StringType,
			"server":                             types.StringType,
			"sender_id":                          types.StringType,
			"sender_number":                      types.StringType,
			"receiver_id":                        types.StringType,
			"bot_token":                          types.StringType,
			"sender_domain":                      types.StringType,
			"icon":                               types.StringType,
			"host":                               types.StringType,
			"from":                               types.StringType,
			"expires":                            types.StringType,
			"access_token":                       types.StringType,
			"access_token_secret":                types.StringType,
			"api_key":                            types.StringType,
			"app_token":                          types.StringType,
			"arguments":                          types.StringType,
			"author":                             types.StringType,
			"auth_token":                         types.StringType,
			"auth_user":                          types.StringType,
			"implementation":                     types.StringType,
			"password":                           types.StringType,
			"channel":                            types.StringType,
			"chat_id":                            types.StringType,
			"consumer_key":                       types.StringType,
			"consumer_secret":                    types.StringType,
			"device_names":                       types.StringType,
			"auth_username":                      types.StringType,
			"auth_password":                      types.StringType,
			"configuration_key":                  types.StringType,
			"key":                                types.StringType,
			"event":                              types.StringType,
			"notification_type":                  types.Int64Type,
			"expire":                             types.Int64Type,
			"display_time":                       types.Int64Type,
			"priority":                           types.Int64Type,
			"port":                               types.Int64Type,
			"method":                             types.Int64Type,
			"retry":                              types.Int64Type,
			"use_encryption":                     types.Int64Type,
			"id":                                 types.Int64Type,
			"update_library":                     types.BoolType,
			"on_grab":                            types.BoolType,
			"use_eu_endpoint":                    types.BoolType,
			"notify":                             types.BoolType,
			"use_ssl":                            types.BoolType,
			"on_episode_file_delete_for_upgrade": types.BoolType,
			"send_silently":                      types.BoolType,
			"direct_message":                     types.BoolType,
			"clean_library":                      types.BoolType,
			"always_update":                      types.BoolType,
			"on_episode_file_delete":             types.BoolType,
			"include_health_warnings":            types.BoolType,
			"on_application_update":              types.BoolType,
			"on_health_issue":                    types.BoolType,
			"on_health_restored":                 types.BoolType,
			"on_manual_interaction_required":     types.BoolType,
			"on_series_add":                      types.BoolType,
			"on_series_delete":                   types.BoolType,
			"on_rename":                          types.BoolType,
			"on_upgrade":                         types.BoolType,
			"on_download":                        types.BoolType,
			"on_import_complete":                 types.BoolType,
		})
}

func (r *NotificationResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_" + notificationResourceName
}

func (r *NotificationResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "<!-- subcategory:Notifications -->\nGeneric Notification resource. When possible use a specific resource instead.\nFor more information refer to [Notification](https://wiki.servarr.com/sonarr/settings#connect).",
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
			"on_import_complete": schema.BoolAttribute{
				MarkdownDescription: "On import complete flag.",
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
			"on_application_update": schema.BoolAttribute{
				MarkdownDescription: "On application update flag.",
				Optional:            true,
				Computed:            true,
			},
			"on_manual_interaction_required": schema.BoolAttribute{
				MarkdownDescription: "On manual interaction required flag.",
				Optional:            true,
				Computed:            true,
			},
			"include_health_warnings": schema.BoolAttribute{
				MarkdownDescription: "Include health warnings.",
				Optional:            true,
				Computed:            true,
			},
			"config_contract": schema.StringAttribute{
				MarkdownDescription: "Notification configuration template.",
				Required:            true,
			},
			"implementation": schema.StringAttribute{
				MarkdownDescription: "Notification implementation name.",
				Required:            true,
			},
			"name": schema.StringAttribute{
				MarkdownDescription: "Notification name.",
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
			"always_update": schema.BoolAttribute{
				MarkdownDescription: "Always update flag.",
				Optional:            true,
				Computed:            true,
			},
			"clean_library": schema.BoolAttribute{
				MarkdownDescription: "Clean library flag.",
				Optional:            true,
				Computed:            true,
			},
			"direct_message": schema.BoolAttribute{
				MarkdownDescription: "Direct message flag.",
				Optional:            true,
				Computed:            true,
			},
			"notify": schema.BoolAttribute{
				MarkdownDescription: "Notify flag.",
				Optional:            true,
				Computed:            true,
			},
			"send_silently": schema.BoolAttribute{
				MarkdownDescription: "Add silently flag.",
				Optional:            true,
				Computed:            true,
			},
			"update_library": schema.BoolAttribute{
				MarkdownDescription: "Update library flag.",
				Optional:            true,
				Computed:            true,
			},
			"use_eu_endpoint": schema.BoolAttribute{
				MarkdownDescription: "Use EU endpoint flag.",
				Optional:            true,
				Computed:            true,
			},
			"use_ssl": schema.BoolAttribute{
				MarkdownDescription: "Use SSL flag.",
				Optional:            true,
				Computed:            true,
			},
			"expire": schema.Int64Attribute{
				MarkdownDescription: "Expire.",
				Optional:            true,
				Computed:            true,
			},
			"port": schema.Int64Attribute{
				MarkdownDescription: "Port.",
				Optional:            true,
				Computed:            true,
			},
			"method": schema.Int64Attribute{
				MarkdownDescription: "Method. `1` POST, `2` PUT.",
				Optional:            true,
				Computed:            true,
				Validators: []validator.Int64{
					int64validator.OneOf(1, 2),
				},
			},
			"use_encryption": schema.Int64Attribute{
				MarkdownDescription: "Require encryption. `0` Preferred, `1` Always, `2` Never.",
				Optional:            true,
				Computed:            true,
				Validators: []validator.Int64{
					int64validator.OneOf(0, 1, 2),
				},
			},
			"priority": schema.Int64Attribute{
				MarkdownDescription: "Priority.", // TODO: add values in description
				Optional:            true,
				Computed:            true,
				Validators: []validator.Int64{
					int64validator.OneOf(-2, -1, 0, 1, 2, 3, 4, 5, 7, 8),
				},
			},
			"retry": schema.Int64Attribute{
				MarkdownDescription: "Retry.",
				Optional:            true,
				Computed:            true,
			},
			"notification_type": schema.Int64Attribute{
				MarkdownDescription: "Notification type. `0` Info, `1` Success, `2` Warning, `3` Failure.",
				Optional:            true,
				Computed:            true,
				Validators: []validator.Int64{
					int64validator.OneOf(0, 1, 2, 3),
				},
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
			"access_token": schema.StringAttribute{
				MarkdownDescription: "Access token.",
				Optional:            true,
				Computed:            true,
			},
			"access_token_secret": schema.StringAttribute{
				MarkdownDescription: "Access token secret.",
				Optional:            true,
				Computed:            true,
				Sensitive:           true,
			},
			"api_key": schema.StringAttribute{
				MarkdownDescription: "API key.",
				Optional:            true,
				Computed:            true,
				Sensitive:           true,
			},
			"app_token": schema.StringAttribute{
				MarkdownDescription: "App token.",
				Optional:            true,
				Computed:            true,
				Sensitive:           true,
			},
			"key": schema.StringAttribute{
				MarkdownDescription: "Key.",
				Optional:            true,
				Computed:            true,
				Sensitive:           true,
			},
			"event": schema.StringAttribute{
				MarkdownDescription: "Event.",
				Optional:            true,
				Computed:            true,
			},
			"arguments": schema.StringAttribute{
				MarkdownDescription: "Arguments.",
				Optional:            true,
				Computed:            true,
			},
			"author": schema.StringAttribute{
				MarkdownDescription: "Author.",
				Optional:            true,
				Computed:            true,
			},
			"server_url": schema.StringAttribute{
				MarkdownDescription: "Server URL.",
				Optional:            true,
				Computed:            true,
			},
			"click_url": schema.StringAttribute{
				MarkdownDescription: "Click URL.",
				Optional:            true,
				Computed:            true,
			},
			"auth_token": schema.StringAttribute{
				MarkdownDescription: "Auth token.",
				Optional:            true,
				Computed:            true,
				Sensitive:           true,
			},
			"auth_user": schema.StringAttribute{
				MarkdownDescription: "Auth user.",
				Optional:            true,
				Computed:            true,
			},
			"avatar": schema.StringAttribute{
				MarkdownDescription: "Avatar.",
				Optional:            true,
				Computed:            true,
			},
			"bot_token": schema.StringAttribute{
				MarkdownDescription: "Bot token.",
				Optional:            true,
				Computed:            true,
			},
			"channel": schema.StringAttribute{
				MarkdownDescription: "Channel.",
				Optional:            true,
				Computed:            true,
			},
			"chat_id": schema.StringAttribute{
				MarkdownDescription: "Chat ID.",
				Optional:            true,
				Computed:            true,
			},
			"consumer_key": schema.StringAttribute{
				MarkdownDescription: "Consumer key.",
				Optional:            true,
				Computed:            true,
			},
			"consumer_secret": schema.StringAttribute{
				MarkdownDescription: "Consumer secret.",
				Optional:            true,
				Computed:            true,
				Sensitive:           true,
			},
			"device_names": schema.StringAttribute{
				MarkdownDescription: "Device names.",
				Optional:            true,
				Computed:            true,
			},
			"display_time": schema.Int64Attribute{
				MarkdownDescription: "Display time.",
				Optional:            true,
				Computed:            true,
			},
			"expires": schema.StringAttribute{
				MarkdownDescription: "Expires.",
				Optional:            true,
				Computed:            true,
			},
			"from": schema.StringAttribute{
				MarkdownDescription: "From.",
				Optional:            true,
				Computed:            true,
			},
			"host": schema.StringAttribute{
				MarkdownDescription: "Host.",
				Optional:            true,
				Computed:            true,
			},
			"icon": schema.StringAttribute{
				MarkdownDescription: "Icon.",
				Optional:            true,
				Computed:            true,
			},
			"mention": schema.StringAttribute{
				MarkdownDescription: "Mention.",
				Optional:            true,
				Computed:            true,
			},
			"password": schema.StringAttribute{
				MarkdownDescription: "password.",
				Optional:            true,
				Computed:            true,
				Sensitive:           true,
			},
			"path": schema.StringAttribute{
				MarkdownDescription: "Path.",
				Optional:            true,
				Computed:            true,
			},
			"refresh_token": schema.StringAttribute{
				MarkdownDescription: "Refresh token.",
				Optional:            true,
				Computed:            true,
			},
			"sender_domain": schema.StringAttribute{
				MarkdownDescription: "Sender domain.",
				Optional:            true,
				Computed:            true,
			},
			"sender_id": schema.StringAttribute{
				MarkdownDescription: "Sender ID.",
				Optional:            true,
				Computed:            true,
			},
			"sender_number": schema.StringAttribute{
				MarkdownDescription: "Sender Number.",
				Optional:            true,
				Computed:            true,
			},
			"receiver_id": schema.StringAttribute{
				MarkdownDescription: "Receiver ID.",
				Optional:            true,
				Computed:            true,
			},
			"server": schema.StringAttribute{
				MarkdownDescription: "server.",
				Optional:            true,
				Computed:            true,
			},
			"sign_in": schema.StringAttribute{
				MarkdownDescription: "Sign in.",
				Optional:            true,
				Computed:            true,
			},
			"sound": schema.StringAttribute{
				MarkdownDescription: "Sound.",
				Optional:            true,
				Computed:            true,
			},
			"token": schema.StringAttribute{
				MarkdownDescription: "Token.",
				Optional:            true,
				Computed:            true,
				Sensitive:           true,
			},
			"url": schema.StringAttribute{
				MarkdownDescription: "URL.",
				Optional:            true,
				Computed:            true,
			},
			"user_key": schema.StringAttribute{
				MarkdownDescription: "User key.",
				Optional:            true,
				Computed:            true,
			},
			"username": schema.StringAttribute{
				MarkdownDescription: "Username.",
				Optional:            true,
				Computed:            true,
			},
			"web_hook_url": schema.StringAttribute{
				MarkdownDescription: "Web hook url.",
				Optional:            true,
				Computed:            true,
			},
			"grab_fields": schema.SetAttribute{
				MarkdownDescription: "Grab fields. `0` Overview, `1` Rating, `2` Genres, `3` Quality, `4` Group, `5` Size, `6` Links, `7` Release, `8` Poster, `9` Fanart.",
				Optional:            true,
				Computed:            true,
				ElementType:         types.Int64Type,
			},
			"import_fields": schema.SetAttribute{
				MarkdownDescription: "Import fields. `0` Overview, `1` Rating, `2` Genres, `3` Quality, `4` Codecs, `5` Group, `6` Size, `7` Languages, `8` Subtitles, `9` Links, `10` Release, `11` Poster, `12` Fanart.",
				Optional:            true,
				Computed:            true,
				ElementType:         types.Int64Type,
			},
			"channel_tags": schema.SetAttribute{
				MarkdownDescription: "Channel tags.",
				Optional:            true,
				Computed:            true,
				ElementType:         types.StringType,
			},
			"device_ids": schema.SetAttribute{
				MarkdownDescription: "Device IDs.",
				Optional:            true,
				Computed:            true,
				ElementType:         types.StringType,
			},
			"devices": schema.SetAttribute{
				MarkdownDescription: "Devices.",
				Optional:            true,
				Computed:            true,
				ElementType:         types.StringType,
			},
			"recipients": schema.SetAttribute{
				MarkdownDescription: "Recipients.",
				Optional:            true,
				Computed:            true,
				ElementType:         types.StringType,
			},
			"to": schema.SetAttribute{
				MarkdownDescription: "To.",
				Optional:            true,
				Computed:            true,
				ElementType:         types.StringType,
			},
			"cc": schema.SetAttribute{
				MarkdownDescription: "Cc.",
				Optional:            true,
				Computed:            true,
				ElementType:         types.StringType,
			},
			"bcc": schema.SetAttribute{
				MarkdownDescription: "Bcc.",
				Optional:            true,
				Computed:            true,
				ElementType:         types.StringType,
			},
			"topics": schema.SetAttribute{
				MarkdownDescription: "Topics.",
				Computed:            true,
				Optional:            true,
				ElementType:         types.StringType,
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

func (r *NotificationResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	if auth, client := resourceConfigure(ctx, req, resp); client != nil {
		r.client = client
		r.auth = auth
	}
}

func (r *NotificationResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	// Retrieve values from plan
	var notification *Notification

	resp.Diagnostics.Append(req.Plan.Get(ctx, &notification)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Create new Notification
	request := notification.read(ctx, &resp.Diagnostics)

	response, _, err := r.client.NotificationAPI.CreateNotification(r.auth).NotificationResource(*request).Execute()
	if err != nil {
		resp.Diagnostics.AddError(helpers.ClientError, helpers.ParseClientError(helpers.Create, notificationResourceName, err))

		return
	}

	tflog.Trace(ctx, "created "+notificationResourceName+": "+strconv.Itoa(int(response.GetId())))
	// Generate resource state struct
	// this is needed because of many empty fields are unknown in both plan and read
	var state Notification

	state.writeSensitive(notification)
	state.write(ctx, response, &resp.Diagnostics)
	resp.Diagnostics.Append(resp.State.Set(ctx, state)...)
}

func (r *NotificationResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	// Get current state
	var notification *Notification

	resp.Diagnostics.Append(req.State.Get(ctx, &notification)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Get Notification current value
	response, _, err := r.client.NotificationAPI.GetNotificationById(r.auth, int32(notification.ID.ValueInt64())).Execute()
	if err != nil {
		resp.Diagnostics.AddError(helpers.ClientError, helpers.ParseClientError(helpers.Read, notificationResourceName, err))

		return
	}

	tflog.Trace(ctx, "read "+notificationResourceName+": "+strconv.Itoa(int(response.GetId())))
	// Map response body to resource schema attribute
	// this is needed because of many empty fields are unknown in both plan and read
	var state Notification

	state.writeSensitive(notification)
	state.write(ctx, response, &resp.Diagnostics)
	resp.Diagnostics.Append(resp.State.Set(ctx, state)...)
}

func (r *NotificationResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	// Get plan values
	var notification *Notification

	resp.Diagnostics.Append(req.Plan.Get(ctx, &notification)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Update Notification
	request := notification.read(ctx, &resp.Diagnostics)

	response, _, err := r.client.NotificationAPI.UpdateNotification(r.auth, strconv.Itoa(int(request.GetId()))).NotificationResource(*request).Execute()
	if err != nil {
		resp.Diagnostics.AddError(helpers.ClientError, helpers.ParseClientError(helpers.Update, notificationResourceName, err))

		return
	}

	tflog.Trace(ctx, "updated "+notificationResourceName+": "+strconv.Itoa(int(response.GetId())))
	// Generate resource state struct
	// this is needed because of many empty fields are unknown in both plan and read
	var state Notification

	state.writeSensitive(notification)
	state.write(ctx, response, &resp.Diagnostics)
	resp.Diagnostics.Append(resp.State.Set(ctx, state)...)
}

func (r *NotificationResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var ID int64

	resp.Diagnostics.Append(req.State.GetAttribute(ctx, path.Root("id"), &ID)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Delete Notification current value
	_, err := r.client.NotificationAPI.DeleteNotification(r.auth, int32(ID)).Execute()
	if err != nil {
		resp.Diagnostics.AddError(helpers.ClientError, helpers.ParseClientError(helpers.Delete, notificationResourceName, err))

		return
	}

	tflog.Trace(ctx, "deleted "+notificationResourceName+": "+strconv.Itoa(int(ID)))
	resp.State.RemoveResource(ctx)
}

func (r *NotificationResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	helpers.ImportStatePassthroughIntID(ctx, path.Root("id"), req, resp)
	tflog.Trace(ctx, "imported "+notificationResourceName+": "+req.ID)
}

func (n *Notification) write(ctx context.Context, notification *sonarr.NotificationResource, diags *diag.Diagnostics) {
	var localDiag diag.Diagnostics

	n.Tags, localDiag = types.SetValueFrom(ctx, types.Int64Type, notification.Tags)
	diags.Append(localDiag...)

	n.OnGrab = types.BoolValue(notification.GetOnGrab())
	n.OnDownload = types.BoolValue(notification.GetOnDownload())
	n.OnUpgrade = types.BoolValue(notification.GetOnUpgrade())
	n.OnRename = types.BoolValue(notification.GetOnRename())
	n.OnSeriesAdd = types.BoolValue(notification.GetOnSeriesAdd())
	n.OnSeriesDelete = types.BoolValue(notification.GetOnSeriesDelete())
	n.OnEpisodeFileDelete = types.BoolValue(notification.GetOnEpisodeFileDelete())
	n.OnEpisodeFileDeleteForUpgrade = types.BoolValue(notification.GetOnEpisodeFileDeleteForUpgrade())
	n.OnHealthIssue = types.BoolValue(notification.GetOnHealthIssue())
	n.OnHealthRestored = types.BoolValue(notification.GetOnHealthRestored())
	n.OnApplicationUpdate = types.BoolValue(notification.GetOnApplicationUpdate())
	n.OnManualInteractionRequired = types.BoolValue(notification.GetOnManualInteractionRequired())
	n.OnImportComplete = types.BoolValue(notification.GetOnImportComplete())
	n.IncludeHealthWarnings = types.BoolValue(notification.GetIncludeHealthWarnings())
	n.ID = types.Int64Value(int64(notification.GetId()))
	n.Name = types.StringValue(notification.GetName())
	n.Implementation = types.StringValue(notification.GetImplementation())
	n.ConfigContract = types.StringValue(notification.GetConfigContract())
	n.ImportFields = types.SetValueMust(types.Int64Type, nil)
	n.GrabFields = types.SetValueMust(types.Int64Type, nil)
	n.ChannelTags = types.SetValueMust(types.StringType, nil)
	n.DeviceIDs = types.SetValueMust(types.StringType, nil)
	n.Devices = types.SetValueMust(types.StringType, nil)
	n.Recipients = types.SetValueMust(types.StringType, nil)
	n.To = types.SetValueMust(types.StringType, nil)
	n.Cc = types.SetValueMust(types.StringType, nil)
	n.Bcc = types.SetValueMust(types.StringType, nil)
	n.Topics = types.SetValueMust(types.StringType, nil)
	n.FieldTags = types.SetValueMust(types.StringType, nil)
	helpers.WriteFields(ctx, n, notification.GetFields(), notificationFields)
}

func (n *Notification) read(ctx context.Context, diags *diag.Diagnostics) *sonarr.NotificationResource {
	notification := sonarr.NewNotificationResource()
	notification.SetOnGrab(n.OnGrab.ValueBool())
	notification.SetOnDownload(n.OnDownload.ValueBool())
	notification.SetOnUpgrade(n.OnUpgrade.ValueBool())
	notification.SetOnRename(n.OnRename.ValueBool())
	notification.SetOnSeriesAdd(n.OnSeriesAdd.ValueBool())
	notification.SetOnSeriesDelete(n.OnSeriesDelete.ValueBool())
	notification.SetOnEpisodeFileDelete(n.OnEpisodeFileDelete.ValueBool())
	notification.SetOnEpisodeFileDeleteForUpgrade(n.OnEpisodeFileDeleteForUpgrade.ValueBool())
	notification.SetOnHealthIssue(n.OnHealthIssue.ValueBool())
	notification.SetOnHealthRestored(n.OnHealthRestored.ValueBool())
	notification.SetOnApplicationUpdate(n.OnApplicationUpdate.ValueBool())
	notification.SetOnManualInteractionRequired(n.OnManualInteractionRequired.ValueBool())
	notification.SetOnImportComplete(n.OnImportComplete.ValueBool())
	notification.SetIncludeHealthWarnings(n.IncludeHealthWarnings.ValueBool())
	notification.SetId(int32(n.ID.ValueInt64()))
	notification.SetName(n.Name.ValueString())
	notification.SetImplementation(n.Implementation.ValueString())
	notification.SetConfigContract(n.ConfigContract.ValueString())
	diags.Append(n.Tags.ElementsAs(ctx, &notification.Tags, true)...)
	notification.SetFields(helpers.ReadFields(ctx, n, notificationFields))

	return notification
}

// writeSensitive copy sensitive data from another resource.
func (n *Notification) writeSensitive(notification *Notification) {
	if !notification.Token.IsUnknown() {
		n.Token = notification.Token
	}

	if !notification.APIKey.IsUnknown() {
		n.APIKey = notification.APIKey
	}

	if !notification.Password.IsUnknown() {
		n.Password = notification.Password
	}

	if !notification.AppToken.IsUnknown() {
		n.AppToken = notification.AppToken
	}

	if !notification.BotToken.IsUnknown() {
		n.BotToken = notification.BotToken
	}

	if !notification.AccessToken.IsUnknown() {
		n.AccessToken = notification.AccessToken
	}

	if !notification.AccessTokenSecret.IsUnknown() {
		n.AccessTokenSecret = notification.AccessTokenSecret
	}

	if !notification.ConsumerKey.IsUnknown() {
		n.ConsumerKey = notification.ConsumerKey
	}

	if !notification.ConsumerSecret.IsUnknown() {
		n.ConsumerSecret = notification.ConsumerSecret
	}

	if !notification.ConfigurationKey.IsUnknown() {
		n.ConfigurationKey = notification.ConfigurationKey
	}

	if !notification.AuthPassword.IsUnknown() {
		n.AuthPassword = notification.AuthPassword
	}
}

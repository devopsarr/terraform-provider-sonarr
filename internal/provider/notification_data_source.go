package provider

import (
	"context"
	"fmt"

	"github.com/devopsarr/terraform-provider-sonarr/internal/helpers"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"golift.io/starr/sonarr"
)

const notificationDataSourceName = "notification"

// Ensure provider defined types fully satisfy framework interfaces.
var _ datasource.DataSource = &NotificationDataSource{}

func NewNotificationDataSource() datasource.DataSource {
	return &NotificationDataSource{}
}

// NotificationDataSource defines the notification implementation.
type NotificationDataSource struct {
	client *sonarr.Sonarr
}

func (d *NotificationDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_" + notificationDataSourceName
}

func (d *NotificationDataSource) GetSchema(ctx context.Context) (tfsdk.Schema, diag.Diagnostics) {
	return tfsdk.Schema{
		// This description is used by the documentation generator and the delay server.
		MarkdownDescription: "<!-- subcategory:Notifications -->Single [Notification](../resources/notification).",
		Attributes: map[string]tfsdk.Attribute{
			"on_grab": {
				MarkdownDescription: "On grab flag.",
				Computed:            true,
				Type:                types.BoolType,
			},
			"on_download": {
				MarkdownDescription: "On download flag.",
				Computed:            true,
				Type:                types.BoolType,
			},
			"on_upgrade": {
				MarkdownDescription: "On upgrade flag.",
				Computed:            true,
				Type:                types.BoolType,
			},
			"on_rename": {
				MarkdownDescription: "On rename flag.",
				Computed:            true,
				Type:                types.BoolType,
			},
			"on_series_delete": {
				MarkdownDescription: "On series delete flag.",
				Computed:            true,
				Type:                types.BoolType,
			},
			"on_episode_file_delete": {
				MarkdownDescription: "On episode file delete flag.",
				Computed:            true,
				Type:                types.BoolType,
			},
			"on_episode_file_delete_for_upgrade": {
				MarkdownDescription: "On episode file delete for upgrade flag.",
				Computed:            true,
				Type:                types.BoolType,
			},
			"on_health_issue": {
				MarkdownDescription: "On health issue flag.",
				Computed:            true,
				Type:                types.BoolType,
			},
			"on_application_update": {
				MarkdownDescription: "On application update flag.",
				Computed:            true,
				Type:                types.BoolType,
			},
			"include_health_warnings": {
				MarkdownDescription: "Include health warnings.",
				Computed:            true,
				Type:                types.BoolType,
			},
			"config_contract": {
				MarkdownDescription: "Notification configuration template.",
				Computed:            true,
				Type:                types.StringType,
			},
			"implementation": {
				MarkdownDescription: "Notification implementation name.",
				Computed:            true,
				Type:                types.StringType,
			},
			"name": {
				MarkdownDescription: "Notification name.",
				Required:            true,
				Type:                types.StringType,
			},
			"tags": {
				MarkdownDescription: "List of associated tags.",
				Computed:            true,
				Type: types.SetType{
					ElemType: types.Int64Type,
				},
			},
			"id": {
				MarkdownDescription: "Notification ID.",
				Computed:            true,
				Type:                types.Int64Type,
				PlanModifiers: tfsdk.AttributePlanModifiers{
					resource.UseStateForUnknown(),
				},
			},
			// Field values
			"always_update": {
				MarkdownDescription: "Always update flag.",
				Computed:            true,
				Type:                types.BoolType,
			},
			"clean_library": {
				MarkdownDescription: "Clean library flag.",
				Computed:            true,
				Type:                types.BoolType,
			},
			"direct_message": {
				MarkdownDescription: "Direct message flag.",
				Computed:            true,
				Type:                types.BoolType,
			},
			"notify": {
				MarkdownDescription: "Notify flag.",
				Computed:            true,
				Type:                types.BoolType,
			},
			"require_encryption": {
				MarkdownDescription: "Require encryption flag.",
				Computed:            true,
				Type:                types.BoolType,
			},
			"send_silently": {
				MarkdownDescription: "Add silently flag.",
				Computed:            true,
				Type:                types.BoolType,
			},
			"update_library": {
				MarkdownDescription: "Update library flag.",
				Computed:            true,
				Type:                types.BoolType,
			},
			"use_eu_endpoint": {
				MarkdownDescription: "Use EU endpoint flag.",
				Computed:            true,
				Type:                types.BoolType,
			},
			"use_ssl": {
				MarkdownDescription: "Use SSL flag.",
				Computed:            true,
				Type:                types.BoolType,
			},
			"port": {
				MarkdownDescription: "Port.",
				Computed:            true,
				Type:                types.Int64Type,
			},
			"grab_fields": {
				MarkdownDescription: "Grab fields. `0` Overview, `1` Rating, `2` Genres, `3` Quality, `4` Group, `5` Size, `6` Links, `7` Release, `8` Poster, `9` Fanart.",
				Computed:            true,
				Type:                types.Int64Type,
			},
			"import_fields": {
				MarkdownDescription: "Import fields. `0` Overview, `1` Rating, `2` Genres, `3` Quality, `4` Codecs, `5` Group, `6` Size, `7` Languages, `8` Subtitles, `9` Links, `10` Release, `11` Poster, `12` Fanart.",
				Computed:            true,
				Type:                types.Int64Type,
			},
			"method": {
				MarkdownDescription: "Method. `1` POST, `2` PUT.",
				Computed:            true,
				Type:                types.Int64Type,
				Validators: []tfsdk.AttributeValidator{
					helpers.IntMatch([]int64{1, 2}),
				},
			},
			"priority": {
				MarkdownDescription: "Priority.", // TODO: add values in description
				Computed:            true,
				Type:                types.Int64Type,
			},
			"access_token": {
				MarkdownDescription: "Access token.",
				Computed:            true,
				Type:                types.StringType,
			},
			"access_token_secret": {
				MarkdownDescription: "Access token secret.",
				Computed:            true,
				Type:                types.StringType,
			},
			"api_key": {
				MarkdownDescription: "API key.",
				Computed:            true,
				Type:                types.StringType,
			},
			"app_token": {
				MarkdownDescription: "App token.",
				Computed:            true,
				Type:                types.StringType,
			},
			"arguments": {
				MarkdownDescription: "Arguments.",
				Computed:            true,
				Type:                types.StringType,
			},
			"author": {
				MarkdownDescription: "Author.",
				Computed:            true,
				Type:                types.StringType,
			},
			"auth_token": {
				MarkdownDescription: "Auth token.",
				Computed:            true,
				Type:                types.StringType,
			},
			"auth_user": {
				MarkdownDescription: "Auth user.",
				Computed:            true,
				Type:                types.StringType,
			},
			"avatar": {
				MarkdownDescription: "Avatar.",
				Computed:            true,
				Type:                types.StringType,
			},
			"bcc": {
				MarkdownDescription: "BCC.",
				Computed:            true,
				Type:                types.StringType,
			},
			"bot_token": {
				MarkdownDescription: "Bot token.",
				Computed:            true,
				Type:                types.StringType,
			},
			"cc": {
				MarkdownDescription: "CC.",
				Computed:            true,
				Type:                types.StringType,
			},
			"channel": {
				MarkdownDescription: "Channel.",
				Computed:            true,
				Type:                types.StringType,
			},
			"chat_id": {
				MarkdownDescription: "Chat ID.",
				Computed:            true,
				Type:                types.StringType,
			},
			"consumer_key": {
				MarkdownDescription: "Consumer key.",
				Computed:            true,
				Type:                types.StringType,
			},
			"consumer_secret": {
				MarkdownDescription: "Consumer secret.",
				Computed:            true,
				Type:                types.StringType,
			},
			"device_names": {
				MarkdownDescription: "Device names.",
				Computed:            true,
				Type:                types.StringType,
			},
			"display_time": {
				MarkdownDescription: "Display time.",
				Computed:            true,
				Type:                types.StringType,
			},
			"expire": {
				MarkdownDescription: "Expire.",
				Computed:            true,
				Type:                types.StringType,
			},
			"expires": {
				MarkdownDescription: "Expires.",
				Computed:            true,
				Type:                types.StringType,
			},
			"from": {
				MarkdownDescription: "From.",
				Computed:            true,
				Type:                types.StringType,
			},
			"host": {
				MarkdownDescription: "Host.",
				Computed:            true,
				Type:                types.StringType,
			},
			"icon": {
				MarkdownDescription: "Icon.",
				Computed:            true,
				Type:                types.StringType,
			},
			"mention": {
				MarkdownDescription: "Mention.",
				Computed:            true,
				Type:                types.StringType,
			},
			"password": {
				MarkdownDescription: "password.",
				Computed:            true,
				Type:                types.StringType,
			},
			"path": {
				MarkdownDescription: "Path.",
				Computed:            true,
				Type:                types.StringType,
			},
			"refresh_token": {
				MarkdownDescription: "Refresh token.",
				Computed:            true,
				Type:                types.StringType,
			},
			"retry": {
				MarkdownDescription: "Retry.",
				Computed:            true,
				Type:                types.StringType,
			},
			"sender_domain": {
				MarkdownDescription: "Sender domain.",
				Computed:            true,
				Type:                types.StringType,
			},
			"sender_id": {
				MarkdownDescription: "Sender ID.",
				Computed:            true,
				Type:                types.StringType,
			},
			"server": {
				MarkdownDescription: "server.",
				Computed:            true,
				Type:                types.StringType,
			},
			"sign_in": {
				MarkdownDescription: "Sign in.",
				Computed:            true,
				Type:                types.StringType,
			},
			"sound": {
				MarkdownDescription: "Sound.",
				Computed:            true,
				Type:                types.StringType,
			},
			"to": {
				MarkdownDescription: "To.",
				Computed:            true,
				Type:                types.StringType,
			},
			"token": {
				MarkdownDescription: "Token.",
				Computed:            true,
				Type:                types.StringType,
			},
			"url": {
				MarkdownDescription: "URL.",
				Computed:            true,
				Type:                types.StringType,
			},
			"user_key": {
				MarkdownDescription: "User key.",
				Computed:            true,
				Type:                types.StringType,
			},
			"username": {
				MarkdownDescription: "Username.",
				Computed:            true,
				Type:                types.StringType,
			},
			"web_hook_url": {
				MarkdownDescription: "Web hook url.",
				Computed:            true,
				Type:                types.StringType,
			},
			"channel_tags": {
				MarkdownDescription: "Channel tags.",
				Computed:            true,
				Type: types.SetType{
					ElemType: types.StringType,
				},
			},
			"device_ids": {
				MarkdownDescription: "Device IDs.",
				Computed:            true,
				Type: types.SetType{
					ElemType: types.StringType,
				},
			},
			"devices": {
				MarkdownDescription: "Devices.",
				Computed:            true,
				Type: types.SetType{
					ElemType: types.StringType,
				},
			},
			"recipients": {
				MarkdownDescription: "Recipients.",
				Computed:            true,
				Type: types.SetType{
					ElemType: types.StringType,
				},
			},
		},
	}, nil
}

func (d *NotificationDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *NotificationDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data *Notification

	resp.Diagnostics.Append(resp.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}
	// Get notification current value
	response, err := d.client.GetNotificationsContext(ctx)
	if err != nil {
		resp.Diagnostics.AddError(helpers.ClientError, fmt.Sprintf("Unable to read %s, got error: %s", notificationDataSourceName, err))

		return
	}

	notification, err := findNotification(data.Name.ValueString(), response)
	if err != nil {
		resp.Diagnostics.AddError(helpers.DataSourceError, fmt.Sprintf("Unable to find %s, got error: %s", notificationDataSourceName, err))

		return
	}

	tflog.Trace(ctx, "read "+notificationDataSourceName)
	data.write(ctx, notification)
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func findNotification(name string, notifications []*sonarr.NotificationOutput) (*sonarr.NotificationOutput, error) {
	for _, i := range notifications {
		if i.Name == name {
			return i, nil
		}
	}

	return nil, helpers.ErrDataNotFoundError(notificationDataSourceName, "name", name)
}

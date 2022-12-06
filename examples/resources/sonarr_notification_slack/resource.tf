resource "sonarr_notification_slack" "example" {
  on_grab                            = false
  on_download                        = true
  on_upgrade                         = true
  on_rename                          = false
  on_series_delete                   = false
  on_episode_file_delete             = false
  on_episode_file_delete_for_upgrade = true
  on_health_issue                    = false
  on_application_update              = false

  include_health_warnings = false
  name                    = "Example"

  web_hook_url = "http://my.slack.com/test"
  username     = "user"
  channel      = "example-channel"
}
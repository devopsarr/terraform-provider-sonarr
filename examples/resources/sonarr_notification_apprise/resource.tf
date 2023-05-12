resource "sonarr_notification_apprise" "example" {
  on_grab                            = false
  on_download                        = false
  on_upgrade                         = true
  on_series_delete                   = false
  on_episode_file_delete             = false
  on_episode_file_delete_for_upgrade = false
  on_health_issue                    = true
  on_application_update              = false

  include_health_warnings = false
  name                    = "Example"

  notification_type = 1
  server_url        = "https://apprise.go"
  auth_username     = "User"
  auth_password     = "Password"
  field_tags        = ["warning", "skull"]
}
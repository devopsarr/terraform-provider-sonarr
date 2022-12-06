resource "sonarr_notification_join" "example" {
  on_grab                            = false
  on_download                        = false
  on_upgrade                         = false
  on_series_delete                   = false
  on_episode_file_delete             = false
  on_episode_file_delete_for_upgrade = false
  on_health_issue                    = false
  on_application_update              = false

  include_health_warnings = false
  name                    = "Example"

  api_key  = "Key"
  priority = 2
}
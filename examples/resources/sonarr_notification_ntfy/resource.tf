resource "sonarr_notification_ntfy" "example" {
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

  priority   = 1
  server_url = "https://ntfy.sh"
  username   = "User"
  password   = "Pass"
  topics     = ["Topic1234", "Topic4321"]
  field_tags = ["warning", "skull"]
}
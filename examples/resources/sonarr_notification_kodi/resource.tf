resource "sonarr_notification_kodi" "example" {
  on_grab                            = false
  on_download                        = false
  on_upgrade                         = false
  on_rename                          = false
  on_series_delete                   = false
  on_episode_file_delete             = false
  on_episode_file_delete_for_upgrade = false
  on_health_issue                    = false
  on_application_update              = false

  include_health_warnings = false
  name                    = "Example"

  host     = "http://kodi.com"
  port     = 8080
  username = "User"
  password = "MyPass"
  notify   = true
}
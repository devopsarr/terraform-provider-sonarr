resource "sonarr_notification_boxcar" "example" {
  on_grab                            = false
  on_download                        = true
  on_upgrade                         = true
  on_series_delete                   = false
  on_episode_file_delete             = false
  on_episode_file_delete_for_upgrade = true
  on_health_issue                    = false
  on_application_update              = false

  include_health_warnings = false
  name                    = "Example"

  auth_username = "User"
  auth_password = "Password"

  host          = "localhost"
  port          = 8080
  use_ssl       = true
  sender_number = "1234"
  receiver_id   = "4321"
}
resource "sonarr_notification_mailgun" "example" {
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

  api_key    = "APIkey"
  from       = "from_mailgun@example.com"
  recipients = ["user1@example.com", "user2@example.com"]
}
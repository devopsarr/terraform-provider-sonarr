resource "sonarr_notification" "example" {
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

  implementation  = "CustomScript"
  config_contract = "CustomScriptSettings"

  path = "/scripts/sonarr.sh"
}
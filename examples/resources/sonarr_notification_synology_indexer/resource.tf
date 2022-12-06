resource "sonarr_notification_synology_indexer" "example" {
  on_download                        = true
  on_upgrade                         = true
  on_rename                          = false
  on_series_delete                   = false
  on_episode_file_delete             = false
  on_episode_file_delete_for_upgrade = true

  include_health_warnings = false
  name                    = "Example"

  update_library = true
}
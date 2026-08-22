resource "sonarr_notification_plex" "example" {
  on_download                        = true
  on_upgrade                         = true
  on_rename                          = false
  on_series_delete                   = false
  on_episode_file_delete             = false
  on_episode_file_delete_for_upgrade = true

  include_health_warnings = false
  name                    = "Example"

  host       = "plex.lcl"
  port       = 32400
  auth_token = "AuthTOKEN"

  # optional path mapping, for when Sonarr and Plex see the library at different
  # paths. only applied when update_library is enabled.
  update_library = true
  map_from       = "/tv"
  map_to         = "/media/tv"
}

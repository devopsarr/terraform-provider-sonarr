resource "sonarr_import_list_plex_rss" "example" {
  enable_automatic_add = true
  season_folder        = true
  should_monitor       = "all"
  series_type          = "standard"
  root_folder_path     = sonarr_root_folder.example.path
  quality_profile_id   = 1
  name                 = "Example"
  url                  = "https://app.plex.tv/desktop/#!/settings/watchlist"
}
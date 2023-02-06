resource "sonarr_import_list_user" "example" {
  name                 = "Example"
  enable_automatic_add = true
  season_folder        = true
  should_monitor       = "all"
  series_type          = "standard"
  root_folder_path     = sonarr_root_folder.example.path
  quality_profile_id   = 1
  username             = "User"
  access_token         = "Token"
  trakt_list_type      = 0
}
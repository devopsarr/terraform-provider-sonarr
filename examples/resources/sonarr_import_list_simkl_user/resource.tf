resource "sonarr_import_list_simkl_user" "example" {
  name                 = "Example"
  enable_automatic_add = true
  season_folder        = true
  should_monitor       = "all"
  series_type          = "standard"
  root_folder_path     = sonarr_root_folder.example.path
  quality_profile_id   = 1
  access_token         = "Token"
  list_type            = 0
}
resource "sonarr_import_list_sonarr" "example" {
  enable_automatic_add = true
  season_folder        = true
  should_monitor       = "all"
  series_type          = "standard"
  root_folder_path     = sonarr_root_folder.example.path
  quality_profile_id   = 1
  name                 = "Example"
  base_url             = "http://127.0.0.1:8989"
  api_key              = "APIKey"
  tags                 = [1, 2, 3]
  quality_profile_ids  = [1, 2]
  language_profile_ids = [1]
  tag_ids              = [1, 2, 3]
}
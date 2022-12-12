resource "sonarr_import_list" "example" {
  enable_automatic_add = true
  season_folder        = true
  should_monitor       = "all"
  series_type          = "standard"
  root_folder_path     = sonarr_root_folder.example.path
  quality_profile_id   = sonarr_quality_profile.example.id
  language_profile_id  = sonarr_language_profile.example.id
  name                 = "Esample"
  implementation       = "SonarrImport"
  config_contract      = "SonarrSettings"
  base_url             = "http://127.0.0.1:8989"
  api_key              = "b01df9fca2e64e459d64a09888ce7451"
  tags                 = [1, 2]
}
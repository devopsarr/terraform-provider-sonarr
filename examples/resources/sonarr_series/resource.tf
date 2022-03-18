resource "sonarr_series" "example" {
  title      = "Breaking Bad"
  title_slug = "breaking-bad"
  tvdb_id    = 81189

  monitored           = true
  season_folder       = true
  use_scene_numbering = false
  path                = "/tmp/breaking_bad"
  root_folder_path    = "/tmp/"

  language_profile_id = 1
  quality_profile_id  = 1
  tags                = [1]
}

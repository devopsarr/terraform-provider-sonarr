resource "sonarr_metadata_roksbox" "example" {
  enable           = true
  name             = "Example"
  episode_metadata = true
  series_images    = false
  season_images    = true
  episode_images   = false
}
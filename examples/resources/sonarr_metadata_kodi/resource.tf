resource "sonarr_metadata_kodi" "example" {
  enable              = true
  name                = "Example"
  series_metadata     = true
  series_images       = true
  episode_images      = true
  series_metadata_url = false
  season_images       = true
  episode_metadata    = false
}
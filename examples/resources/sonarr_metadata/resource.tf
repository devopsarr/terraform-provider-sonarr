resource "sonarr_metadata" "example" {
  enable           = true
  name             = "Example"
  implementation   = "MediaBrowserMetadata"
  config_contract  = "MediaBrowserMetadataSettings"
  episode_metadata = true
  series_images    = false
  season_images    = true
  episode_images   = false
  tags             = [1, 2]
}
resource "sonarr_download_client_torrent_download_station" "example" {
  enable   = true
  priority = 1
  name     = "Example"
  host     = "downloadstation"
  port     = 5000
}
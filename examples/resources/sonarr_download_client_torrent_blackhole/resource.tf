resource "sonarr_download_client_torrent_blackhole" "example" {
  enable                = true
  priority              = 1
  name                  = "Example"
  magnet_file_extension = ".magnet"
  watch_folder          = "/watch/"
  torrent_folder        = "/torrent/"
}
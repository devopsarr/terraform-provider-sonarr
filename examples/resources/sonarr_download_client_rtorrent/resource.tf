resource "sonarr_download_client_rtorrent" "example" {
  enable   = true
  priority = 1
  name     = "Example"
  host     = "rtorrent"
  url_base = "/rtorrent/"
  port     = 9091
}
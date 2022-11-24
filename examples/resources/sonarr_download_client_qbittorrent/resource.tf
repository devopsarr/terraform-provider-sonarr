resource "sonarr_download_client_qbittorrent" "example" {
  enable         = true
  priority       = 1
  name           = "Example"
  host           = "qbittorrent"
  url_base       = "/qbittorrent/"
  port           = 9091
  tv_category    = "tv-sonarr"
  first_and_last = true
}
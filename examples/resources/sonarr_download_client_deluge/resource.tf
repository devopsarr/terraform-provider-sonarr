resource "sonarr_download_client_deluge" "example" {
  enable   = true
  priority = 1
  name     = "Example"
  host     = "deluge"
  url_base = "/deluge/"
  port     = 9091
}
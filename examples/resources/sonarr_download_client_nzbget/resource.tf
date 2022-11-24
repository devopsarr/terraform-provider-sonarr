resource "sonarr_download_client_nzbget" "example" {
  enable   = true
  priority = 1
  name     = "Example"
  host     = "nzbget"
  url_base = "/nzbget/"
  port     = 6789
}
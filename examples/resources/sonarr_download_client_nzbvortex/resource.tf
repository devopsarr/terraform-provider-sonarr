resource "sonarr_download_client_nzbvortex" "example" {
  enable   = true
  priority = 1
  name     = "Example"
  host     = "nzbvortex"
  url_base = "/nzbvortex/"
  port     = 6789
}
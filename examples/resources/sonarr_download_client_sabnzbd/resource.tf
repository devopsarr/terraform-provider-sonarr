resource "sonarr_download_client_sabnzbd" "example" {
  enable   = true
  priority = 1
  name     = "Example"
  host     = "sabnzbd"
  url_base = "/sabnzbd/"
  port     = 9091
  api_key  = "test"
}
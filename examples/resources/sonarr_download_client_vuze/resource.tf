resource "sonarr_download_client_vuze" "example" {
  enable   = true
  priority = 1
  name     = "Example"
  host     = "vuze"
  url_base = "/vuze/"
  port     = 9091
}
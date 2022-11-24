resource "sonarr_download_client_hadouken" "example" {
  enable   = true
  priority = 1
  name     = "Example"
  host     = "hadouken"
  url_base = "/hadouken/"
  port     = 9091
  username = "username"
  password = "password"
}
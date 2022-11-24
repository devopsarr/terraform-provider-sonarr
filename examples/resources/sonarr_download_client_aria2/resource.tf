resource "sonarr_download_client_aria2" "example" {
  enable   = true
  priority = 1
  name     = "Example"
  host     = "aria2"
  rpc_path = "/aria2/"
  port     = 6800
}
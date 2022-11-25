resource "sonarr_download_client_usenet_blackhole" "example" {
  enable       = true
  priority     = 1
  name         = "Example"
  watch_folder = "/watch/"
  nzb_folder   = "/nzb/"
}
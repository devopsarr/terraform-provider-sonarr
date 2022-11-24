resource "sonarr_download_client_pneumatic" "example" {
  enable      = true
  priority    = 1
  name        = "Example"
  nzb_folder  = "/nzb/"
  strm_folder = "/strm/"
}
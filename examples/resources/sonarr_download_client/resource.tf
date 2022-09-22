resource "sonarr_download_client" "example" {
  enable          = true
  priority        = 1
  name            = "Example"
  implementation  = "Transmission"
  protocol        = "torrent"
  config_contract = "TransmissionSettings"
  host            = "transmission"
  url_base        = "/transmission/"
  port            = 9091
}
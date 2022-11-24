resource "sonarr_download_client_flood" "example" {
  enable          = true
  priority        = 1
  name            = "Example"
  host            = "flood"
  url_base        = "/flood/"
  port            = 9091
  start_on_add    = true
  additional_tags = [0, 1]
  field_tags      = ["sonarr"]
}
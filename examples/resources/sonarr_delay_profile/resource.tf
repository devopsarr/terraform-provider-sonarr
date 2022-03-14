resource "sonarr_delay_profile" "example" {
  enable_usenet             = true
  enable_torrent            = true
  bypass_if_highest_quality = true
  usenet_delay              = 0
  torrent_delay             = 0
  tags                      = [1, 2]
  preferred_protocol        = "torrent"
}
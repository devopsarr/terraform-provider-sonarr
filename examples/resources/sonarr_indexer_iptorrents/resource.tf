resource "sonarr_indexer_iptorrents" "example" {
  enable_automatic_search = true
  name                    = "Example"
  base_url                = "https://iptorrent.io"
  minimum_seeders         = 1
}

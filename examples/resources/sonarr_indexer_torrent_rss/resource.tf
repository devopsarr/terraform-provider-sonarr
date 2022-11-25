resource "sonarr_indexer_torrent_rss" "example" {
  enable_automatic_search = true
  name                    = "Example"
  base_url                = "https://rss.io"
  allow_zero_size         = true
  minimum_seeders         = 1
}

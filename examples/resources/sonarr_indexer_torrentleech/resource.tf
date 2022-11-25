resource "sonarr_indexer_hdbits" "example" {
  enable_automatic_search = true
  name                    = "Example"
  base_url                = "http://rss.torrentleech.org"
  api_key                 = "APIKey"
  minimum_seeders         = 1
}

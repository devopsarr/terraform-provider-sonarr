resource "sonarr_indexer_fanzub" "example" {
  enable_automatic_search = true
  name                    = "Example"
  base_url                = "http://fanzub.com/rss/"
}
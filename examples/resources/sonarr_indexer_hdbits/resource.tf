resource "sonarr_indexer_hdbits" "example" {
  enable_automatic_search = true
  name                    = "Example"
  base_url                = "https://hdbits.org"
  username                = "User"
  api_key                 = "APIKey"
  minimum_seeders         = 1
}

resource "sonarr_indexer_rargb" "example" {
  enable_automatic_search = true
  name                    = "Example"
  base_url                = "https://api.broadcasthe.net/"
  api_key                 = "APIKey"
}

resource "sonarr_indexer_torznab" "example" {
  enable_automatic_search = true
  name                    = "Example"
  base_url                = "https://feed.animetosho.org"
  api_path                = "/nabapi"
  anime_categories        = [5070]
  minimum_seeders         = 1
}

resource "sonarr_indexer_newznab" "example" {
  enable_automatic_search = true
  name                    = "Example"
  base_url                = "https://lolo.sickbeard.com"
  api_path                = "/api"
  categories              = [5030, 5040]
}
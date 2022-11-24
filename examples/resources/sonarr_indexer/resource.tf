resource "sonarr_indexer" "example" {
  enable_automatic_search = true
  name                    = "Example"
  implementation          = "Newznab"
  protocol                = "usenet"
  config_contract         = "NewznabSettings"
  base_url                = "https://lolo.sickbeard.com"
  api_path                = "/api"
  categories              = [5030, 5040]
  tags                    = [1, 2]
}
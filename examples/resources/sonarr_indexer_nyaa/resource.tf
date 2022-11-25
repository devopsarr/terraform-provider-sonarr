resource "sonarr_indexer_nyaa" "example" {
  enable_automatic_search = true
  name                    = "Example"
  base_url                = "https://nyaa.io"
  additional_parameters   = "&cats=1_0&filter=1"
  minimum_seeders         = 1
}

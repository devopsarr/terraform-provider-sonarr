resource "sonarr_indexer_rargb" "example" {
  enable_automatic_search = true
  name                    = "Example"
  implementation          = "Newznab"
  base_url                = "https://torrentapi.org"
  ranked_only             = "false"
  minimum_seeders         = 1
}

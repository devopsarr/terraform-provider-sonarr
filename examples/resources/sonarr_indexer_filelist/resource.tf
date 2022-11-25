resource "sonarr_indexer_filelist" "example" {
  enable_automatic_search = true
  name                    = "Example"
  base_url                = "https://filelist.io"
  categories              = [21, 23, 27]
  username                = "User"
  passkey                 = "PassKey"
  minimum_seeders         = 1
}

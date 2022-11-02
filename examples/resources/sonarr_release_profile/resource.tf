resource "sonarr_release_profile" "example" {
  enabled    = true
  name       = "Example"
  required   = ["proper"]
  ignored    = ["repack"]
  indexer_id = 1
}

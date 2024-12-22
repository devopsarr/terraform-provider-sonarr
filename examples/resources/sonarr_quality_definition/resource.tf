resource "sonarr_quality_definition" "example" {
  id             = 21
  title          = "Example"
  min_size       = 35.0
  max_size       = 400
  preferred_size = 200
}

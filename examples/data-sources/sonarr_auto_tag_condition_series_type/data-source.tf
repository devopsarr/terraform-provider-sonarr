data "sonarr_auto_tag_condition_series_type" "example" {
  name     = "Example"
  negate   = false
  required = false
  value    = "1"
}

resource "sonarr_custom_format" "example" {
  remove_tags_automatically = false
  name                      = "Example"

  tags = [1, 2]

  specifications = [data.sonarr_auto_tag_condition_series_type.example]
}
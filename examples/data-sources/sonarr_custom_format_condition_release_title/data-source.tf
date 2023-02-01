data "sonarr_custom_format_condition_release_title" "example" {
  name     = "x265"
  negate   = false
  required = false
  value    = "(((x|h)\\.?265)|(HEVC))"
}

resource "sonarr_custom_format" "example" {
  include_custom_format_when_renaming = false
  name                                = "Example"

  specifications = [data.sonarr_custom_format_condition_release_title.example]
}
data "sonarr_custom_format_condition_release_group" "example" {
  name     = "HDBits"
  negate   = false
  required = false
  value    = ".*HDBits.*"
}

resource "sonarr_custom_format" "example" {
  include_custom_format_when_renaming = false
  name                                = "Example"

  specifications = [data.sonarr_custom_format_condition_release_group.example]
}
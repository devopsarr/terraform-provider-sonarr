data "sonarr_auto_tag_condition_root_folder" "example" {
  name     = "Example"
  negate   = false
  required = false
  value    = "/series"
}

resource "sonarr_custom_format" "example" {
  remove_tags_automatically = false
  name                      = "Example"

  tags = [1, 2]

  specifications = [data.sonarr_auto_tag_condition_root_folder.example]
}
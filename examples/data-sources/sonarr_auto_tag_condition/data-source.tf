data "sonarr_auto_tag_condition" "example" {
  name           = "Example"
  implementation = "SeriesTypeSpecification"
  negate         = false
  required       = false
  value          = "2"
}

resource "sonarr_auto_tag" "example" {
  remove_tags_automatically = false
  name                      = "Example"

  tags = [1, 2]

  specifications = [data.sonarr_auto_tag.example]
}
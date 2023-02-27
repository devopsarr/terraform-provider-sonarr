resource "sonarr_auto_tag" "example" {
  name                      = "Example"
  remove_tags_automatically = true
  tags                      = [1, 2]

  specifications = [
    {
      name           = "folder"
      implementation = "RootFolderSpecification"
      negate         = true
      required       = false
      value          = "/series"
    },
    {
      name           = "type"
      implementation = "SeriesTypeSpecification"
      negate         = true
      required       = false
      value          = "2"
    },
    {
      name           = "genre"
      implementation = "GenreSpecification"
      negate         = false
      required       = false
      value          = "horror comedy"
    },
  ]
}
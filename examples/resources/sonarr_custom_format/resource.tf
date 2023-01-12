resource "sonarr_custom_format" "example" {
  include_custom_format_when_renaming = true
  name                                = "Example"

  specifications = [
    {
      name           = "Surround Sound"
      implementation = "ReleaseTitleSpecification"
      negate         = false
      required       = false
      value          = "DTS.?(HD|ES|X(?!\\D))|TRUEHD|ATMOS|DD(\\+|P).?([5-9])|EAC3.?([5-9])"
    },
    {
      name           = "Arabic"
      implementation = "LanguageSpecification"
      negate         = false
      required       = false
      value          = "31"
    }
  ]
}
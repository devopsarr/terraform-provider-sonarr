resource "sonarr_language_profile" "example" {
  upgrade_allowed = true
  name            = "Eng"
  cutoff_language = "English"
  languages       = ["English", "Italian"]
}
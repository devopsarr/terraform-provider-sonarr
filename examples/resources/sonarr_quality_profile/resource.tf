resource "sonarr_quality_profile" "example" {
  name            = "example-4k"
  upgrade_allowed = true
  cutoff          = 1100

  quality_groups = [
    {
      id   = 1100
      name = "4k"
      qualities = [
        {
          id         = 18
          name       = "WEBDL-2160p"
          source     = "web"
          resolution = 2160
        },
        {
          id         = 19
          name       = "Bluray-2160p"
          source     = "bluray"
          resolution = 2160
        }
      ]
    }
  ]
}
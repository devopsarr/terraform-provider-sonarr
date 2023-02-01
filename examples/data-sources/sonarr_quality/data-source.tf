data "sonarr_quality" "bluray" {
  name = "Bluray-2160p"
}

data "sonarr_quality" "webdl" {
  name = "WEBDL-2160p"
}

data "sonarr_quality" "webrip" {
  name = "WEBRip-2160p"
}

resource "sonarr_quality_profile" "Example" {
  name            = "Example"
  upgrade_allowed = true
  cutoff          = 2000

  quality_groups = [
    {
      id   = 2000
      name = "WEB 2160p"
      qualities = [
        data.sonarr_quality.webdl,
        data.sonarr_quality.webrip,
      ]
    },
    {
      qualities = [data.sonarr_quality.bluray]
    }
  ]
}
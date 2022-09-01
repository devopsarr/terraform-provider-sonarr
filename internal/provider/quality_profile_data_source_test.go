package provider

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccQualityProfileDataSource(t *testing.T) {
	t.Parallel()

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Read testing
			{
				Config: testAccQualityProfileDataSourceConfig,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.sonarr_quality_profile.test", "id"),
					resource.TestCheckResourceAttr("data.sonarr_quality_profile.test", "cutoff", "1100")),
			},
		},
	})
}

const testAccQualityProfileDataSourceConfig = `
resource "sonarr_quality_profile" "test" {
	name            = "qpdata"
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

data "sonarr_quality_profile" "test" {
	name = sonarr_quality_profile.test.name
}
`

package provider

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccDelayProfilesDataSource(t *testing.T) {
	t.Parallel()

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Read testing
			{
				Config: testAccDelayProfilesDataSourceConfig,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckTypeSetElemNestedAttrs("data.sonarr_delay_profiles.test", "delay_profiles.*", map[string]string{"preferred_protocol": "torrent"}),
				),
			},
		},
	})
}

const testAccDelayProfilesDataSourceConfig = `
resource "sonarr_tag" "test" {
	label = "delay_profiles_datasource"
}

resource "sonarr_delay_profile" "test" {
	enable_usenet = true
	enable_torrent = true
	bypass_if_highest_quality = true
	usenet_delay = 0
	torrent_delay = 0
	preferred_protocol= "torrent"
	tags = [sonarr_tag.test.id]
}

data "sonarr_delay_profiles" "test" {
	depends_on = [sonarr_delay_profile.test]
}
`

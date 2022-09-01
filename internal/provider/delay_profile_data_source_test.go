package provider

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccDelayProfileDataSource(t *testing.T) {
	t.Parallel()

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Read testing
			{
				Config: testAccDelayProfileDataSourceConfig,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.sonarr_delay_profile.test", "id"),
					resource.TestCheckResourceAttr("data.sonarr_delay_profile.test", "preferred_protocol", "torrent")),
			},
		},
	})
}

const testAccDelayProfileDataSourceConfig = `
resource "sonarr_tag" "test" {
	label = "dpdata"
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

data "sonarr_delay_profile" "test" {
	id = sonarr_delay_profile.test.id
}
`

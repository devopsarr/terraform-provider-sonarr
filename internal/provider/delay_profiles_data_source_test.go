package provider

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccDelayProfilesDataSource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create a delay profile to have a value to check
			{
				Config: testAccDelayProfileResourceConfig("torrent"),
			},
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
data "sonarr_delay_profiles" "test" {
}
`

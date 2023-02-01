package provider

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccMetadataDataSource(t *testing.T) {
	t.Parallel()

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Read testing
			{
				Config: testAccMetadataDataSourceConfig,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.sonarr_metadata.test", "id"),
					resource.TestCheckResourceAttr("data.sonarr_metadata.test", "episode_metadata", "false")),
			},
		},
	})
}

const testAccMetadataDataSourceConfig = `
resource "sonarr_metadata" "test" {
	enable = true
	name = "metadataData"
	implementation = "WdtvMetadata"
	config_contract = "WdtvMetadataSettings"
	episode_metadata = false
	series_images = false
	season_images = true
	episode_images = false
}

data "sonarr_metadata" "test" {
	name = sonarr_metadata.test.name
}
`

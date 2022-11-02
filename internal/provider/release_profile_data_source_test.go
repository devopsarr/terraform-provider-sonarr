package provider

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccReleaseProfileDataSource(t *testing.T) {
	t.Parallel()

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Read testing
			{
				Config: testAccReleaseProfileDataSourceConfig,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.sonarr_release_profile.test", "id"),
					resource.TestCheckResourceAttr("data.sonarr_release_profile.test", "name", "dataSourceTestSingle")),
			},
		},
	})
}

const testAccReleaseProfileDataSourceConfig = `
resource "sonarr_release_profile" "test" {
	name = "dataSourceTestSingle"
	indexer_id = 0
	required= ["notreally"]
}

data "sonarr_release_profile" "test" {
	id = sonarr_release_profile.test.id
}
`

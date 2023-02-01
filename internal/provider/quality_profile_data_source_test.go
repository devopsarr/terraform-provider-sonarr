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
					resource.TestCheckResourceAttr("data.sonarr_quality_profile.test", "cutoff", "1")),
			},
		},
	})
}

const testAccQualityProfileDataSourceConfig = `
data "sonarr_quality_profile" "test" {
	name = "Any"
}
`

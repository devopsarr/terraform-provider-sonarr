package provider

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccQualityDataSource(t *testing.T) {
	t.Parallel()

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Read testing
			{
				Config: testAccQualityDataSourceConfig,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.sonarr_quality.test", "id"),
					resource.TestCheckResourceAttr("data.sonarr_quality.test", "resolution", "2160")),
			},
		},
	})
}

const testAccQualityDataSourceConfig = `
data "sonarr_quality" "test" {
	name = "Bluray-2160p"
}
`

package provider

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccSeriesDataSource(t *testing.T) {
	t.Parallel()

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create a series to have a value to check
			{
				Config: testAccSeriesResourceConfig(332606, "Friends", "friends", "false"),
			},
			// Read testing
			{
				Config: testAccSeriesDataSourceConfig,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckTypeSetElemNestedAttrs("data.sonarr_series.test", "series.*", map[string]string{"monitored": "false"}),
				),
			},
		},
	})
}

const testAccSeriesDataSourceConfig = `
data "sonarr_series" "test" {
}
`

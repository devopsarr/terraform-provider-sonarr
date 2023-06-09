package provider

import (
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccAllSeriesDataSource(t *testing.T) {
	t.Parallel()

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Unauthorized
			{
				Config:      testAccAllSeriesDataSourceConfig + testUnauthorizedProvider,
				ExpectError: regexp.MustCompile("Client Error"),
			},
			// Create a resource to test
			{
				Config: testAccSeriesResourceConfig(332606, "Friends", "friends", "false"),
			},
			// Read testing
			{
				Config: testAccAllSeriesDataSourceConfig,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckTypeSetElemNestedAttrs("data.sonarr_all_series.test", "series.*", map[string]string{"monitored": "false"}),
				),
			},
		},
	})
}

const testAccAllSeriesDataSourceConfig = `
data "sonarr_all_series" "test" {
}
`

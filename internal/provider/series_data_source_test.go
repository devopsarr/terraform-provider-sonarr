package provider

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccSeriesDataSource(t *testing.T) {
	t.Parallel()

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// // Unauthorized
			{
				Config:      testAccSeriesDataSourceConfig("\"Error\"") + testUnauthorizedProvider,
				ExpectError: regexp.MustCompile("Client Error"),
			},
			// Not found testing
			{
				Config:      testAccSeriesDataSourceConfig("\"Error\""),
				ExpectError: regexp.MustCompile("Unable to find series"),
			},
			// Read testing
			{
				Config: testAccSeriesResourceConfig(153021, "The Walking Dead", "the-walking-dead", "false") + testAccSeriesDataSourceConfig("sonarr_series.test.title"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.sonarr_series.test", "id"),
					resource.TestCheckResourceAttr("data.sonarr_series.test", "path", "/config/the-walking-dead")),
			},
		},
	})
}

func testAccSeriesDataSourceConfig(title string) string {
	return fmt.Sprintf(`
	data "sonarr_series" "test" {
		title = %s
	}
	`, title)
}

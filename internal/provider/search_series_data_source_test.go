package provider

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccSearchSeriesDataSource(t *testing.T) {
	t.Parallel()

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Read testing
			{
				Config: testAccSearchSeriesDataSourceConfig,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.sonarr_search_series.test", "id"),
					resource.TestCheckResourceAttr("data.sonarr_search_series.test", "tvdb_id", "153021")),
			},
		},
	})
}

const testAccSearchSeriesDataSourceConfig = `
data "sonarr_search_series" "test" {
	tvdb_id    = 153021
}
`

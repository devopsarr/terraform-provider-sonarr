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
			// Read testing
			{
				Config: testAccSeriesDataSourceConfig,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.sonarr_series.test", "id"),
					resource.TestCheckResourceAttr("data.sonarr_series.test", "path", "/tmp/the-walking-dead")),
			},
		},
	})
}

const testAccSeriesDataSourceConfig = `
resource "sonarr_series" "test" {
	title      = "The Walking Dead"
	title_slug = "the-walking-dead"
	tvdb_id    = 153021
  
	monitored           = false
	season_folder       = true
	use_scene_numbering = false
	path                = "/tmp/the-walking-dead"
	root_folder_path    = "/tmp"
  
	language_profile_id = 1
	quality_profile_id  = 1
	tags                = []
}

data "sonarr_series" "test" {
	title = sonarr_series.test.title
}
`

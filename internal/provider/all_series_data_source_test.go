package provider

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccAllSeriesDataSource(t *testing.T) {
	t.Parallel()

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
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
resource "sonarr_tag" "test" {
	label = "all_series_datasource"
}

resource "sonarr_series" "test" {
	title      = "Friends"
	title_slug = "friends"
	tvdb_id    = 332606
  
	monitored           = false
	season_folder       = true
	use_scene_numbering = false
	path                = "/tmp/friends"
	root_folder_path    = "/tmp"
  
	language_profile_id = 1
	quality_profile_id  = 1
	tags                = [sonarr_tag.test.id]
}

data "sonarr_all_series" "test" {
	depends_on = [sonarr_series.test]
}
`

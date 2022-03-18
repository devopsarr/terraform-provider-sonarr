package provider

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccSeriesResource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Init a Tag to be there for testing
			{
				Config: testAccTagResourceConfig("test", "eng"),
			},
			// Create and Read testing
			{
				Config: testAccSeriesResourceConfig("true"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("sonarr_series.test", "monitored", "true"),
					resource.TestCheckResourceAttrSet("sonarr_series.test", "id"),
				),
			},
			// Update and Read testing
			{
				Config: testAccSeriesResourceConfig("false"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("sonarr_series.test", "monitored", "false"),
				),
			},
			// ImportState testing
			{
				ResourceName:      "sonarr_series.test",
				ImportState:       true,
				ImportStateVerify: true,
			},
			// Delete testing automatically occurs in TestCase
		},
	})
}

func testAccSeriesResourceConfig(monitored string) string {
	return fmt.Sprintf(`
	resource "sonarr_series" "test" {
		title      = "Breaking Bad"
		title_slug = "breaking-bad"
		tvdb_id    = 81189
	  
		monitored           = %s
		season_folder       = true
		use_scene_numbering = false
		path                = "/tmp/breaking_bad"
		root_folder_path    = "/tmp"
	  
		language_profile_id = 1
		quality_profile_id  = 1
		tags                = [1]
	}
	`, monitored)
}

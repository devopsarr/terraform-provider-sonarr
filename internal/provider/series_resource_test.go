package provider

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccSeriesResource(t *testing.T) {
	t.Parallel()

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Unauthorized Create
			{
				Config:      testAccSeriesResourceConfig(81189, "Breaking Bad", "breaking-bad", "false") + testUnauthorizedProvider,
				ExpectError: regexp.MustCompile("Client Error"),
			},
			// Create and Read testing
			{
				Config: testAccSeriesResourceConfig(81189, "Breaking Bad", "breaking-bad", "false"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("sonarr_series.test", "monitored", "false"),
					resource.TestCheckResourceAttrSet("sonarr_series.test", "id"),
				),
			},
			// Unauthorized Read
			{
				Config:      testAccSeriesResourceConfig(81189, "Breaking Bad", "breaking-bad", "false") + testUnauthorizedProvider,
				ExpectError: regexp.MustCompile("Client Error"),
			},
			// Update and Read testing
			{
				Config: testAccSeriesResourceConfig(81189, "Breaking Bad", "breaking-bad", "true"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("sonarr_series.test", "monitored", "true"),
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

func testAccSeriesResourceConfig(id int, title, slug, monitored string) string {
	return fmt.Sprintf(`
	resource "sonarr_series" "test" {
		title      = "%s"
		title_slug = "%s"
		tvdb_id    = %d
	  
		monitored           = %s
		season_folder       = true
		use_scene_numbering = false
		path                = "/config/%s"
		root_folder_path    = "/config"
	  
		quality_profile_id  = 1
	}
	`, title, slug, id, monitored, slug)
}

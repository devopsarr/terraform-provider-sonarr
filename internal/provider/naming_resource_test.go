package provider

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccNamingResource(t *testing.T) {
	t.Parallel()

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Unauthorized Create
			{
				Config:      testAccNamingResourceConfig("Specials") + testUnauthorizedProvider,
				ExpectError: regexp.MustCompile("Client Error"),
			},
			// Create and Read testing
			{
				Config: testAccNamingResourceConfig("Specials"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("sonarr_naming.test", "specials_folder_format", "Specials"),
					resource.TestCheckResourceAttrSet("sonarr_naming.test", "id"),
				),
			},
			// Unauthorized Read
			{
				Config:      testAccNamingResourceConfig("Specials") + testUnauthorizedProvider,
				ExpectError: regexp.MustCompile("Client Error"),
			},
			// Update and Read testing
			{
				Config: testAccNamingResourceConfig("S0"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("sonarr_naming.test", "specials_folder_format", "S0"),
				),
			},
			// ImportState testing
			{
				ResourceName:      "sonarr_naming.test",
				ImportState:       true,
				ImportStateVerify: true,
			},
			// Delete testing automatically occurs in TestCase
		},
	})
}

func testAccNamingResourceConfig(specials string) string {
	return fmt.Sprintf(`
	resource "sonarr_naming" "test" {
		rename_episodes            = true
		replace_illegal_characters = true
		multi_episode_style        = 0
		colon_replacement_format   = 4
		daily_episode_format       = "{Series Title} - {Air-Date} - {Episode Title} {Quality Full}"
		anime_episode_format       = "{Series Title} - S{season:00}E{episode:00} - {Episode Title} {Quality Full}"
		series_folder_format       = "{Series Title}"
		season_folder_format       = "Season {season}"
		specials_folder_format     = "%s"
		standard_episode_format    = "{Series Title} - S{season:00}E{episode:00} - {Episode Title} {Quality Full}"
	}`, specials)
}

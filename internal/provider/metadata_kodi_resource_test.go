package provider

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccMetadataKodiResource(t *testing.T) {
	t.Parallel()

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Unauthorized Create
			{
				Config:      testAccMetadataKodiResourceConfig("kodiResourceTest", "false") + testUnauthorizedProvider,
				ExpectError: regexp.MustCompile("Client Error"),
			},
			// Create and Read testing
			{
				Config: testAccMetadataKodiResourceConfig("kodiResourceTest", "false"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("sonarr_metadata_kodi.test", "series_metadata", "false"),
					resource.TestCheckResourceAttrSet("sonarr_metadata_kodi.test", "id"),
				),
			},
			// Unauthorized Read
			{
				Config:      testAccMetadataKodiResourceConfig("kodiResourceTest", "false") + testUnauthorizedProvider,
				ExpectError: regexp.MustCompile("Client Error"),
			},
			// Update and Read testing
			{
				Config: testAccMetadataKodiResourceConfig("kodiResourceTest", "true"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("sonarr_metadata_kodi.test", "series_metadata", "true"),
				),
			},
			// ImportState testing
			{
				ResourceName:      "sonarr_metadata_kodi.test",
				ImportState:       true,
				ImportStateVerify: true,
			},
			// Delete testing automatically occurs in TestCase
		},
	})
}

func testAccMetadataKodiResourceConfig(name, metadata string) string {
	return fmt.Sprintf(`
	resource "sonarr_metadata_kodi" "test" {
		enable = false
		name = "%s"
		series_metadata = %s
		series_images = true
		episode_images = true
		series_metadata_url = false
		season_images = true
		episode_metadata = false
	}`, name, metadata)
}

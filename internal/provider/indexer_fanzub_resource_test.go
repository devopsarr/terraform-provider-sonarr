package provider

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccIndexerFanzubResource(t *testing.T) {
	t.Parallel()

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Unauthorized Create
			{
				Config:      testAccIndexerFanzubResourceConfig("fanzubResourceTest", "false") + testUnauthorizedProvider,
				ExpectError: regexp.MustCompile("Client Error"),
			},
			// Create and Read testing
			{
				Config: testAccIndexerFanzubResourceConfig("fanzubResourceTest", "false"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("sonarr_indexer_fanzub.test", "anime_standard_format_search", "false"),
					resource.TestCheckResourceAttr("sonarr_indexer_fanzub.test", "base_url", "http://fanzub.com/rss/"),
					resource.TestCheckResourceAttrSet("sonarr_indexer_fanzub.test", "id"),
				),
			},
			// Unauthorized Read
			{
				Config:      testAccIndexerFanzubResourceConfig("fanzubResourceTest", "false") + testUnauthorizedProvider,
				ExpectError: regexp.MustCompile("Client Error"),
			},
			// Update and Read testing
			{
				Config: testAccIndexerFanzubResourceConfig("fanzubResourceTest", "true"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("sonarr_indexer_fanzub.test", "anime_standard_format_search", "true"),
				),
			},
			// ImportState testing
			{
				ResourceName:      "sonarr_indexer_fanzub.test",
				ImportState:       true,
				ImportStateVerify: true,
			},
			// Delete testing automatically occurs in TestCase
		},
	})
}

func testAccIndexerFanzubResourceConfig(name, aSearch string) string {
	return fmt.Sprintf(`
	resource "sonarr_indexer_fanzub" "test" {
		enable_automatic_search = false
		name = "%s"
		base_url = "http://fanzub.com/rss/"
		anime_standard_format_search = %s
	}`, name, aSearch)
}

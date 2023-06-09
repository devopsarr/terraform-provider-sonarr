package provider

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccIndexerNyaaResource(t *testing.T) {
	t.Parallel()

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Unauthorized Create
			{
				Config:      testAccIndexerNyaaResourceConfig("nyaaResourceTest", "https://nyaa.org") + testUnauthorizedProvider,
				ExpectError: regexp.MustCompile("Client Error"),
			},
			// Create and Read testing
			{
				Config: testAccIndexerNyaaResourceConfig("nyaaResourceTest", "https://nyaa.org"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("sonarr_indexer_nyaa.test", "base_url", "https://nyaa.org"),
					resource.TestCheckResourceAttrSet("sonarr_indexer_nyaa.test", "id"),
				),
			},
			// Unauthorized Read
			{
				Config:      testAccIndexerNyaaResourceConfig("nyaaResourceTest", "https://nyaa.org") + testUnauthorizedProvider,
				ExpectError: regexp.MustCompile("Client Error"),
			},
			// Update and Read testing
			{
				Config: testAccIndexerNyaaResourceConfig("nyaaResourceTest", "https://nyaa.net"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("sonarr_indexer_nyaa.test", "base_url", "https://nyaa.net"),
				),
			},
			// ImportState testing
			{
				ResourceName:      "sonarr_indexer_nyaa.test",
				ImportState:       true,
				ImportStateVerify: true,
			},
			// Delete testing automatically occurs in TestCase
		},
	})
}

func testAccIndexerNyaaResourceConfig(name, url string) string {
	return fmt.Sprintf(`
	resource "sonarr_indexer_nyaa" "test" {
		enable_automatic_search = false
		name = "%s"
		base_url = "%s"
		minimum_seeders = 1
	}`, name, url)
}

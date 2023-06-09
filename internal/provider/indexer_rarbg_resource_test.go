package provider

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccIndexerRarbgResource(t *testing.T) {
	t.Parallel()

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Unauthorized Create
			{
				Config:      testAccIndexerRarbgResourceConfig("rarbgResourceTest", "false") + testUnauthorizedProvider,
				ExpectError: regexp.MustCompile("Client Error"),
			},
			// Create and Read testing
			{
				Config: testAccIndexerRarbgResourceConfig("rarbgResourceTest", "false"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("sonarr_indexer_rarbg.test", "ranked_only", "false"),
					resource.TestCheckResourceAttr("sonarr_indexer_rarbg.test", "base_url", "https://torrentapi.org"),
					resource.TestCheckResourceAttrSet("sonarr_indexer_rarbg.test", "id"),
				),
			},
			// Unauthorized Read
			{
				Config:      testAccIndexerRarbgResourceConfig("rarbgResourceTest", "false") + testUnauthorizedProvider,
				ExpectError: regexp.MustCompile("Client Error"),
			},
			// Update and Read testing
			{
				Config: testAccIndexerRarbgResourceConfig("rarbgResourceTest", "true"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("sonarr_indexer_rarbg.test", "ranked_only", "true"),
				),
			},
			// ImportState testing
			{
				ResourceName:      "sonarr_indexer_rarbg.test",
				ImportState:       true,
				ImportStateVerify: true,
			},
			// Delete testing automatically occurs in TestCase
		},
	})
}

func testAccIndexerRarbgResourceConfig(name, ranked string) string {
	return fmt.Sprintf(`
	resource "sonarr_indexer_rarbg" "test" {
		enable_automatic_search = false
		name = "%s"
		base_url = "https://torrentapi.org"
		ranked_only = %s
		minimum_seeders = 1
	}`, name, ranked)
}

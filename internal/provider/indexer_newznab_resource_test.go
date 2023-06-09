package provider

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccIndexerNewznabResource(t *testing.T) {
	t.Parallel()

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Unauthorized Create
			{
				Config:      testAccIndexerNewznabResourceConfig("newzabResourceTest", "false") + testUnauthorizedProvider,
				ExpectError: regexp.MustCompile("Client Error"),
			},
			// Create and Read testing
			{
				Config: testAccIndexerNewznabResourceConfig("newzabResourceTest", "false"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("sonarr_indexer_newznab.test", "enable_automatic_search", "false"),
					resource.TestCheckResourceAttr("sonarr_indexer_newznab.test", "base_url", "https://lolo.sickbeard.com"),
					resource.TestCheckResourceAttrSet("sonarr_indexer_newznab.test", "id"),
				),
			},
			// Unauthorized Read
			{
				Config:      testAccIndexerNewznabResourceConfig("newzabResourceTest", "false") + testUnauthorizedProvider,
				ExpectError: regexp.MustCompile("Client Error"),
			},
			// Update and Read testing
			{
				Config: testAccIndexerNewznabResourceConfig("newzabResourceTest", "true"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("sonarr_indexer_newznab.test", "enable_automatic_search", "true"),
				),
			},
			// ImportState testing
			{
				ResourceName:      "sonarr_indexer_newznab.test",
				ImportState:       true,
				ImportStateVerify: true,
			},
			// Delete testing automatically occurs in TestCase
		},
	})
}

func testAccIndexerNewznabResourceConfig(name, aSearch string) string {
	return fmt.Sprintf(`
	resource "sonarr_indexer_newznab" "test" {
		enable_automatic_search = %s
		name = "%s"
		base_url = "https://lolo.sickbeard.com"
		api_path = "/api"
		categories = [5030, 5040]
	}`, aSearch, name)
}

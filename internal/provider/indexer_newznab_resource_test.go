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
				Config:      testAccIndexerNewznabResourceConfig("newzabResourceTest", "/api", "[5030, 5040]") + testUnauthorizedProvider,
				ExpectError: regexp.MustCompile("Client Error"),
			},
			// Create and Read testing
			{
				Config: testAccIndexerNewznabResourceConfig("newzabResourceTest", "/api", "[5030, 5040]"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("sonarr_indexer_newznab.test", "enable_automatic_search", "false"),
					resource.TestCheckResourceAttr("sonarr_indexer_newznab.test", "base_url", "https://lolo.sickbeard.com"),
					resource.TestCheckResourceAttrSet("sonarr_indexer_newznab.test", "id"),
				),
			},
			// Unauthorized Read
			{
				Config:      testAccIndexerNewznabResourceConfig("newzabResourceTest", "/api", "[5030, 5040]") + testUnauthorizedProvider,
				ExpectError: regexp.MustCompile("Client Error"),
			},
			// Update and Read testing
			{
				Config: testAccIndexerNewznabResourceConfig("newzabResourceTest", "/apis", "[5030, 5040]"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("sonarr_indexer_newznab.test", "api_path", "/apis"),
				),
			},
			// Update to an empty set and Read testing
			{
				Config: testAccIndexerNewznabResourceConfig("newzabResourceTest", "/apis", "[]"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("sonarr_indexer_newznab.test", "categories.#", "0"),
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

func testAccIndexerNewznabResourceConfig(name, apiPath, categories string) string {
	return fmt.Sprintf(`
	resource "sonarr_indexer_newznab" "test" {
		enable_automatic_search = false
		name = "%s"
		base_url = "https://lolo.sickbeard.com"
		api_path = "%s"
		categories = %s
	}`, name, apiPath, categories)
}

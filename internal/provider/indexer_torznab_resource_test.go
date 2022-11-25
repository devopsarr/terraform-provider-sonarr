package provider

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccIndexerTorznabResource(t *testing.T) {
	t.Parallel()

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read testing
			{
				Config: testAccIndexerTorznabResourceConfig("torznabResourceTest", 1),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("sonarr_indexer_torznab.test", "minimum_seeders", "1"),
					resource.TestCheckResourceAttrSet("sonarr_indexer_torznab.test", "id"),
				),
			},
			// Update and Read testing
			{
				Config: testAccIndexerTorznabResourceConfig("torznabResourceTest", 2),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("sonarr_indexer_torznab.test", "minimum_seeders", "2"),
				),
			},
			// ImportState testing
			{
				ResourceName:      "sonarr_indexer_torznab.test",
				ImportState:       true,
				ImportStateVerify: true,
			},
			// Delete testing automatically occurs in TestCase
		},
	})
}

func testAccIndexerTorznabResourceConfig(name string, seeders int) string {
	return fmt.Sprintf(`
	resource "sonarr_indexer_torznab" "test" {
		enable_automatic_search = false
		name = "%s"
		base_url = "https://feed.animetosho.org"
		api_path = "/nabapi"
		minimum_seeders = %d
		anime_categories = [6070]
	}`, name, seeders)
}

package provider

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccIndexerNyaaResource(t *testing.T) {
	t.Parallel()

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read testing
			{
				Config: testAccIndexerNyaaResourceConfig("nyaaResourceTest", "https://nyaa.org"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("sonarr_indexer_nyaa.test", "base_url", "https://nyaa.org"),
					resource.TestCheckResourceAttrSet("sonarr_indexer_nyaa.test", "id"),
				),
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

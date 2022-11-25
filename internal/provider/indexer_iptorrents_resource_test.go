package provider

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccIndexerIptorrentsResource(t *testing.T) {
	t.Parallel()

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read testing
			{
				Config: testAccIndexerIptorrentsResourceConfig("iptorrentsResourceTest", "https://iptorrents.org"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("sonarr_indexer_iptorrents.test", "base_url", "https://iptorrents.org"),
					resource.TestCheckResourceAttrSet("sonarr_indexer_iptorrents.test", "id"),
				),
			},
			// Update and Read testing
			{
				Config: testAccIndexerIptorrentsResourceConfig("iptorrentsResourceTest", "https://iptorrents.net"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("sonarr_indexer_iptorrents.test", "base_url", "https://iptorrents.net"),
				),
			},
			// ImportState testing
			{
				ResourceName:      "sonarr_indexer_iptorrents.test",
				ImportState:       true,
				ImportStateVerify: true,
			},
			// Delete testing automatically occurs in TestCase
		},
	})
}

func testAccIndexerIptorrentsResourceConfig(name, url string) string {
	return fmt.Sprintf(`
	resource "sonarr_indexer_iptorrents" "test" {
		enable_automatic_search = false
		name = "%s"
		base_url = "%s"
		minimum_seeders = 1
	}`, name, url)
}

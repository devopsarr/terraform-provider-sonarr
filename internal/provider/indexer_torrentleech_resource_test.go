package provider

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccIndexerTorrentleechResource(t *testing.T) {
	t.Parallel()

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Unauthorized Create
			{
				Config:      testAccIndexerTorrentleechResourceConfig("torrentleechResourceTest", 1) + testUnauthorizedProvider,
				ExpectError: regexp.MustCompile("Client Error"),
			},
			// Create and Read testing
			{
				Config: testAccIndexerTorrentleechResourceConfig("torrentleechResourceTest", 1),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("sonarr_indexer_torrentleech.test", "minimum_seeders", "1"),
					resource.TestCheckResourceAttrSet("sonarr_indexer_torrentleech.test", "id"),
				),
			},
			// Unauthorized Read
			{
				Config:      testAccIndexerTorrentleechResourceConfig("torrentleechResourceTest", 1) + testUnauthorizedProvider,
				ExpectError: regexp.MustCompile("Client Error"),
			},
			// Update and Read testing
			{
				Config: testAccIndexerTorrentleechResourceConfig("torrentleechResourceTest", 2),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("sonarr_indexer_torrentleech.test", "minimum_seeders", "2"),
				),
			},
			// ImportState testing
			{
				ResourceName:            "sonarr_indexer_torrentleech.test",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"api_key"},
			},
			// Delete testing automatically occurs in TestCase
		},
	})
}

func testAccIndexerTorrentleechResourceConfig(name string, seeders int) string {
	return fmt.Sprintf(`
	resource "sonarr_indexer_torrentleech" "test" {
		enable_automatic_search = false
		name = "%s"
		base_url = "http://rss.torrentleech.org"
		api_key = "Key"
		minimum_seeders = %d
	}`, name, seeders)
}

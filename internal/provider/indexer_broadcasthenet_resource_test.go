package provider

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccIndexerBroadcastheNetResource(t *testing.T) {
	t.Parallel()

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read testing
			{
				Config: testAccIndexerBroadcastheNetResourceConfig("broadcasthenetResourceTest", 1),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("sonarr_indexer_broadcasthenet.test", "seed_time", "1"),
					resource.TestCheckResourceAttrSet("sonarr_indexer_broadcasthenet.test", "id"),
				),
			},
			// Update and Read testing
			{
				Config: testAccIndexerBroadcastheNetResourceConfig("broadcasthenetResourceTest", 2),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("sonarr_indexer_broadcasthenet.test", "seed_time", "2"),
				),
			},
			// ImportState testing
			{
				ResourceName:      "sonarr_indexer_broadcasthenet.test",
				ImportState:       true,
				ImportStateVerify: true,
			},
			// Delete testing automatically occurs in TestCase
		},
	})
}

func testAccIndexerBroadcastheNetResourceConfig(name string, sTime int) string {
	return fmt.Sprintf(`
	resource "sonarr_indexer_broadcasthenet" "test" {
		enable_automatic_search = false
		name = "%s"
		base_url = "https://api.broadcasthe.net/"
		api_key = "API_key"
		minimum_seeders = 1
		season_pack_seed_time = 1
		seed_time = %d
		seed_ratio = 0.5
	}`, name, sTime)
}

package provider

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccIndexerConfigResource(t *testing.T) {
	t.Parallel()

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read testing
			{
				Config: testAccIndexerConfigResourceConfig(20),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("sonarr_indexer_config.test", "rss_sync_interval", "20"),
					resource.TestCheckResourceAttrSet("sonarr_indexer_config.test", "id"),
				),
			},
			// Update and Read testing
			{
				Config: testAccIndexerConfigResourceConfig(30),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("sonarr_indexer_config.test", "rss_sync_interval", "30"),
				),
			},
			// ImportState testing
			{
				ResourceName:      "sonarr_indexer_config.test",
				ImportState:       true,
				ImportStateVerify: true,
			},
			// Delete testing automatically occurs in TestCase
		},
	})
}

func testAccIndexerConfigResourceConfig(rss int) string {
	return fmt.Sprintf(`
	resource "sonarr_indexer_config" "test" {
		maximum_size = 0
		minimum_age = 0
		retention = 0
		rss_sync_interval = %d
	}`, rss)
}

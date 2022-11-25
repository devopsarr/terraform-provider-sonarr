package provider

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccIndexerOmgwtfnzbsResource(t *testing.T) {
	t.Parallel()

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read testing
			{
				Config: testAccIndexerOmgwtfnzbsResourceConfig("newzabResourceTest", 30),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("sonarr_indexer_omgwtfnzbs.test", "delay", "30"),
					resource.TestCheckResourceAttrSet("sonarr_indexer_omgwtfnzbs.test", "id"),
				),
			},
			// Update and Read testing
			{
				Config: testAccIndexerOmgwtfnzbsResourceConfig("newzabResourceTest", 60),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("sonarr_indexer_omgwtfnzbs.test", "delay", "60"),
				),
			},
			// ImportState testing
			{
				ResourceName:      "sonarr_indexer_omgwtfnzbs.test",
				ImportState:       true,
				ImportStateVerify: true,
			},
			// Delete testing automatically occurs in TestCase
		},
	})
}

func testAccIndexerOmgwtfnzbsResourceConfig(name string, delay int) string {
	return fmt.Sprintf(`
	resource "sonarr_indexer_omgwtfnzbs" "test" {
		enable_automatic_search = false
		name = "%s"
		username = "Username"
		api_key = "API_Key"
		delay = %d
	}`, name, delay)
}

package provider

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccIndexerResource(t *testing.T) {
	t.Parallel()

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read testing
			{
				Config: testAccIndexerResourceConfig("resourceTest", "false"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("sonarr_indexer.test", "enable_automatic_search", "false"),
					resource.TestCheckResourceAttr("sonarr_indexer.test", "base_url", "https://lolo.sickbeard.com"),
					resource.TestCheckResourceAttrSet("sonarr_indexer.test", "id"),
				),
			},
			// Update and Read testing
			{
				Config: testAccIndexerResourceConfig("resourceTest", "true"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("sonarr_indexer.test", "enable_automatic_search", "true"),
				),
			},
			// ImportState testing
			{
				ResourceName:      "sonarr_indexer.test",
				ImportState:       true,
				ImportStateVerify: true,
			},
			// Delete testing automatically occurs in TestCase
		},
	})
}

func testAccIndexerResourceConfig(name, aSearch string) string {
	return fmt.Sprintf(`
	resource "sonarr_indexer" "test" {
		enable_automatic_search = %s
		name = "%s"
		implementation = "Newznab"
		protocol = "usenet"
    	config_contract = "NewznabSettings"
		base_url = "https://lolo.sickbeard.com"
		api_path = "/api"
		categories = [5030, 5040]
		tags = []
	}`, aSearch, name)
}

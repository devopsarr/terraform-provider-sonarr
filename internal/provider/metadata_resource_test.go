package provider

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccMetadataResource(t *testing.T) {
	t.Parallel()

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read testing
			{
				Config: testAccMetadataResourceConfig("resourceTest", "true"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("sonarr_metadata.test", "episode_metadata", "true"),
					resource.TestCheckResourceAttrSet("sonarr_metadata.test", "id"),
				),
			},
			// Update and Read testing
			{
				Config: testAccMetadataResourceConfig("resourceTest", "false"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("sonarr_metadata.test", "episode_metadata", "false"),
				),
			},
			// ImportState testing
			{
				ResourceName:      "sonarr_metadata.test",
				ImportState:       true,
				ImportStateVerify: true,
			},
			// Delete testing automatically occurs in TestCase
		},
	})
}

func testAccMetadataResourceConfig(name, metadata string) string {
	return fmt.Sprintf(`
	resource "sonarr_metadata" "test" {
		enable = true
		name = "%s"
		implementation = "WdtvMetadata"
    	config_contract = "WdtvMetadataSettings"
		episode_metadata = %s
		series_images = false
		season_images = true
		episode_images = false
	}`, name, metadata)
}

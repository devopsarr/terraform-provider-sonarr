package provider

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccMetadataRoksboxResource(t *testing.T) {
	t.Parallel()

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read testing
			{
				Config: testAccMetadataRoksboxResourceConfig("roksboxResourceTest", "false"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("sonarr_metadata_roksbox.test", "episode_metadata", "false"),
					resource.TestCheckResourceAttrSet("sonarr_metadata_roksbox.test", "id"),
				),
			},
			// Update and Read testing
			{
				Config: testAccMetadataRoksboxResourceConfig("roksboxResourceTest", "true"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("sonarr_metadata_roksbox.test", "episode_metadata", "true"),
				),
			},
			// ImportState testing
			{
				ResourceName:      "sonarr_metadata_roksbox.test",
				ImportState:       true,
				ImportStateVerify: true,
			},
			// Delete testing automatically occurs in TestCase
		},
	})
}

func testAccMetadataRoksboxResourceConfig(name, metadata string) string {
	return fmt.Sprintf(`
	resource "sonarr_metadata_roksbox" "test" {
		enable = false
		name = "%s"
		episode_metadata = %s
		series_images = false
		season_images = true
		episode_images = false
	}`, name, metadata)
}

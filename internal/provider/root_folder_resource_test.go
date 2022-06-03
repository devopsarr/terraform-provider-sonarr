package provider

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccRootFolderResource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read testing
			{
				Config: testAccRootFolderResourceConfig("/tmp"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("sonarr_root_folder.test", "path", "/tmp"),
					resource.TestCheckResourceAttrSet("sonarr_root_folder.test", "id"),
				),
			},
			// Update and Read testing
			{
				Config: testAccRootFolderResourceConfig("/defaults"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("sonarr_root_folder.test", "path", "/defaults"),
				),
			},
			// ImportState testing
			{
				ResourceName:      "sonarr_root_folder.test",
				ImportState:       true,
				ImportStateVerify: true,
			},
			// Delete testing automatically occurs in TestCase
		},
	})
}

func testAccRootFolderResourceConfig(path string) string {
	return fmt.Sprintf(`
		resource "sonarr_root_folder" "test" {
  			path = "%s"
		}
	`, path)
}

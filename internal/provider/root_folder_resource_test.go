package provider

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccRootFolderResource(t *testing.T) {
	t.Parallel()

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Unauthorized Create
			{
				Config:      testAccRootFolderResourceConfig("/error") + testUnauthorizedProvider,
				ExpectError: regexp.MustCompile("Client Error"),
			},
			// Create and Read testing
			{
				Config: testAccRootFolderResourceConfig("/config/asp"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("sonarr_root_folder.test", "path", "/config/asp"),
					resource.TestCheckResourceAttrSet("sonarr_root_folder.test", "id"),
				),
			},
			// Unauthorized Read
			{
				Config:      testAccRootFolderResourceConfig("/error") + testUnauthorizedProvider,
				ExpectError: regexp.MustCompile("Client Error"),
			},
			// Update and Read testing
			{
				Config: testAccRootFolderResourceConfig("/config/logs"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("sonarr_root_folder.test", "path", "/config/logs"),
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

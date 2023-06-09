package provider

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccIndexerFilelistResource(t *testing.T) {
	t.Parallel()

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Unauthorized Create
			{
				Config:      testAccIndexerFilelistResourceConfig("filelistResourceTest", "user") + testUnauthorizedProvider,
				ExpectError: regexp.MustCompile("Client Error"),
			},
			// Create and Read testing
			{
				Config: testAccIndexerFilelistResourceConfig("filelistResourceTest", "user"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("sonarr_indexer_filelist.test", "username", "user"),
					resource.TestCheckResourceAttrSet("sonarr_indexer_filelist.test", "id"),
				),
			},
			// Unauthorized Read
			{
				Config:      testAccIndexerFilelistResourceConfig("filelistResourceTest", "user") + testUnauthorizedProvider,
				ExpectError: regexp.MustCompile("Client Error"),
			},
			// Update and Read testing
			{
				Config: testAccIndexerFilelistResourceConfig("filelistResourceTest", "Username"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("sonarr_indexer_filelist.test", "username", "Username"),
				),
			},
			// ImportState testing
			{
				ResourceName:            "sonarr_indexer_filelist.test",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"passkey"},
			},
			// Delete testing automatically occurs in TestCase
		},
	})
}

func testAccIndexerFilelistResourceConfig(name, username string) string {
	return fmt.Sprintf(`
	resource "sonarr_indexer_filelist" "test" {
		enable_automatic_search = false
		name = "%s"
		base_url = "https://filelist.io"
		username = "%s"
		passkey = "Pass"
		categories = [21,23,27]
		minimum_seeders = 1
	}`, name, username)
}

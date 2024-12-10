package provider

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccIndexerResource(t *testing.T) {
	t.Parallel()

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Unauthorized Create
			{
				Config:      testAccIndexerResourceConfig("resourceTest", "/api") + testUnauthorizedProvider,
				ExpectError: regexp.MustCompile("Client Error"),
			},
			// Create and Read testing
			{
				Config: testAccIndexerResourceConfig("resourceTest", "/api"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("sonarr_indexer.test", "enable_automatic_search", "false"),
					resource.TestCheckResourceAttr("sonarr_indexer.test", "base_url", "https://lolo.sickbeard.com"),
					resource.TestCheckResourceAttrSet("sonarr_indexer.test", "id"),
				),
			},
			// Unauthorized Read
			{
				Config:      testAccIndexerResourceConfig("resourceTest", "/api") + testUnauthorizedProvider,
				ExpectError: regexp.MustCompile("Client Error"),
			},
			// Update and Read testing
			{
				Config: testAccIndexerResourceConfig("resourceTest", "/apis"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("sonarr_indexer.test", "api_path", "/apis"),
				),
			},
			// ImportState testing
			{
				ResourceName:      "sonarr_indexer.test",
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				ResourceName:            "sonarr_indexer.test_sensitive",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"passkey"},
			},
			// Delete testing automatically occurs in TestCase
		},
	})
}

func testAccIndexerResourceConfig(name, apiPath string) string {
	return fmt.Sprintf(`
	resource "sonarr_indexer" "test" {
		enable_automatic_search = false
		name = "%s"
		implementation = "Newznab"
		protocol = "usenet"
    	config_contract = "NewznabSettings"
		base_url = "https://lolo.sickbeard.com"
		api_path = "%s"
		categories = [5030, 5040]
		tags = []
	}

	resource "sonarr_indexer" "test_sensitive" {
		enable_automatic_search = false
		name = "%sWithSensitive"
		base_url = "https://filelist.io"
		username = "test"
		passkey = "Pass"
		categories = [21,23,27]
		minimum_seeders = 1
		implementation = "FileList"
		protocol = "torrent"
    	config_contract = "FileListSettings"
	}
	`, name, apiPath, name)
}

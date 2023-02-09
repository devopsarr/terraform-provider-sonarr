package provider

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccIndexerHdbitsResource(t *testing.T) {
	t.Parallel()

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Unauthorized Create
			{
				Config:      testAccIndexerHdbitsResourceConfig("hdbitsResourceTest", "user") + testUnauthorizedProvider,
				ExpectError: regexp.MustCompile("Client Error"),
			},
			// Create and Read testing
			{
				Config: testAccIndexerHdbitsResourceConfig("hdbitsResourceTest", "user"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("sonarr_indexer_hdbits.test", "username", "user"),
					resource.TestCheckResourceAttrSet("sonarr_indexer_hdbits.test", "id"),
				),
			},
			// Unauthorized Read
			{
				Config:      testAccIndexerHdbitsResourceConfig("hdbitsResourceTest", "user") + testUnauthorizedProvider,
				ExpectError: regexp.MustCompile("Client Error"),
			},
			// Update and Read testing
			{
				Config: testAccIndexerHdbitsResourceConfig("hdbitsResourceTest", "Username"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("sonarr_indexer_hdbits.test", "username", "Username"),
				),
			},
			// ImportState testing
			{
				ResourceName:            "sonarr_indexer_hdbits.test",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"api_key"},
			},
			// Delete testing automatically occurs in TestCase
		},
	})
}

func testAccIndexerHdbitsResourceConfig(name, username string) string {
	return fmt.Sprintf(`
	resource "sonarr_indexer_hdbits" "test" {
		enable_automatic_search = false
		name = "%s"
		base_url = "https://hdbits.org"
		username = "%s"
		api_key = "Key"
		minimum_seeders = 1
	}`, name, username)
}

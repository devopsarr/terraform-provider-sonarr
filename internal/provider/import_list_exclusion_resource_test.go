package provider

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccImportListExclusionResource(t *testing.T) {
	t.Parallel()

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Unauthorized Create
			{
				Config:      testAccImportListExclusionResourceConfig("test", 123) + testUnauthorizedProvider,
				ExpectError: regexp.MustCompile("Client Error"),
			},
			// Create and Read testing
			{
				Config: testAccImportListExclusionResourceConfig("test", 123),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("sonarr_import_list_exclusion.test", "tvdb_id", "123"),
					resource.TestCheckResourceAttrSet("sonarr_import_list_exclusion.test", "id"),
				),
			},
			// Unauthorized Read
			{
				Config:      testAccImportListExclusionResourceConfig("test", 123) + testUnauthorizedProvider,
				ExpectError: regexp.MustCompile("Client Error"),
			},
			// Update and Read testing
			{
				Config: testAccImportListExclusionResourceConfig("test", 1234),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("sonarr_import_list_exclusion.test", "tvdb_id", "1234"),
				),
			},
			// ImportState testing
			{
				ResourceName:      "sonarr_import_list_exclusion.test",
				ImportState:       true,
				ImportStateVerify: true,
			},
			// Delete testing automatically occurs in TestCase
		},
	})
}

func testAccImportListExclusionResourceConfig(name string, tvID int) string {
	return fmt.Sprintf(`
		resource "sonarr_import_list_exclusion" "%s" {
  			title = "Friends"
			tvdb_id = %d
		}
	`, name, tvID)
}

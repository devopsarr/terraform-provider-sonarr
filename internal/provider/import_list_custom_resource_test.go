package provider

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccImportListCustomResource(t *testing.T) {
	t.Parallel()

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Unauthorized Create
			{
				Config:      testAccImportListCustomResourceConfig("resourceCustomTest", "false") + testUnauthorizedProvider,
				ExpectError: regexp.MustCompile("Client Error"),
			},
			// Create and Read testing
			{
				PreConfig: rootFolderDSInit,
				Config:    testAccImportListCustomResourceConfig("resourceCustomTest", "false"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("sonarr_import_list_custom.test", "season_folder", "false"),
					resource.TestCheckResourceAttrSet("sonarr_import_list_custom.test", "id"),
				),
			},
			// Unauthorized Read
			{
				Config:      testAccImportListCustomResourceConfig("resourceCustomTest", "false") + testUnauthorizedProvider,
				ExpectError: regexp.MustCompile("Client Error"),
			},
			// Update and Read testing
			{
				Config: testAccImportListCustomResourceConfig("resourceCustomTest", "true"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("sonarr_import_list_custom.test", "season_folder", "true"),
				),
			},
			// ImportState testing
			{
				ResourceName:      "sonarr_import_list_custom.test",
				ImportState:       true,
				ImportStateVerify: true,
			},
			// Delete testing automatically occurs in TestCase
		},
	})
}

func testAccImportListCustomResourceConfig(name, folder string) string {
	return fmt.Sprintf(`

	resource "sonarr_import_list_custom" "test" {
		enable_automatic_add = false
		season_folder = %s
		should_monitor = "all"
		series_type = "standard"
		root_folder_path = "/config"
		quality_profile_id = 1
		name = "%s"
		base_url = "localhost"
		tags = []
	}`, folder, name)
}

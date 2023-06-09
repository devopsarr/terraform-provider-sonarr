package provider

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccImportListTraktListResource(t *testing.T) {
	t.Parallel()

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Unauthorized Create
			{
				Config:      testAccImportListTraktListResourceConfig("resourceTraktListTest", "false") + testUnauthorizedProvider,
				ExpectError: regexp.MustCompile("Client Error"),
			},
			// Create and Read testing
			{
				PreConfig: rootFolderDSInit,
				Config:    testAccImportListTraktListResourceConfig("resourceTraktListTest", "false"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("sonarr_import_list_trakt_list.test", "season_folder", "false"),
					resource.TestCheckResourceAttrSet("sonarr_import_list_trakt_list.test", "id"),
				),
			},
			// Unauthorized Read
			{
				Config:      testAccImportListTraktListResourceConfig("resourceTraktListTest", "false") + testUnauthorizedProvider,
				ExpectError: regexp.MustCompile("Client Error"),
			},
			// Update and Read testing
			{
				Config: testAccImportListTraktListResourceConfig("resourceTraktListTest", "true"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("sonarr_import_list_trakt_list.test", "season_folder", "true"),
				),
			},
			// ImportState testing
			{
				ResourceName:      "sonarr_import_list_trakt_list.test",
				ImportState:       true,
				ImportStateVerify: true,
			},
			// Delete testing automatically occurs in TestCase
		},
	})
}

func testAccImportListTraktListResourceConfig(name, folder string) string {
	return fmt.Sprintf(`
	resource "sonarr_import_list_trakt_list" "test" {
		enable_automatic_add = false
		season_folder = %s
		should_monitor = "all"
		series_type = "standard"
		root_folder_path = "/config"
		quality_profile_id = 1
		name = "%s"
		access_token = "Token"
		username = "User"
		listname = "test"
		limit = 100
		tags = []
	}`, folder, name)
}

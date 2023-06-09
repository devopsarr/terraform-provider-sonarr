package provider

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccImportListTraktUserResource(t *testing.T) {
	t.Parallel()

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Unauthorized Create
			{
				Config:      testAccImportListTraktUserResourceConfig("resourceTraktUserTest", "false") + testUnauthorizedProvider,
				ExpectError: regexp.MustCompile("Client Error"),
			},
			// Create and Read testing
			{
				PreConfig: rootFolderDSInit,
				Config:    testAccImportListTraktUserResourceConfig("resourceTraktUserTest", "false"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("sonarr_import_list_trakt_user.test", "season_folder", "false"),
					resource.TestCheckResourceAttrSet("sonarr_import_list_trakt_user.test", "id"),
				),
			},
			// Unauthorized Read
			{
				Config:      testAccImportListTraktUserResourceConfig("resourceTraktUserTest", "false") + testUnauthorizedProvider,
				ExpectError: regexp.MustCompile("Client Error"),
			},
			// Update and Read testing
			{
				Config: testAccImportListTraktUserResourceConfig("resourceTraktUserTest", "true"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("sonarr_import_list_trakt_user.test", "season_folder", "true"),
				),
			},
			// ImportState testing
			{
				ResourceName:      "sonarr_import_list_trakt_user.test",
				ImportState:       true,
				ImportStateVerify: true,
			},
			// Delete testing automatically occurs in TestCase
		},
	})
}

func testAccImportListTraktUserResourceConfig(name, folder string) string {
	return fmt.Sprintf(`
	resource "sonarr_import_list_trakt_user" "test" {
		enable_automatic_add = false
		season_folder = %s
		should_monitor = "all"
		series_type = "standard"
		root_folder_path = "/config"
		quality_profile_id = 1
		name = "%s"
		access_token = "Token"
		username = "User"
		trakt_list_type = 0
		limit = 100
		tags = []
	}`, folder, name)
}

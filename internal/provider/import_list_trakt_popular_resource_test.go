package provider

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccImportListTraktPopularResource(t *testing.T) {
	t.Parallel()

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read testing
			{
				PreConfig: rootFolderDSInit,
				Config:    testAccImportListTraktPopularResourceConfig("resourceTraktPopularTest", "false"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("sonarr_import_list_trakt_popular.test", "season_folder", "false"),
					resource.TestCheckResourceAttrSet("sonarr_import_list_trakt_popular.test", "id"),
				),
			},
			// Update and Read testing
			{
				Config: testAccImportListTraktPopularResourceConfig("resourceTraktPopularTest", "true"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("sonarr_import_list_trakt_popular.test", "season_folder", "true"),
				),
			},
			// ImportState testing
			{
				ResourceName:      "sonarr_import_list_trakt_popular.test",
				ImportState:       true,
				ImportStateVerify: true,
			},
			// Delete testing automatically occurs in TestCase
		},
	})
}

func testAccImportListTraktPopularResourceConfig(name, folder string) string {
	return fmt.Sprintf(`
	data "sonarr_root_folder" "test" {
		path = "/defaults"
  	}

	resource "sonarr_import_list_trakt_popular" "test" {
		enable_automatic_add = false
		season_folder = %s
		should_monitor = "all"
		series_type = "standard"
		root_folder_path = data.sonarr_root_folder.test.path
		quality_profile_id = 1
		name = "%s"
		access_token = "Token"
		trakt_list_type = 0
		limit = 100
		tags = []
	}`, folder, name)
}

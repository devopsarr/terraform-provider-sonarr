package provider

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccImportListSonarrResource(t *testing.T) {
	t.Parallel()

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read testing
			{
				PreConfig: rootFolderDSInit,
				Config:    testAccImportListSonarrResourceConfig("resourceSonarrTest", "false"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("sonarr_import_list_sonarr.test", "enable_automatic_add", "false"),
					resource.TestCheckResourceAttrSet("sonarr_import_list_sonarr.test", "id"),
				),
			},
			// Update and Read testing
			{
				Config: testAccImportListSonarrResourceConfig("resourceSonarrTest", "true"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("sonarr_import_list_sonarr.test", "enable_automatic_add", "true"),
				),
			},
			// ImportState testing
			{
				ResourceName:      "sonarr_import_list_sonarr.test",
				ImportState:       true,
				ImportStateVerify: true,
			},
			// Delete testing automatically occurs in TestCase
		},
	})
}

func testAccImportListSonarrResourceConfig(name, add string) string {
	return fmt.Sprintf(`
	resource "sonarr_import_list_sonarr" "test" {
		enable_automatic_add = %s
		season_folder = true
		should_monitor = "all"
		series_type = "standard"
		root_folder_path = "/config"
		quality_profile_id = 1
		name = "%s"
		base_url = "http://127.0.0.1:8989"
		api_key = "b01df9fca2e64e459d64a09888ce7451"
		tags = []
	}`, add, name)
}

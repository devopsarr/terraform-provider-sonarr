package provider

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccImportListResource(t *testing.T) {
	t.Parallel()

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read testing
			{
				Config: testAccImportListResourceConfig("/config/.config", "importListResourceTest", "false"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("sonarr_import_list.test", "enable_automatic_add", "false"),
					resource.TestCheckResourceAttrSet("sonarr_import_list.test", "id"),
				),
			},
			// Update and Read testing
			{
				Config: testAccImportListResourceConfig("/config/.config", "importListResourceTest", "true"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("sonarr_import_list.test", "enable_automatic_add", "true"),
				),
			},
			// ImportState testing
			{
				ResourceName:      "sonarr_import_list.test",
				ImportState:       true,
				ImportStateVerify: true,
			},
			// Delete testing automatically occurs in TestCase
		},
	})
}

func testAccImportListResourceConfig(path, name, enable string) string {
	return fmt.Sprintf(`
	resource "sonarr_root_folder" "test" {
		path = "%s"
  	}

	  resource "sonarr_quality_profile" "test" {
		name            = "%s"
		upgrade_allowed = true
		cutoff          = 1100

		quality_groups = [
			{
				id   = 1100
				name = "4k"
				qualities = [
					{
						id         = 18
						name       = "WEBDL-2160p"
						source     = "web"
						resolution = 2160
					},
					{
						id         = 19
						name       = "Bluray-2160p"
						source     = "bluray"
						resolution = 2160
					}
				]
			}
		]
	}

	resource "sonarr_language_profile" "test" {
		upgrade_allowed = true
		name = "%s"
		cutoff_language = "English"
		languages = [ "English" ]
	}

	resource "sonarr_import_list" "test" {
		enable_automatic_add = %s
		season_folder = true
		should_monitor = "all"
		series_type = "standard"
		root_folder_path = sonarr_root_folder.test.path
		quality_profile_id = sonarr_quality_profile.test.id
		language_profile_id = sonarr_language_profile.test.id
		name = "%s"
		implementation = "SonarrImport"
    	config_contract = "SonarrSettings"
		base_url = "http://127.0.0.1:8989"
		api_key = "b01df9fca2e64e459d64a09888ce7451"
		tags = []
	}`, path, name, name, enable, name)
}

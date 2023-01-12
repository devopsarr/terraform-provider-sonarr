package provider

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccImportListDataSource(t *testing.T) {
	t.Parallel()

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Read testing
			{
				PreConfig: rootFolderDSInit,
				Config:    testAccImportListDataSourceConfig,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.sonarr_import_list.test", "id"),
					resource.TestCheckResourceAttr("data.sonarr_import_list.test", "should_monitor", "all")),
			},
		},
	})
}

const testAccImportListDataSourceConfig = `
resource "sonarr_import_list" "test" {
	enable_automatic_add = false
	season_folder = true
	should_monitor = "all"
	series_type = "standard"
	root_folder_path = "/config"
	quality_profile_id = 1
	name = "importListDataTest"
	implementation = "SonarrImport"
	config_contract = "SonarrSettings"
	base_url = "http://127.0.0.1:8989"
	api_key = "b01df9fca2e64e459d64a09888ce7451"
	tags = []
}

data "sonarr_import_list" "test" {
	name = sonarr_import_list.test.name
}
`

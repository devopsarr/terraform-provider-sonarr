package provider

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccImportListExclusionDataSource(t *testing.T) {
	t.Parallel()

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccImportListExclusionDataSourceConfig,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.sonarr_import_list_exclusion.test", "id"),
					resource.TestCheckResourceAttr("data.sonarr_import_list_exclusion.test", "title", "testDS"),
				),
			},
		},
	})
}

const testAccImportListExclusionDataSourceConfig = `
resource "sonarr_import_list_exclusion" "test" {
	title = "testDS"
	tvdb_id = 987
}

data "sonarr_import_list_exclusion" "test" {
	tvdb_id = sonarr_import_list_exclusion.test.tvdb_id
}
`

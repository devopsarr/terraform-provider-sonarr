package provider

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccImportListsDataSource(t *testing.T) {
	t.Parallel()

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create a delay profile to have a value to check
			{
				Config: testAccImportListResourceConfig("/config/.config/.mono", "importListsDataTest", "false"),
			},
			// Read testing
			{
				Config: testAccImportListsDataSourceConfig,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckTypeSetElemNestedAttrs("data.sonarr_import_lists.test", "import_lists.*", map[string]string{"base_url": "http://127.0.0.1:8989"}),
				),
			},
		},
	})
}

const testAccImportListsDataSourceConfig = `
data "sonarr_import_lists" "test" {
}
`

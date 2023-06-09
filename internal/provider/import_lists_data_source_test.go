package provider

import (
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccImportListsDataSource(t *testing.T) {
	t.Parallel()

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Unauthorized
			{
				Config:      testAccImportListsDataSourceConfig + testUnauthorizedProvider,
				ExpectError: regexp.MustCompile("Client Error"),
			},
			// Create a resource to have a value to check
			{
				PreConfig: rootFolderDSInit,
				Config:    testAccImportListResourceConfig("importListsDataTest", "false"),
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

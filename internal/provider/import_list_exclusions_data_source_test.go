package provider

import (
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccImportListExclusionsDataSource(t *testing.T) {
	t.Parallel()

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Unauthorized
			{
				Config:      testAccImportListExclusionsDataSourceConfig + testUnauthorizedProvider,
				ExpectError: regexp.MustCompile("Client Error"),
			},
			// Create a resource to have a value to check
			{
				Config: testAccImportListExclusionResourceConfig("testList", 321),
			},
			// Read testing
			{
				Config: testAccImportListExclusionsDataSourceConfig,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckTypeSetElemNestedAttrs("data.sonarr_import_list_exclusions.test", "import_list_exclusions.*", map[string]string{"tvdb_id": "321"}),
				),
			},
		},
	})
}

const testAccImportListExclusionsDataSourceConfig = `
data "sonarr_import_list_exclusions" "test" {
}
`

package provider

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccRootFoldersDataSource(t *testing.T) {
	t.Parallel()

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create a root folder to have a value to check
			{
				Config: testAccRootFolderResourceConfig("/tmp"),
			},
			// Read testing
			{
				Config: testAccRootFoldersDataSourceConfig,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckTypeSetElemNestedAttrs("data.sonarr_root_folders.test", "root_folders.*", map[string]string{"path": "/tmp"}),
				),
			},
		},
	})
}

const testAccRootFoldersDataSourceConfig = `
data "sonarr_root_folders" "test" {
}
`

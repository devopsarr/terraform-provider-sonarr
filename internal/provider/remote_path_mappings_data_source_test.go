package provider

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccRemotePathMappingsDataSource(t *testing.T) {
	t.Parallel()

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create a tag to have a value to check
			{
				Config: testAccRemotePathMappingResourceConfig("remotemapDataSourceTest", "/test3/"),
			},
			// Read testing
			{
				Config: testAccRemotePathMappingsDataSourceConfig,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckTypeSetElemNestedAttrs("data.sonarr_remote_path_mappings.test", "remote_path_mappings.*", map[string]string{"remote_path": "/test3/"}),
				),
			},
		},
	})
}

const testAccRemotePathMappingsDataSourceConfig = `
data "sonarr_remote_path_mappings" "test" {
}
`

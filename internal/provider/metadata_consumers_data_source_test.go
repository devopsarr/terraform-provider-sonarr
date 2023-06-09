package provider

import (
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccMetadataConsumersDataSource(t *testing.T) {
	t.Parallel()

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Unauthorized
			{
				Config:      testAccMetadataConsumersDataSourceConfig + testUnauthorizedProvider,
				ExpectError: regexp.MustCompile("Client Error"),
			},
			// Create a resource to have a value to check
			{
				Config: testAccMetadataResourceConfig("datasourceTest", "false"),
			},
			// Read testing
			{
				Config: testAccMetadataConsumersDataSourceConfig,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckTypeSetElemNestedAttrs("data.sonarr_metadata_consumers.test", "metadata_consumers.*", map[string]string{"episode_metadata": "false"}),
				),
			},
		},
	})
}

const testAccMetadataConsumersDataSourceConfig = `
data "sonarr_metadata_consumers" "test" {
}
`

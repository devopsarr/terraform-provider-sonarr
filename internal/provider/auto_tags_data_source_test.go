package provider

import (
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccAutoTagsDataSource(t *testing.T) {
	t.Parallel()

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Unauthorized
			{
				Config:      testAccAutoTagsDataSourceConfig + testUnauthorizedProvider,
				ExpectError: regexp.MustCompile("Client Error"),
			},
			// Create a resource to have a value to check
			{
				Config: testAccTagResourceConfig("test", "dataautotag") + testAccAutoTagResourceConfig("datasourceTest", "true"),
			},
			// Read testing
			{
				Config: testAccAutoTagsDataSourceConfig,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckTypeSetElemNestedAttrs("data.sonarr_auto_tags.test", "auto_tags.*", map[string]string{"remove_tags_automatically": "true"}),
				),
			},
		},
	})
}

const testAccAutoTagsDataSourceConfig = `
data "sonarr_auto_tags" "test" {
}
`

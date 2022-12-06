package provider

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccQualityDefinitionDataSource(t *testing.T) {
	t.Parallel()

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Read testing
			{
				Config: testAccQualityDefinitionDataSourceConfig,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.sonarr_quality_definition.test", "title"),
					resource.TestCheckResourceAttr("data.sonarr_quality_definition.test", "id", "21")),
			},
		},
	})
}

const testAccQualityDefinitionDataSourceConfig = `
data "sonarr_quality_definition" "test" {
	id = 21
}
`

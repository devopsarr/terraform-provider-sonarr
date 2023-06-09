package provider

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccCustomFormatConditionResolutionDataSource(t *testing.T) {
	t.Parallel()

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Read testing
			{
				Config: testAccCustomFormatConditionResolutionDataSourceConfig,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.sonarr_custom_format_condition_resolution.test", "id"),
					resource.TestCheckResourceAttr("data.sonarr_custom_format_condition_resolution.test", "name", "4K"),
					resource.TestCheckResourceAttr("sonarr_custom_format.test", "specifications.0.value", "2160")),
			},
		},
	})
}

const testAccCustomFormatConditionResolutionDataSourceConfig = `
data  "sonarr_custom_format_condition_resolution" "test" {
	name = "4K"
	negate = false
	required = false
	value = "2160"
}

resource "sonarr_custom_format" "test" {
	include_custom_format_when_renaming = false
	name = "TestWithDSResolution"
	
	specifications = [data.sonarr_custom_format_condition_resolution.test]	
}`

package provider

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccCustomFormatConditionSourceDataSource(t *testing.T) {
	t.Parallel()

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Read testing
			{
				Config: testAccCustomFormatConditionSourceDataSourceConfig,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.sonarr_custom_format_condition_source.test", "id"),
					resource.TestCheckResourceAttr("data.sonarr_custom_format_condition_source.test", "name", "WEBDL"),
					resource.TestCheckResourceAttr("sonarr_custom_format.test", "specifications.0.value", "7")),
			},
		},
	})
}

const testAccCustomFormatConditionSourceDataSourceConfig = `
data  "sonarr_custom_format_condition_source" "test" {
	name = "WEBDL"
	negate = false
	required = false
	value = 7
}

resource "sonarr_custom_format" "test" {
	include_custom_format_when_renaming = false
	name = "TestWithDSSource"
	
	specifications = [data.sonarr_custom_format_condition_source.test]	
}`

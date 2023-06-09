package provider

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccCustomFormatConditionReleaseTitleDataSource(t *testing.T) {
	t.Parallel()

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Read testing
			{
				Config: testAccCustomFormatConditionReleaseTitleDataSourceConfig,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.sonarr_custom_format_condition_release_title.test", "id"),
					resource.TestCheckResourceAttr("data.sonarr_custom_format_condition_release_title.test", "name", "x265"),
					resource.TestCheckResourceAttr("sonarr_custom_format.test", "specifications.0.value", "(((x|h)\\.?265)|(HEVC))")),
			},
		},
	})
}

const testAccCustomFormatConditionReleaseTitleDataSourceConfig = `
data  "sonarr_custom_format_condition_release_title" "test" {
	name = "x265"
	negate = false
	required = false
	value = "(((x|h)\\.?265)|(HEVC))"
}

resource "sonarr_custom_format" "test" {
	include_custom_format_when_renaming = false
	name = "TestWithDSReleaseTitle"
	
	specifications = [data.sonarr_custom_format_condition_release_title.test]	
}`

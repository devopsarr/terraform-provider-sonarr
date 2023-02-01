package provider

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccCustomFormatConditionLanguageDataSource(t *testing.T) {
	t.Parallel()

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Read testing
			{
				Config: testAccCustomFormatConditionLanguageDataSourceConfig,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.sonarr_custom_format_condition_language.test", "id"),
					resource.TestCheckResourceAttr("data.sonarr_custom_format_condition_language.test", "name", "Arabic"),
					resource.TestCheckResourceAttr("sonarr_custom_format.test", "specifications.0.value", "31")),
			},
		},
	})
}

const testAccCustomFormatConditionLanguageDataSourceConfig = `
data  "sonarr_custom_format_condition_language" "test" {
	name = "Arabic"
	negate = false
	required = false
	value = "31"
}

resource "sonarr_custom_format" "test" {
	include_custom_format_when_renaming = false
	name = "TestWithDSLanguage"
	
	specifications = [data.sonarr_custom_format_condition_language.test]	
}`

package provider

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccAutoTagConditionSeriesTypeDataSource(t *testing.T) {
	t.Parallel()

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Read testing
			{
				Config: testAccAutoTagConditionSeriesTypeDataSourceConfig,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.sonarr_auto_tag_condition_series_type.test", "id"),
					resource.TestCheckResourceAttr("data.sonarr_auto_tag_condition_series_type.test", "name", "Test"),
					resource.TestCheckResourceAttr("sonarr_auto_tag.test", "specifications.0.value", "1")),
			},
		},
	})
}

const testAccAutoTagConditionSeriesTypeDataSourceConfig = `
resource "sonarr_tag" "test" {
	label = "atconditiontype"
}

data  "sonarr_auto_tag_condition_series_type" "test" {
	name = "Test"
	negate = false
	required = false
	value = "1"
}

resource "sonarr_auto_tag" "test" {
	remove_tags_automatically = false
	name = "TestWithDSSeriesType"

	tags = [sonarr_tag.test.id]
	
	specifications = [data.sonarr_auto_tag_condition_series_type.test]	
}`

package provider

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccAutoTagConditionDataSource(t *testing.T) {
	t.Parallel()

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Read testing
			{
				Config: testAccAutoTagConditionDataSourceConfig,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.sonarr_auto_tag_condition.test", "id"),
					resource.TestCheckResourceAttr("data.sonarr_auto_tag_condition.test", "name", "type"),
					resource.TestCheckResourceAttr("sonarr_auto_tag.test", "specifications.0.implementation", "SeriesTypeSpecification")),
			},
		},
	})
}

const testAccAutoTagConditionDataSourceConfig = `
resource "sonarr_tag" "test" {
	label = "atcondition"
}

data  "sonarr_auto_tag_condition" "test" {
	name = "type"
	implementation = "SeriesTypeSpecification"
	value = "1"
	required = false
	negate = true
}

resource "sonarr_auto_tag" "test" {
	remove_tags_automatically = false
	name = "TestWithDS"
	tags = [sonarr_tag.test.id]
	
	specifications = [data.sonarr_auto_tag_condition.test]	
}`

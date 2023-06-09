package provider

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccAutoTagConditionGenresDataSource(t *testing.T) {
	t.Parallel()

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Read testing
			{
				Config: testAccAutoTagConditionGenresDataSourceConfig,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.sonarr_auto_tag_condition_genres.test", "id"),
					resource.TestCheckResourceAttr("data.sonarr_auto_tag_condition_genres.test", "name", "Test"),
					resource.TestCheckResourceAttr("sonarr_auto_tag.test", "specifications.0.value", "horror comedy")),
			},
		},
	})
}

const testAccAutoTagConditionGenresDataSourceConfig = `
resource "sonarr_tag" "test" {
	label = "atconditiongenre"
}

data  "sonarr_auto_tag_condition_genres" "test" {
	name = "Test"
	negate = false
	required = false
	value = "horror comedy"
}

resource "sonarr_auto_tag" "test" {
	remove_tags_automatically = false
	name = "TestWithDSGenres"

	tags = [sonarr_tag.test.id]
	
	specifications = [data.sonarr_auto_tag_condition_genres.test]	
}`

package provider

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccAutoTagDataSource(t *testing.T) {
	t.Parallel()

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Unauthorized
			{
				Config:      testAccAutoTagDataSourceConfig("\"Error\"") + testUnauthorizedProvider,
				ExpectError: regexp.MustCompile("Client Error"),
			},
			// Not found testing
			{
				Config:      testAccAutoTagDataSourceConfig("\"Error\""),
				ExpectError: regexp.MustCompile("Unable to find auto_tag"),
			},
			// Read testing
			{
				Config: testAccTagResourceConfig("test", "singledataautotag") + testAccAutoTagResourceConfig("dataTest", "false") + testAccAutoTagDataSourceConfig("sonarr_auto_tag.test.name"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.sonarr_auto_tag.test", "id"),
					resource.TestCheckResourceAttr("data.sonarr_auto_tag.test", "remove_tags_automatically", "false")),
			},
		},
	})
}

func testAccAutoTagDataSourceConfig(name string) string {
	return fmt.Sprintf(`
	data "sonarr_auto_tag" "test" {
		name = %s
	}
	`, name)
}

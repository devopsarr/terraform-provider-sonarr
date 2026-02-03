package provider

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccTagDataSource(t *testing.T) {
	t.Parallel()

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Unauthorized
			{
				Config:      testAccTagDataSourceConfig("error") + testUnauthorizedProvider,
				ExpectError: regexp.MustCompile("Client Error"),
			},
			// Not found testing
			{
				Config:      testAccTagDataSourceConfig("error"),
				ExpectError: regexp.MustCompile("Unable to find tag"),
			},
			// Create a resource be read
			{
				Config: testAccTagResourceConfig("test", "tag-datasource"),
			},
			// Read testing
			{
				Config: testAccTagResourceConfig("test", "tag-datasource") + testAccTagDataSourceConfig("tag-datasource"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.sonarr_tag.test", "id"),
					resource.TestCheckResourceAttr("data.sonarr_tag.test", "label", "tag-datasource"),
				),
			},
		},
	})
}

func testAccTagDataSourceConfig(label string) string {
	return fmt.Sprintf(`
	data "sonarr_tag" "test" {
		label = "%s"
	}
	`, label)
}

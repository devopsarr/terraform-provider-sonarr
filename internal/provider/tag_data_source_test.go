package provider

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
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
				Config: testAccTagResourceConfig("test", "tag_datasource"),
			},
			// Read testing
			{
				Config: testAccTagResourceConfig("test", "tag_datasource") + testAccTagDataSourceConfig("tag_datasource"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.sonarr_tag.test", "id"),
					resource.TestCheckResourceAttr("data.sonarr_tag.test", "label", "tag_datasource"),
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

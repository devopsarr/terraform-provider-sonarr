package provider

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccReleaseProfileDataSource(t *testing.T) {
	t.Parallel()

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Unauthorized
			{
				Config:      testAccReleaseProfileDataSourceConfig("999") + testUnauthorizedProvider,
				ExpectError: regexp.MustCompile("Client Error"),
			},
			// Read testing
			{
				Config: testAccReleaseProfileResourceConfig("dataSourceTestSingle", "notreally") + testAccReleaseProfileDataSourceConfig("sonarr_release_profile.test.id"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.sonarr_release_profile.test", "id"),
					resource.TestCheckResourceAttr("data.sonarr_release_profile.test", "name", "dataSourceTestSingle")),
			},
		},
	})
}

func testAccReleaseProfileDataSourceConfig(id string) string {
	return fmt.Sprintf(`
data "sonarr_release_profile" "test" {
	id = %s
}
`, id)
}

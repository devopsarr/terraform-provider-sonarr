package provider

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccQualityProfileDataSource(t *testing.T) {
	t.Parallel()

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Unauthorized
			{
				Config:      testAccQualityProfileDataSourceConfig("Error") + testUnauthorizedProvider,
				ExpectError: regexp.MustCompile("Client Error"),
			},
			// Not found testing
			{
				Config:      testAccQualityProfileDataSourceConfig("Error"),
				ExpectError: regexp.MustCompile("Unable to find quality_profile"),
			},
			// Read testing
			{
				Config: testAccQualityProfileDataSourceConfig("Any"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.sonarr_quality_profile.test", "id"),
					resource.TestCheckResourceAttr("data.sonarr_quality_profile.test", "cutoff", "1")),
			},
		},
	})
}

func testAccQualityProfileDataSourceConfig(name string) string {
	return fmt.Sprintf(`
	data "sonarr_quality_profile" "test" {
		name = "%s"
	}
	`, name)
}

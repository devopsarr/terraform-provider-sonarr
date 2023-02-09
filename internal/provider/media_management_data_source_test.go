package provider

import (
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccMediaManagementDataSource(t *testing.T) {
	t.Parallel()

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Unauthorized
			{
				Config:      testAccMediaManagementDataSourceConfig + testUnauthorizedProvider,
				ExpectError: regexp.MustCompile("Client Error"),
			},
			// Read testing
			{
				Config: testAccMediaManagementDataSourceConfig,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.sonarr_media_management.test", "id")),
			},
		},
	})
}

const testAccMediaManagementDataSourceConfig = `
data "sonarr_media_management" "test" {
}
`

package provider

import (
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccDownloadClientsDataSource(t *testing.T) {
	t.Parallel()

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Unauthorized
			{
				Config:      testAccDownloadClientsDataSourceConfig + testUnauthorizedProvider,
				ExpectError: regexp.MustCompile("Client Error"),
			},
			// Create a resource to have a value to check
			{
				Config: testAccDownloadClientResourceConfig("datasourceTest", "false"),
			},
			// Read testing
			{
				Config: testAccDownloadClientsDataSourceConfig,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckTypeSetElemNestedAttrs("data.sonarr_download_clients.test", "download_clients.*", map[string]string{"port": "9091"}),
				),
			},
		},
	})
}

const testAccDownloadClientsDataSourceConfig = `
data "sonarr_download_clients" "test" {
}
`

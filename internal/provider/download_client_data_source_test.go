package provider

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccDownloadClientDataSource(t *testing.T) {
	t.Parallel()

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Unauthorized
			{
				Config:      testAccDownloadClientDataSourceConfig("\"Error\"") + testUnauthorizedProvider,
				ExpectError: regexp.MustCompile("Client Error"),
			},
			// Not found testing
			{
				Config:      testAccDownloadClientDataSourceConfig("\"Error\""),
				ExpectError: regexp.MustCompile("Unable to find download_client"),
			},
			// Read testing
			{
				Config: testAccDownloadClientResourceConfig("dataTest", "true") + testAccDownloadClientDataSourceConfig("sonarr_download_client.test.name"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.sonarr_download_client.test", "id"),
					resource.TestCheckResourceAttr("data.sonarr_download_client.test", "protocol", "torrent")),
			},
		},
	})
}

func testAccDownloadClientDataSourceConfig(name string) string {
	return fmt.Sprintf(`
	data "sonarr_download_client" "test" {
		name = %s
	}
	`, name)
}

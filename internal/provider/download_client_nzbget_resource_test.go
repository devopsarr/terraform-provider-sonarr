package provider

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccDownloadClientNzbgetResource(t *testing.T) {
	t.Parallel()

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read testing
			{
				Config: testAccDownloadClientNzbgetResourceConfig("resourceNzbgetTest", "nzbget"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("sonarr_download_client_nzbget.test", "host", "nzbget"),
					resource.TestCheckResourceAttr("sonarr_download_client_nzbget.test", "url_base", "/nzbget/"),
					resource.TestCheckResourceAttrSet("sonarr_download_client_nzbget.test", "id"),
				),
			},
			// Update and Read testing
			{
				Config: testAccDownloadClientNzbgetResourceConfig("resourceNzbgetTest", "nzbget-host"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("sonarr_download_client_nzbget.test", "host", "nzbget-host"),
				),
			},
			// ImportState testing
			{
				ResourceName:      "sonarr_download_client_nzbget.test",
				ImportState:       true,
				ImportStateVerify: true,
			},
			// Delete testing automatically occurs in TestCase
		},
	})
}

func testAccDownloadClientNzbgetResourceConfig(name, host string) string {
	return fmt.Sprintf(`
	resource "sonarr_download_client_nzbget" "test" {
		enable = false
		priority = 1
		name = "%s"
		host = "%s"
		url_base = "/nzbget/"
		port = 9091
	}`, name, host)
}

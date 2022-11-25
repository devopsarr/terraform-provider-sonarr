package provider

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccDownloadClientUsenetDownloadStationResource(t *testing.T) {
	t.Parallel()

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read testing
			{
				Config: testAccDownloadClientUsenetDownloadStationResourceConfig("resourceUsenetDownloadStationTest", "false"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("sonarr_download_client_usenet_download_station.test", "use_ssl", "false"),
					resource.TestCheckResourceAttrSet("sonarr_download_client_usenet_download_station.test", "id"),
				),
			},
			// Update and Read testing
			{
				Config: testAccDownloadClientUsenetDownloadStationResourceConfig("resourceUsenetDownloadStationTest", "true"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("sonarr_download_client_usenet_download_station.test", "use_ssl", "true"),
				),
			},
			// ImportState testing
			{
				ResourceName:      "sonarr_download_client_usenet_download_station.test",
				ImportState:       true,
				ImportStateVerify: true,
			},
			// Delete testing automatically occurs in TestCase
		},
	})
}

func testAccDownloadClientUsenetDownloadStationResourceConfig(name, ssl string) string {
	return fmt.Sprintf(`
	resource "sonarr_download_client_usenet_download_station" "test" {
		enable = false
		use_ssl = %s
		priority = 1
		name = "%s"
		host = "usenet-download-station"
		port = 9091
	}`, ssl, name)
}

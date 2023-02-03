package provider

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccDownloadClientDelugeResource(t *testing.T) {
	t.Parallel()

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read testing
			{
				Config: testAccDownloadClientDelugeResourceConfig("resourceDelugeTest", "deluge"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("sonarr_download_client_deluge.test", "host", "deluge"),
					resource.TestCheckResourceAttr("sonarr_download_client_deluge.test", "url_base", "/deluge/"),
					resource.TestCheckResourceAttrSet("sonarr_download_client_deluge.test", "id"),
				),
			},
			// Update and Read testing
			{
				Config: testAccDownloadClientDelugeResourceConfig("resourceDelugeTest", "deluge-host"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("sonarr_download_client_deluge.test", "host", "deluge-host"),
				),
			},
			// ImportState testing
			{
				ResourceName:            "sonarr_download_client_deluge.test",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"password"},
			},
			// Delete testing automatically occurs in TestCase
		},
	})
}

func testAccDownloadClientDelugeResourceConfig(name, host string) string {
	return fmt.Sprintf(`
	resource "sonarr_download_client_deluge" "test" {
		enable = false
		priority = 1
		name = "%s"
		host = "%s"
		url_base = "/deluge/"
		port = 9091
		password = ""
	}`, name, host)
}

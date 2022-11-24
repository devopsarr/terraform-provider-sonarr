package provider

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccDownloadClientRtorrentResource(t *testing.T) {
	t.Parallel()

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read testing
			{
				Config: testAccDownloadClientRtorrentResourceConfig("resourceRtorrentTest", "rtorrent"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("sonarr_download_client_rtorrent.test", "host", "rtorrent"),
					resource.TestCheckResourceAttr("sonarr_download_client_rtorrent.test", "url_base", "/rtorrent/"),
					resource.TestCheckResourceAttrSet("sonarr_download_client_rtorrent.test", "id"),
				),
			},
			// Update and Read testing
			{
				Config: testAccDownloadClientRtorrentResourceConfig("resourceRtorrentTest", "rtorrent-host"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("sonarr_download_client_rtorrent.test", "host", "rtorrent-host"),
				),
			},
			// ImportState testing
			{
				ResourceName:      "sonarr_download_client_rtorrent.test",
				ImportState:       true,
				ImportStateVerify: true,
			},
			// Delete testing automatically occurs in TestCase
		},
	})
}

func testAccDownloadClientRtorrentResourceConfig(name, host string) string {
	return fmt.Sprintf(`
	resource "sonarr_download_client_rtorrent" "test" {
		enable = false
		priority = 1
		name = "%s"
		host = "%s"
		url_base = "/rtorrent/"
		port = 9091
	}`, name, host)
}

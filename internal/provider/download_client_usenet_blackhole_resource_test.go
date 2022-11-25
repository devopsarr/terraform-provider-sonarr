package provider

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccDownloadClientUsenetBlackholeResource(t *testing.T) {
	t.Parallel()

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read testing
			{
				Config: testAccDownloadClientUsenetBlackholeResourceConfig("resourceUsenetBlackholeTest", "/config/"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("sonarr_download_client_usenet_blackhole.test", "nzb_folder", "/config/"),
					resource.TestCheckResourceAttr("sonarr_download_client_usenet_blackhole.test", "watch_folder", "/config/"),
					resource.TestCheckResourceAttrSet("sonarr_download_client_usenet_blackhole.test", "id"),
				),
			},
			// Update and Read testing
			{
				Config: testAccDownloadClientUsenetBlackholeResourceConfig("resourceUsenetBlackholeTest", "/config/logs/"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("sonarr_download_client_usenet_blackhole.test", "nzb_folder", "/config/logs/"),
				),
			},
			// ImportState testing
			{
				ResourceName:      "sonarr_download_client_usenet_blackhole.test",
				ImportState:       true,
				ImportStateVerify: true,
			},
			// Delete testing automatically occurs in TestCase
		},
	})
}

func testAccDownloadClientUsenetBlackholeResourceConfig(name, folder string) string {
	return fmt.Sprintf(`
	resource "sonarr_download_client_usenet_blackhole" "test" {
		enable = false
		priority = 1
		name = "%s"
		watch_folder = "/config/"
		nzb_folder = "%s"
	}`, name, folder)
}

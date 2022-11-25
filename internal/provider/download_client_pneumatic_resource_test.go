package provider

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccDownloadClientPneumaticResource(t *testing.T) {
	t.Parallel()

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read testing
			{
				Config: testAccDownloadClientPneumaticResourceConfig("resourcePneumaticTest", "/config/"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("sonarr_download_client_pneumatic.test", "nzb_folder", "/config/"),
					resource.TestCheckResourceAttrSet("sonarr_download_client_pneumatic.test", "id"),
				),
			},
			// Update and Read testing
			{
				Config: testAccDownloadClientPneumaticResourceConfig("resourcePneumaticTest", "/config/logs/"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("sonarr_download_client_pneumatic.test", "nzb_folder", "/config/logs/"),
				),
			},
			// ImportState testing
			{
				ResourceName:      "sonarr_download_client_pneumatic.test",
				ImportState:       true,
				ImportStateVerify: true,
			},
			// Delete testing automatically occurs in TestCase
		},
	})
}

func testAccDownloadClientPneumaticResourceConfig(name, folder string) string {
	return fmt.Sprintf(`
	resource "sonarr_download_client_pneumatic" "test" {
		enable = false
		priority = 1
		name = "%s"
		nzb_folder = "%s"
		strm_folder = "/config/"
	}`, name, folder)
}

package provider

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccDownloadClientResource(t *testing.T) {
	t.Parallel()

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read testing
			{
				Config: testAccDownloadClientResourceConfig("resourceTest", "false"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("sonarr_download_client.test", "enable", "false"),
					resource.TestCheckResourceAttr("sonarr_download_client.test", "url_base", "/transmission/"),
					resource.TestCheckResourceAttrSet("sonarr_download_client.test", "id"),
				),
			},
			// Update and Read testing
			{
				Config: testAccDownloadClientResourceConfig("resourceTest", "true"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("sonarr_download_client.test", "enable", "true"),
				),
			},
			// ImportState testing
			{
				ResourceName:      "sonarr_download_client.test",
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				ResourceName:            "sonarr_download_client.test_sensitive",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"password"},
			},
			// Delete testing automatically occurs in TestCase
		},
	})
}

func testAccDownloadClientResourceConfig(name, enable string) string {
	return fmt.Sprintf(`
	resource "sonarr_download_client" "test" {
		enable = %s
		priority = 1
		name = "%s"
		implementation = "Transmission"
		protocol = "torrent"
    	config_contract = "TransmissionSettings"
		host = "transmission"
		url_base = "/transmission/"
		port = 9091
		tags = []
	}
	
	resource "sonarr_download_client" "test_sensitive" {
		enable = false
		priority = 1
		name = "%sWithSensitive"
		host = "hadouken"
		url_base = "/hadouken/"
		port = 9091
		category = "sonarr-tv"
		username = "username"
		password = "password"
		protocol = "torrent"
    	config_contract = "HadoukenSettings"
		implementation = "Hadouken"
	}
	`, enable, name, name)
}

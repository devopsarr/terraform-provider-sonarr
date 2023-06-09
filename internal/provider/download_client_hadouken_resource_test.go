package provider

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccDownloadClientHadoukenResource(t *testing.T) {
	t.Parallel()

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Unauthorized Create
			{
				Config:      testAccDownloadClientHadoukenResourceConfig("resourceHadoukenTest", "hadouken") + testUnauthorizedProvider,
				ExpectError: regexp.MustCompile("Client Error"),
			},
			// Create and Read testing
			{
				Config: testAccDownloadClientHadoukenResourceConfig("resourceHadoukenTest", "hadouken"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("sonarr_download_client_hadouken.test", "host", "hadouken"),
					resource.TestCheckResourceAttr("sonarr_download_client_hadouken.test", "url_base", "/hadouken/"),
					resource.TestCheckResourceAttrSet("sonarr_download_client_hadouken.test", "id"),
				),
			},
			// Unauthorized Read
			{
				Config:      testAccDownloadClientHadoukenResourceConfig("resourceHadoukenTest", "hadouken") + testUnauthorizedProvider,
				ExpectError: regexp.MustCompile("Client Error"),
			},
			// Update and Read testing
			{
				Config: testAccDownloadClientHadoukenResourceConfig("resourceHadoukenTest", "hadouken-host"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("sonarr_download_client_hadouken.test", "host", "hadouken-host"),
				),
			},
			// ImportState testing
			{
				ResourceName:            "sonarr_download_client_hadouken.test",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"password"},
			},
			// Delete testing automatically occurs in TestCase
		},
	})
}

func testAccDownloadClientHadoukenResourceConfig(name, host string) string {
	return fmt.Sprintf(`
	resource "sonarr_download_client_hadouken" "test" {
		enable = false
		priority = 1
		name = "%s"
		host = "%s"
		url_base = "/hadouken/"
		port = 9091
		category = "sonarr-tv"
		username = "username"
		password = "password"
	}`, name, host)
}

package provider

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccDownloadClientAria2Resource(t *testing.T) {
	t.Parallel()

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Unauthorized Create
			{
				Config:      testAccDownloadClientAria2ResourceConfig("resourceAria2Test", "aria2") + testUnauthorizedProvider,
				ExpectError: regexp.MustCompile("Client Error"),
			},
			// Create and Read testing
			{
				Config: testAccDownloadClientAria2ResourceConfig("resourceAria2Test", "aria2"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("sonarr_download_client_aria2.test", "host", "aria2"),
					resource.TestCheckResourceAttr("sonarr_download_client_aria2.test", "rpc_path", "/aria2/"),
					resource.TestCheckResourceAttrSet("sonarr_download_client_aria2.test", "id"),
				),
			},
			// Unauthorized Read
			{
				Config:      testAccDownloadClientAria2ResourceConfig("resourceAria2Test", "aria2") + testUnauthorizedProvider,
				ExpectError: regexp.MustCompile("Client Error"),
			},
			// Update and Read testing
			{
				Config: testAccDownloadClientAria2ResourceConfig("resourceAria2Test", "aria2-host"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("sonarr_download_client_aria2.test", "host", "aria2-host"),
				),
			},
			// ImportState testing
			{
				ResourceName:      "sonarr_download_client_aria2.test",
				ImportState:       true,
				ImportStateVerify: true,
			},
			// Delete testing automatically occurs in TestCase
		},
	})
}

func testAccDownloadClientAria2ResourceConfig(name, host string) string {
	return fmt.Sprintf(`
	resource "sonarr_download_client_aria2" "test" {
		enable = false
		priority = 1
		name = "%s"
		host = "%s"
		rpc_path = "/aria2/"
		port = 6800
	}`, name, host)
}

package provider

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccDownloadClientConfigResource(t *testing.T) {
	t.Parallel()

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Unauthorized Create
			{
				Config:      testAccDownloadClientConfigResourceConfig("true") + testUnauthorizedProvider,
				ExpectError: regexp.MustCompile("Client Error"),
			},
			// Create and Read testing
			{
				Config: testAccDownloadClientConfigResourceConfig("true"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("sonarr_download_client_config.test", "auto_redownload_failed", "true"),
					resource.TestCheckResourceAttrSet("sonarr_download_client_config.test", "id"),
				),
			},
			// Unauthorized Read
			{
				Config:      testAccDownloadClientConfigResourceConfig("true") + testUnauthorizedProvider,
				ExpectError: regexp.MustCompile("Client Error"),
			},
			// Update and Read testing
			{
				Config: testAccDownloadClientConfigResourceConfig("false"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("sonarr_download_client_config.test", "auto_redownload_failed", "false"),
				),
			},
			// ImportState testing
			{
				ResourceName:      "sonarr_download_client_config.test",
				ImportState:       true,
				ImportStateVerify: true,
			},
			// Delete testing automatically occurs in TestCase
		},
	})
}

func testAccDownloadClientConfigResourceConfig(redownload string) string {
	return fmt.Sprintf(`
	resource "sonarr_download_client_config" "test" {
		enable_completed_download_handling = true
		auto_redownload_failed = %s
	}`, redownload)
}

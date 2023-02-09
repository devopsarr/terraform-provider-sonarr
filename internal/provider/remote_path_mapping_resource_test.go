package provider

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccRemotePathMappingResource(t *testing.T) {
	t.Parallel()

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Unauthorized Create
			{
				Config:      testAccRemotePathMappingResourceConfig("error", "/error/") + testUnauthorizedProvider,
				ExpectError: regexp.MustCompile("Client Error"),
			},
			// Create and Read testing
			{
				Config: testAccRemotePathMappingResourceConfig("remotemapResourceTest", "/test1/"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("sonarr_remote_path_mapping.test", "remote_path", "/test1/"),
					resource.TestCheckResourceAttrSet("sonarr_remote_path_mapping.test", "id"),
				),
			},
			// Unauthorized Read
			{
				Config:      testAccRemotePathMappingResourceConfig("error", "/error/") + testUnauthorizedProvider,
				ExpectError: regexp.MustCompile("Client Error"),
			},
			// Update and Read testing
			{
				Config: testAccRemotePathMappingResourceConfig("remotemapResourceTest", "/test2/"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("sonarr_remote_path_mapping.test", "remote_path", "/test2/"),
				),
			},
			// ImportState testing
			{
				ResourceName:      "sonarr_remote_path_mapping.test",
				ImportState:       true,
				ImportStateVerify: true,
			},
			// Delete testing automatically occurs in TestCase
		},
	})
}

func testAccRemotePathMappingResourceConfig(host, remote string) string {
	return fmt.Sprintf(`
		resource "sonarr_remote_path_mapping" "test" {
  			host = "%s"
			remote_path = "%s"
			local_path = "/config/"
		}
	`, host, remote)
}

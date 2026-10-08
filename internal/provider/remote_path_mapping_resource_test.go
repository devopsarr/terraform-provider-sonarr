package provider

import (
	"context"
	"fmt"
	"net/http"
	"regexp"
	"testing"

	"github.com/devopsarr/sonarr-go/sonarr"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
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

//nolint:paralleltest // deletes an object outside Terraform, a parallel test could otherwise take over its freed ID
func TestAccRemotePathMappingResourceDisappears(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create, then delete outside Terraform: Read removes it from state instead of failing the plan
			{
				Config: testAccRemotePathMappingResourceConfig("remotemapResourceTest", "/test2/"),
				Check: testAccCheckResourceDisappears("sonarr_remote_path_mapping.test", func(client *sonarr.APIClient, id int32) (*http.Response, error) {
					return client.RemotePathMappingAPI.DeleteRemotePathMapping(context.TODO(), id).Execute()
				}),
				ExpectNonEmptyPlan: true,
			},
			// Create again after the deletion outside Terraform
			{
				Config: testAccRemotePathMappingResourceConfig("remotemapResourceTest", "/test2/"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet("sonarr_remote_path_mapping.test", "id"),
				),
			},
		},
	})
}

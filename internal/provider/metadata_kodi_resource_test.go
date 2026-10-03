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

func TestAccMetadataKodiResource(t *testing.T) {
	t.Parallel()

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Unauthorized Create
			{
				Config:      testAccMetadataKodiResourceConfig("kodiResourceTest", "false") + testUnauthorizedProvider,
				ExpectError: regexp.MustCompile("Client Error"),
			},
			// Create and Read testing
			{
				Config: testAccMetadataKodiResourceConfig("kodiResourceTest", "false"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("sonarr_metadata_kodi.test", "series_metadata", "false"),
					resource.TestCheckResourceAttrSet("sonarr_metadata_kodi.test", "id"),
				),
			},
			// Unauthorized Read
			{
				Config:      testAccMetadataKodiResourceConfig("kodiResourceTest", "false") + testUnauthorizedProvider,
				ExpectError: regexp.MustCompile("Client Error"),
			},
			// Update and Read testing
			{
				Config: testAccMetadataKodiResourceConfig("kodiResourceTest", "true"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("sonarr_metadata_kodi.test", "series_metadata", "true"),
				),
			},
			// ImportState testing
			{
				ResourceName:      "sonarr_metadata_kodi.test",
				ImportState:       true,
				ImportStateVerify: true,
			},
			// Delete testing automatically occurs in TestCase
		},
	})
}

func testAccMetadataKodiResourceConfig(name, metadata string) string {
	return fmt.Sprintf(`
	resource "sonarr_metadata_kodi" "test" {
		enable = false
		name = "%s"
		series_metadata = %s
		series_images = true
		episode_images = true
		series_metadata_url = false
		season_images = true
		episode_metadata = false
	}`, name, metadata)
}

//nolint:paralleltest // deletes an object outside Terraform, a parallel test could otherwise take over its freed ID
func TestAccMetadataKodiResourceDisappears(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create, then delete outside Terraform: Read removes it from state instead of failing the plan
			{
				Config: testAccMetadataKodiResourceConfig("kodiResourceTest", "true"),
				Check: testAccCheckResourceDisappears("sonarr_metadata_kodi.test", func(client *sonarr.APIClient, id int32) (*http.Response, error) {
					return client.MetadataAPI.DeleteMetadata(context.TODO(), id).Execute()
				}),
				ExpectNonEmptyPlan: true,
			},
			// Create again after the deletion outside Terraform
			{
				Config: testAccMetadataKodiResourceConfig("kodiResourceTest", "true"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet("sonarr_metadata_kodi.test", "id"),
				),
			},
		},
	})
}

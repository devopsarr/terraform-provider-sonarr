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

func TestAccTagResource(t *testing.T) {
	t.Parallel()

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Unauthorized Create
			{
				Config:      testAccTagResourceConfig("test", "error") + testUnauthorizedProvider,
				ExpectError: regexp.MustCompile("Client Error"),
			},
			// Create and Read testing
			{
				Config: testAccTagResourceConfig("test", "eng"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("sonarr_tag.test", "label", "eng"),
					resource.TestCheckResourceAttrSet("sonarr_tag.test", "id"),
				),
			},
			// Unauthorized Read
			{
				Config:      testAccTagResourceConfig("test", "error") + testUnauthorizedProvider,
				ExpectError: regexp.MustCompile("Client Error"),
				Destroy:     true,
			},
			// Update and Read testing
			{
				Config: testAccTagResourceConfig("test", "1080p"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("sonarr_tag.test", "label", "1080p"),
				),
			},
			// ImportState testing
			{
				ResourceName:      "sonarr_tag.test",
				ImportState:       true,
				ImportStateVerify: true,
			},
			// Delete testing automatically occurs in TestCase
		},
	})
}

func testAccTagResourceConfig(name, label string) string {
	return fmt.Sprintf(`
		resource "sonarr_tag" "%s" {
  			label = "%s"
		}
	`, name, label)
}

//nolint:paralleltest // deletes an object outside Terraform, a parallel test could otherwise take over its freed ID
func TestAccTagResourceDisappears(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create, then delete outside Terraform: Read removes it from state instead of failing the plan
			{
				Config: testAccTagResourceConfig("test", "1080p"),
				Check: testAccCheckResourceDisappears("sonarr_tag.test", func(client *sonarr.APIClient, id int32) (*http.Response, error) {
					return client.TagAPI.DeleteTag(context.TODO(), id).Execute()
				}),
				ExpectNonEmptyPlan: true,
			},
			// Create again after the deletion outside Terraform
			{
				Config: testAccTagResourceConfig("test", "1080p"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet("sonarr_tag.test", "id"),
				),
			},
		},
	})
}

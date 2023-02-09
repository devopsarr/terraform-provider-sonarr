package provider

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccReleaseProfileResource(t *testing.T) {
	t.Parallel()

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Unauthorized Create
			{
				Config:      testAccReleaseProfileResourceConfig("resourceTest", "test1") + testUnauthorizedProvider,
				ExpectError: regexp.MustCompile("Client Error"),
			},
			// Create and Read testing
			{
				Config: testAccReleaseProfileResourceConfig("resourceTest", "test1"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("sonarr_release_profile.test", "required.0", "test1"),
					resource.TestCheckResourceAttrSet("sonarr_release_profile.test", "id"),
				),
			},
			// Unauthorized Read
			{
				Config:      testAccReleaseProfileResourceConfig("resourceTest", "test1") + testUnauthorizedProvider,
				ExpectError: regexp.MustCompile("Client Error"),
			},
			// Update and Read testing
			{
				Config: testAccReleaseProfileResourceConfig("resourceTest", "test2"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("sonarr_release_profile.test", "required.0", "test2"),
				),
			},
			// ImportState testing
			{
				ResourceName:      "sonarr_release_profile.test",
				ImportState:       true,
				ImportStateVerify: true,
			},
			// Delete testing automatically occurs in TestCase
		},
	})
}

func testAccReleaseProfileResourceConfig(name, required string) string {
	return fmt.Sprintf(`
	resource "sonarr_release_profile" "test" {
		name = "%s"
		indexer_id = 0
		required= ["%s"]
	}`, name, required)
}

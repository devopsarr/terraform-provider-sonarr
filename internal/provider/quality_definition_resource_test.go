package provider

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccQualityDefinitionResource(t *testing.T) {
	t.Parallel()

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Unauthorized Create
			{
				Config:      testAccQualityDefinitionResourceConfig("example-4k") + testUnauthorizedProvider,
				ExpectError: regexp.MustCompile("Client Error"),
			},
			// Create and Read testing
			{
				Config: testAccQualityDefinitionResourceConfig("example-4k"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("sonarr_quality_definition.test", "title", "example-4k"),
					resource.TestCheckResourceAttrSet("sonarr_quality_definition.test", "id"),
				),
			},
			// Unauthorized Read
			{
				Config:      testAccQualityDefinitionResourceConfig("example-4k") + testUnauthorizedProvider,
				ExpectError: regexp.MustCompile("Client Error"),
			},
			// Update and Read testing
			{
				Config: testAccQualityDefinitionResourceConfig("example-HD"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("sonarr_quality_definition.test", "title", "example-HD"),
				),
			},
			// ImportState testing
			{
				ResourceName:      "sonarr_quality_definition.test",
				ImportState:       true,
				ImportStateVerify: true,
			},
			// Delete testing automatically occurs in TestCase
		},
	})
}

func testAccQualityDefinitionResourceConfig(name string) string {
	return fmt.Sprintf(`
	resource "sonarr_quality_definition" "test" {
		id = 21
		title    = "%s"
		min_size = 35.0
		max_size = 400
	}
	`, name)
}

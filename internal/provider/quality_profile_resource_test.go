package provider

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccQualityProfileResource(t *testing.T) {
	t.Parallel()

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Unauthorized Create
			{
				Config:      testAccQualityProfileResourceError + testUnauthorizedProvider,
				ExpectError: regexp.MustCompile("Client Error"),
			},
			// Create and Read testing
			{
				Config: testAccQualityProfileResourceConfig("example-4k"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("sonarr_quality_profile.test", "name", "example-4k"),
					resource.TestCheckResourceAttrSet("sonarr_quality_profile.test", "id"),
				),
			},
			// Unauthorized Read
			{
				Config:      testAccQualityProfileResourceError + testUnauthorizedProvider,
				ExpectError: regexp.MustCompile("Client Error"),
			},
			// Update and Read testing
			{
				Config: testAccQualityProfileResourceConfig("example-HD"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("sonarr_quality_profile.test", "name", "example-HD"),
				),
			},
			// ImportState testing
			{
				ResourceName:      "sonarr_quality_profile.test",
				ImportState:       true,
				ImportStateVerify: true,
			},
			// Delete testing automatically occurs in TestCase
		},
	})
}

const testAccQualityProfileResourceError = `
resource "sonarr_quality_profile" "test" {
	name            = "Error"
	upgrade_allowed = true
	cutoff          = 2000
	quality_groups = []
}
`

func testAccQualityProfileResourceConfig(name string) string {
	return fmt.Sprintf(`
	resource "sonarr_custom_format" "test" {
		include_custom_format_when_renaming = false
		name = "QualityFormatTest"

		specifications = [
			{
				name = "Arabic"
				implementation = "LanguageSpecification"
				negate = false
				required = false
				value = "31"
			}
		]
	}

	data "sonarr_quality" "bluray" {
		name = "Bluray-2160p"
	}

	data "sonarr_quality" "webdl" {
		name = "WEBDL-2160p"
	}

	data "sonarr_quality" "webrip" {
		name = "WEBRip-2160p"
	}

	resource "sonarr_quality_profile" "test" {
		name            = "%s"
		upgrade_allowed = true
		cutoff          = 2000

		quality_groups = [
			{
				id   = 2000
				name = "WEB 2160p"
				qualities = [
					data.sonarr_quality.webdl,
					data.sonarr_quality.webrip,
				]
			},
			{
				qualities = [data.sonarr_quality.bluray]
			}
		]

		format_items = [
			{
				name   = sonarr_custom_format.test.name
				format = sonarr_custom_format.test.id
				score  = 10
			}
		]
	}
	`, name)
}

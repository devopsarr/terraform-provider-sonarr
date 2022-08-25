package provider

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccQualityProfileResource(t *testing.T) {
	t.Parallel()

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read testing
			{
				Config: testAccQualityProfileResourceConfig("example-4k"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("sonarr_quality_profile.test", "name", "example-4k"),
					resource.TestCheckResourceAttrSet("sonarr_quality_profile.test", "id"),
				),
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

func testAccQualityProfileResourceConfig(name string) string {
	return fmt.Sprintf(`
	resource "sonarr_quality_profile" "test" {
		name            = "%s"
		upgrade_allowed = true
		cutoff          = 1100

		quality_groups = [
			{
				id   = 1100
				name = "4k"
				qualities = [
					{
						id         = 18
						name       = "WEBDL-2160p"
						source     = "web"
						resolution = 2160
					},
					{
						id         = 19
						name       = "Bluray-2160p"
						source     = "bluray"
						resolution = 2160
					}
				]
			}
		]
	}
	`, name)
}

package provider

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccLanguageProfileResource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read testing
			{
				Config: testAccLanguageProfileResourceConfig("English"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("sonarr_language_profile.test", "cutoff_language", "English"),
					resource.TestCheckResourceAttrSet("sonarr_language_profile.test", "id"),
				),
			},
			// Update and Read testing
			{
				Config: testAccLanguageProfileResourceConfig("Italian"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("sonarr_language_profile.test", "cutoff_language", "Italian"),
				),
			},
			// ImportState testing
			{
				ResourceName:      "sonarr_language_profile.test",
				ImportState:       true,
				ImportStateVerify: true,
			},
			// Delete testing automatically occurs in TestCase
		},
	})
}

func testAccLanguageProfileResourceConfig(cutoff string) string {
	return fmt.Sprintf(`
	resource "sonarr_language_profile" "test" {
		upgrade_allowed = true
		name = "Test"
		cutoff_language = "%s"
		languages = [ "English", "Italian" ]
	}`, cutoff)
}

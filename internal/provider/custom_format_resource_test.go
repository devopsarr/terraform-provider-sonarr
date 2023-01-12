package provider

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccCustomFormatResource(t *testing.T) {
	t.Parallel()

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read testing
			{
				Config: testAccCustomFormatResourceConfig("resourceTest", "false"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("sonarr_custom_format.test", "include_custom_format_when_renaming", "false"),
					resource.TestCheckResourceAttrSet("sonarr_custom_format.test", "id"),
				),
			},
			// Update and Read testing
			{
				Config: testAccCustomFormatResourceConfig("resourceTest", "true"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("sonarr_custom_format.test", "include_custom_format_when_renaming", "true"),
				),
			},
			// ImportState testing
			{
				ResourceName:      "sonarr_custom_format.test",
				ImportState:       true,
				ImportStateVerify: true,
			},
			// Delete testing automatically occurs in TestCase
		},
	})
}

func testAccCustomFormatResourceConfig(name, enable string) string {
	return fmt.Sprintf(`
	resource "sonarr_custom_format" "test" {
		include_custom_format_when_renaming = %s
		name = "%s"
		
		specifications = [
			{
				name = "Surround Sound"
				implementation = "ReleaseTitleSpecification"
				negate = false
				required = false
				value = "DTS.?(HD|ES|X(?!\\D))|TRUEHD|ATMOS|DD(\\+|P).?([5-9])|EAC3.?([5-9])"
			},
			{
				name = "Arabic"
				implementation = "LanguageSpecification"
				negate = false
				required = false
				value = "31"
			}
		]	
	}`, enable, name)
}

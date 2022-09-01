package provider

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccLanguageProfileDataSource(t *testing.T) {
	t.Parallel()

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Read testing
			{
				Config: testAccLanguageProfileDataSourceConfig,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.sonarr_language_profile.test", "id"),
					resource.TestCheckResourceAttr("data.sonarr_language_profile.test", "cutoff_language", "Arabic")),
			},
		},
	})
}

const testAccLanguageProfileDataSourceConfig = `
resource "sonarr_language_profile" "test" {
	upgrade_allowed = true
	name = "lpdata"
	cutoff_language = "Arabic"
	languages = [ "English", "Italian", "Arabic" ]
}

data "sonarr_language_profile" "test" {
	name = sonarr_language_profile.test.name
}
`

package provider

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccLanguageProfilesDataSource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create a delay profile to have a value to check
			{
				Config: testAccLanguageProfileResourceConfig("English"),
			},
			// Read testing
			{
				Config: testAccLanguageProfilesDataSourceConfig,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckTypeSetElemNestedAttrs("data.sonarr_language_profiles.test", "language_profiles.*", map[string]string{"name": "Test"}),
				),
			},
		},
	})
}

const testAccLanguageProfilesDataSourceConfig = `
data "sonarr_language_profiles" "test" {
}
`

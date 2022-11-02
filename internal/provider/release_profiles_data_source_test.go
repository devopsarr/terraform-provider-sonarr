package provider

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccReleaseProfilesDataSource(t *testing.T) {
	t.Parallel()

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create a profile to have a value to check
			{
				Config: testAccReleaseProfileResourceConfig("testDataSources", "sd"),
			},
			// Read testing
			{
				Config: testAccReleaseProfilesDataSourceConfig,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckTypeSetElemNestedAttrs("data.sonarr_release_profiles.test", "release_profiles.*", map[string]string{"name": "testDataSources"}),
				),
			},
		},
	})
}

const testAccReleaseProfilesDataSourceConfig = `
data "sonarr_release_profiles" "test" {
}
`

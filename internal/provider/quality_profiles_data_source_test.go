package provider

import (
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"golift.io/starr"
	"golift.io/starr/sonarr"
)

func TestAccQualityProfilesDataSource(t *testing.T) {
	t.Parallel()

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Read testing
			{
				PreConfig: qualityprofilesDSInit,
				Config:    testAccQualityProfilesDataSourceConfig,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckTypeSetElemNestedAttrs("data.sonarr_quality_profiles.test", "quality_profiles.*", map[string]string{"name": "Any"}),
				),
			},
		},
	})
}

const testAccQualityProfilesDataSourceConfig = `
data "sonarr_quality_profiles" "test" {
}
`

func qualityprofilesDSInit() {
	// keep only first two profiles to avoid longer tests
	client := *sonarr.New(starr.New(os.Getenv("SONARR_API_KEY"), os.Getenv("SONARR_URL"), 0))
	for i := 3; i < 7; i++ {
		_ = client.DeleteQualityProfile(i)
	}
}

package provider

import (
	"context"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccQualityProfilesDataSource(t *testing.T) {
	t.Parallel()

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Unauthorized
			{
				Config:      testAccQualityProfilesDataSourceConfig + testUnauthorizedProvider,
				ExpectError: regexp.MustCompile("Client Error"),
			},
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
	client := testAccAPIClient()
	for i := 3; i < 7; i++ {
		_, _ = client.QualityProfileApi.DeleteQualityProfile(context.TODO(), int32(i)).Execute()
	}
}

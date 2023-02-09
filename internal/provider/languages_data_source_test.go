package provider

import (
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccLanguagesDataSource(t *testing.T) {
	t.Parallel()

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Unauthorized
			{
				Config:      testAccLanguagesDataSourceConfig + testUnauthorizedProvider,
				ExpectError: regexp.MustCompile("Client Error"),
			},
			// Read testing
			{
				Config: testAccLanguagesDataSourceConfig,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckTypeSetElemNestedAttrs("data.sonarr_languages.test", "languages.*", map[string]string{"name": "English"}),
				),
			},
		},
	})
}

const testAccLanguagesDataSourceConfig = `
data "sonarr_languages" "test" {
}
`

package provider

import (
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccLanguageDataSource(t *testing.T) {
	t.Parallel()

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Unauthorized
			{
				Config:      testAccLanguageDataSourceConfig + testUnauthorizedProvider,
				ExpectError: regexp.MustCompile("Client Error"),
			},
			// Read testing
			{
				Config: testAccLanguageDataSourceConfig,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.sonarr_language.test", "id"),
					resource.TestCheckResourceAttr("data.sonarr_language.test", "name_lower", "english"),
				),
			},
		},
	})
}

const testAccLanguageDataSourceConfig = `
data "sonarr_language" "test" {
	name = "English"
}
`

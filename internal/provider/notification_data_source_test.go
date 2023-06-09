package provider

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccNotificationDataSource(t *testing.T) {
	t.Parallel()

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Unauthorized
			{
				Config:      testAccNotificationDataSourceConfig("\"Error\"") + testUnauthorizedProvider,
				ExpectError: regexp.MustCompile("Client Error"),
			},
			// Not found testing
			{
				Config:      testAccNotificationDataSourceConfig("\"Error\""),
				ExpectError: regexp.MustCompile("Unable to find notification"),
			},
			// Read testing
			{
				Config: testAccNotificationResourceConfig("dataTest", "true") + testAccNotificationDataSourceConfig("sonarr_notification.test.name"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.sonarr_notification.test", "id"),
					resource.TestCheckResourceAttr("data.sonarr_notification.test", "path", "/scripts/test.sh")),
			},
		},
	})
}

func testAccNotificationDataSourceConfig(name string) string {
	return fmt.Sprintf(`
	data "sonarr_notification" "test" {
		name = %s
	}
	`, name)
}

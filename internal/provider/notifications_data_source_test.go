package provider

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccNotificationsDataSource(t *testing.T) {
	t.Parallel()

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create a delay profile to have a value to check
			{
				Config: testAccNotificationResourceConfig("datasourceTest", "true"),
			},
			// Read testing
			{
				Config: testAccNotificationsDataSourceConfig,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckTypeSetElemNestedAttrs("data.sonarr_notifications.test", "notifications.*", map[string]string{"path": "/defaults/test.sh"}),
				),
			},
		},
	})
}

const testAccNotificationsDataSourceConfig = `
data "sonarr_notifications" "test" {
}
`

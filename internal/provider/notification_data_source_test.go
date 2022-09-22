package provider

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccNotificationDataSource(t *testing.T) {
	t.Parallel()

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Read testing
			{
				Config: testAccNotificationDataSourceConfig,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.sonarr_notification.test", "id"),
					resource.TestCheckResourceAttr("data.sonarr_notification.test", "path", "/scripts/test.sh")),
			},
		},
	})
}

const testAccNotificationDataSourceConfig = `
resource "sonarr_notification" "test" {
	on_grab                            = false
	on_download                        = true
	on_upgrade                         = true
	on_rename                          = false
	on_series_delete                   = false
	on_episode_file_delete             = false
	on_episode_file_delete_for_upgrade = true
	on_health_issue                    = false
	on_application_update              = false
  
	include_health_warnings = false
	name                    = "notificationData"
  
	implementation  = "CustomScript"
	config_contract = "CustomScriptSettings"
  
	path = "/scripts/test.sh"
}

data "sonarr_notification" "test" {
	name = sonarr_notification.test.name
}
`

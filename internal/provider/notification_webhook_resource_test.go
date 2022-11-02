package provider

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccNotificationWebhookResource(t *testing.T) {
	t.Parallel()

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read testing
			{
				Config: testAccNotificationWebhookResourceConfig("resourceWebhookTest", "false"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("sonarr_notification_webhook.test", "on_upgrade", "false"),
					resource.TestCheckResourceAttrSet("sonarr_notification_webhook.test", "id"),
				),
			},
			// Update and Read testing
			{
				Config: testAccNotificationWebhookResourceConfig("resourceWebhookTest", "true"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("sonarr_notification_webhook.test", "on_upgrade", "true"),
				),
			},
			// ImportState testing
			{
				ResourceName:      "sonarr_notification_webhook.test",
				ImportState:       true,
				ImportStateVerify: true,
			},
			// Delete testing automatically occurs in TestCase
		},
	})
}

func testAccNotificationWebhookResourceConfig(name, upgrade string) string {
	return fmt.Sprintf(`
	resource "sonarr_notification_webhook" "test" {
		on_grab                            = false
		on_download                        = true
		on_upgrade                         = %s
		on_rename                          = false
		on_series_delete                   = false
		on_episode_file_delete             = false
		on_episode_file_delete_for_upgrade = true
		on_health_issue                    = false
		on_application_update              = false
	  
		include_health_warnings = false
		name                    = "%s"
	  
		url = "http://transmission:9091"
		method = 1
	}`, upgrade, name)
}

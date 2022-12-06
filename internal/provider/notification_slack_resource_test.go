package provider

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccNotificationSlackResource(t *testing.T) {
	t.Parallel()

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read testing
			{
				Config: testAccNotificationSlackResourceConfig("resourceSlackTest", "test"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("sonarr_notification_slack.test", "channel", "test"),
					resource.TestCheckResourceAttrSet("sonarr_notification_slack.test", "id"),
				),
			},
			// Update and Read testing
			{
				Config: testAccNotificationSlackResourceConfig("resourceSlackTest", "test1"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("sonarr_notification_slack.test", "channel", "test1"),
				),
			},
			// ImportState testing
			{
				ResourceName:      "sonarr_notification_slack.test",
				ImportState:       true,
				ImportStateVerify: true,
			},
			// Delete testing automatically occurs in TestCase
		},
	})
}

func testAccNotificationSlackResourceConfig(name, channel string) string {
	return fmt.Sprintf(`
	resource "sonarr_notification_slack" "test" {
		on_grab                            = false
		on_download                        = false
		on_upgrade                         = false
		on_rename                          = false
		on_series_delete                   = false
		on_episode_file_delete             = false
		on_episode_file_delete_for_upgrade = false
		on_health_issue                    = false
		on_application_update              = false
	  
		include_health_warnings = false
		name                    = "%s"
	  
		web_hook_url = "http://my.slack.com/test"
		username = "user"
		channel = "%s"
	}`, name, channel)
}

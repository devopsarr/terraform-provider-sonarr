package provider

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccNotificationPushoverResource(t *testing.T) {
	t.Parallel()

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Unauthorized Create
			{
				Config:      testAccNotificationPushoverResourceConfig("resourcePushoverTest", 0) + testUnauthorizedProvider,
				ExpectError: regexp.MustCompile("Client Error"),
			},
			// Create and Read testing
			{
				Config: testAccNotificationPushoverResourceConfig("resourcePushoverTest", 0),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("sonarr_notification_pushover.test", "priority", "0"),
					resource.TestCheckResourceAttrSet("sonarr_notification_pushover.test", "id"),
				),
			},
			// Unauthorized Read
			{
				Config:      testAccNotificationPushoverResourceConfig("resourcePushoverTest", 0) + testUnauthorizedProvider,
				ExpectError: regexp.MustCompile("Client Error"),
			},
			// Update and Read testing
			{
				Config: testAccNotificationPushoverResourceConfig("resourcePushoverTest", 2),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("sonarr_notification_pushover.test", "priority", "2"),
				),
			},
			// ImportState testing
			{
				ResourceName:            "sonarr_notification_pushover.test",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"api_key"},
			},
			// Delete testing automatically occurs in TestCase
		},
	})
}

func testAccNotificationPushoverResourceConfig(name string, priority int) string {
	return fmt.Sprintf(`
	resource "sonarr_notification_pushover" "test" {
		on_grab                            = false
		on_download                        = false
		on_upgrade                         = false
		on_series_delete                   = false
		on_episode_file_delete             = false
		on_episode_file_delete_for_upgrade = false
		on_health_issue                    = false
		on_application_update              = false
	  
		include_health_warnings = false
		name                    = "%s"
	  
		api_key = "Key"
		user_key = "userKey"
		priority = %d
	}`, name, priority)
}

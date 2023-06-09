package provider

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccNotificationTelegramResource(t *testing.T) {
	t.Parallel()

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Unauthorized Create
			{
				Config:      testAccNotificationTelegramResourceConfig("resourceTelegramTest", "chat01") + testUnauthorizedProvider,
				ExpectError: regexp.MustCompile("Client Error"),
			},
			// Create and Read testing
			{
				Config: testAccNotificationTelegramResourceConfig("resourceTelegramTest", "chat01"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("sonarr_notification_telegram.test", "chat_id", "chat01"),
					resource.TestCheckResourceAttrSet("sonarr_notification_telegram.test", "id"),
				),
			},
			// Unauthorized Read
			{
				Config:      testAccNotificationTelegramResourceConfig("resourceTelegramTest", "chat01") + testUnauthorizedProvider,
				ExpectError: regexp.MustCompile("Client Error"),
			},
			// Update and Read testing
			{
				Config: testAccNotificationTelegramResourceConfig("resourceTelegramTest", "chat02"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("sonarr_notification_telegram.test", "chat_id", "chat02"),
				),
			},
			// ImportState testing
			{
				ResourceName:            "sonarr_notification_telegram.test",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"bot_token"},
			},
			// Delete testing automatically occurs in TestCase
		},
	})
}

func testAccNotificationTelegramResourceConfig(name, chat string) string {
	return fmt.Sprintf(`
	resource "sonarr_notification_telegram" "test" {
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
	  
		chat_id = "%s"
		bot_token = "Token"
	}`, name, chat)
}

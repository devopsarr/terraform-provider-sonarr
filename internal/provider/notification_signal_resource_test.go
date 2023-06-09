package provider

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccNotificationSignalResource(t *testing.T) {
	t.Parallel()

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Unauthorized Create
			{
				Config:      testAccNotificationSignalResourceConfig("resourceSignalTest", "token123") + testUnauthorizedProvider,
				ExpectError: regexp.MustCompile("Client Error"),
			},
			// Create and Read testing
			{
				Config: testAccNotificationSignalResourceConfig("resourceSignalTest", "token123"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("sonarr_notification_signal.test", "receiver_id", "token123"),
					resource.TestCheckResourceAttrSet("sonarr_notification_signal.test", "id"),
				),
			},
			// Unauthorized Read
			{
				Config:      testAccNotificationSignalResourceConfig("resourceSignalTest", "token123") + testUnauthorizedProvider,
				ExpectError: regexp.MustCompile("Client Error"),
			},
			// Update and Read testing
			{
				Config: testAccNotificationSignalResourceConfig("resourceSignalTest", "token234"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("sonarr_notification_signal.test", "receiver_id", "token234"),
				),
			},
			// ImportState testing
			{
				ResourceName:            "sonarr_notification_signal.test",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"auth_password", "sender_number"},
			},
			// Delete testing automatically occurs in TestCase
		},
	})
}

func testAccNotificationSignalResourceConfig(name, receiver string) string {
	return fmt.Sprintf(`
	resource "sonarr_notification_signal" "test" {
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
	  
		auth_username = "User"
		auth_password = "Passwordn"

		host = "localhost"
		port = 8080
		use_ssl = true
		sender_number = "1234"
		receiver_id = "%s"
	}`, name, receiver)
}

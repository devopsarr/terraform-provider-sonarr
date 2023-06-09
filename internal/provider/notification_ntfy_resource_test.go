package provider

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccNotificationNtfyResource(t *testing.T) {
	t.Parallel()

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Unauthorized Create
			{
				Config:      testAccNotificationNtfyResourceConfig("resourceNtfyTest", "token123") + testUnauthorizedProvider,
				ExpectError: regexp.MustCompile("Client Error"),
			},
			// Create and Read testing
			{
				Config: testAccNotificationNtfyResourceConfig("resourceNtfyTest", "token123"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("sonarr_notification_ntfy.test", "password", "token123"),
					resource.TestCheckResourceAttrSet("sonarr_notification_ntfy.test", "id"),
				),
			},
			// Unauthorized Read
			{
				Config:      testAccNotificationNtfyResourceConfig("resourceNtfyTest", "token123") + testUnauthorizedProvider,
				ExpectError: regexp.MustCompile("Client Error"),
			},
			// Update and Read testing
			{
				Config: testAccNotificationNtfyResourceConfig("resourceNtfyTest", "token234"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("sonarr_notification_ntfy.test", "password", "token234"),
				),
			},
			// ImportState testing
			{
				ResourceName:            "sonarr_notification_ntfy.test",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"password"},
			},
			// Delete testing automatically occurs in TestCase
		},
	})
}

func testAccNotificationNtfyResourceConfig(name, token string) string {
	return fmt.Sprintf(`
	resource "sonarr_notification_ntfy" "test" {
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
	  
		priority = 1
		server_url = "https://ntfy.sh"
		username = "User"
		password = "%s"
		topics = ["Topic1234","Topic4321"]
		field_tags = ["warning","skull"]
	}`, name, token)
}

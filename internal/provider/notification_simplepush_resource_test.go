package provider

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccNotificationSimplepushResource(t *testing.T) {
	t.Parallel()

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Unauthorized Create
			{
				Config:      testAccNotificationSimplepushResourceConfig("resourceSimplepushTest", "ringtone:default") + testUnauthorizedProvider,
				ExpectError: regexp.MustCompile("Client Error"),
			},
			// Create and Read testing
			{
				Config: testAccNotificationSimplepushResourceConfig("resourceSimplepushTest", "ringtone:default"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("sonarr_notification_simplepush.test", "event", "ringtone:default"),
					resource.TestCheckResourceAttrSet("sonarr_notification_simplepush.test", "id"),
				),
			},
			// Unauthorized Read
			{
				Config:      testAccNotificationSimplepushResourceConfig("resourceSimplepushTest", "ringtone:default") + testUnauthorizedProvider,
				ExpectError: regexp.MustCompile("Client Error"),
			},
			// Update and Read testing
			{
				Config: testAccNotificationSimplepushResourceConfig("resourceSimplepushTest", "ringtone:special"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("sonarr_notification_simplepush.test", "event", "ringtone:special"),
				),
			},
			// ImportState testing
			{
				ResourceName:            "sonarr_notification_simplepush.test",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"key"},
			},
			// Delete testing automatically occurs in TestCase
		},
	})
}

func testAccNotificationSimplepushResourceConfig(name, event string) string {
	return fmt.Sprintf(`
	resource "sonarr_notification_simplepush" "test" {
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
	  
		key = "Key"
		event = "%s"
	}`, name, event)
}

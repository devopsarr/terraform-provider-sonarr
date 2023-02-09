package provider

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccNotificationKodiResource(t *testing.T) {
	t.Parallel()

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Unauthorized Create
			{
				Config:      testAccNotificationKodiResourceConfig("resourceKodiTest", "pass1") + testUnauthorizedProvider,
				ExpectError: regexp.MustCompile("Client Error"),
			},
			// Create and Read testing
			{
				Config: testAccNotificationKodiResourceConfig("resourceKodiTest", "pass1"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("sonarr_notification_kodi.test", "password", "pass1"),
					resource.TestCheckResourceAttrSet("sonarr_notification_kodi.test", "id"),
				),
			},
			// Unauthorized Read
			{
				Config:      testAccNotificationKodiResourceConfig("resourceKodiTest", "pass1") + testUnauthorizedProvider,
				ExpectError: regexp.MustCompile("Client Error"),
			},
			// Update and Read testing
			{
				Config: testAccNotificationKodiResourceConfig("resourceKodiTest", "pass2"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("sonarr_notification_kodi.test", "password", "pass2"),
				),
			},
			// ImportState testing
			{
				ResourceName:            "sonarr_notification_kodi.test",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"password"},
			},
			// Delete testing automatically occurs in TestCase
		},
	})
}

func testAccNotificationKodiResourceConfig(name, avatar string) string {
	return fmt.Sprintf(`
	resource "sonarr_notification_kodi" "test" {
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
	  
		host = "http://kodi.com"
		port = 8080
		username = "User"
		password = "%s"
		notify = true
	}`, name, avatar)
}

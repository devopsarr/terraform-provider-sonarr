package provider

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccNotificationTraktResource(t *testing.T) {
	t.Parallel()

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Unauthorized Create
			{
				Config:      testAccNotificationTraktResourceConfig("resourceTraktTest", "token123") + testUnauthorizedProvider,
				ExpectError: regexp.MustCompile("Client Error"),
			},
			// Create and Read testing
			{
				Config: testAccNotificationTraktResourceConfig("resourceTraktTest", "token123"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("sonarr_notification_trakt.test", "access_token", "token123"),
					resource.TestCheckResourceAttrSet("sonarr_notification_trakt.test", "id"),
				),
			},
			// Unauthorized Read
			{
				Config:      testAccNotificationTraktResourceConfig("resourceTraktTest", "token123") + testUnauthorizedProvider,
				ExpectError: regexp.MustCompile("Client Error"),
			},
			// Update and Read testing
			{
				Config: testAccNotificationTraktResourceConfig("resourceTraktTest", "token234"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("sonarr_notification_trakt.test", "access_token", "token234"),
				),
			},
			// ImportState testing
			{
				ResourceName:            "sonarr_notification_trakt.test",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"access_token"},
			},
			// Delete testing automatically occurs in TestCase
		},
	})
}

func testAccNotificationTraktResourceConfig(name, token string) string {
	return fmt.Sprintf(`
	resource "sonarr_notification_trakt" "test" {
		on_download                        = false
		on_upgrade                         = false
		on_series_delete                   = false
		on_episode_file_delete             = false
		on_episode_file_delete_for_upgrade = false
	  
		include_health_warnings = false
		name                    = "%s"
	  
		auth_user = "User"
		access_token = "%s"
	}`, name, token)
}

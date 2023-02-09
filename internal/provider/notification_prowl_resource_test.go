package provider

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccNotificationProwlResource(t *testing.T) {
	t.Parallel()

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Unauthorized Create
			{
				Config:      testAccNotificationProwlResourceConfig("resourceProwlTest", 0) + testUnauthorizedProvider,
				ExpectError: regexp.MustCompile("Client Error"),
			},
			// Create and Read testing
			{
				Config: testAccNotificationProwlResourceConfig("resourceProwlTest", 0),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("sonarr_notification_prowl.test", "priority", "0"),
					resource.TestCheckResourceAttrSet("sonarr_notification_prowl.test", "id"),
				),
			},
			// Unauthorized Read
			{
				Config:      testAccNotificationProwlResourceConfig("resourceProwlTest", 0) + testUnauthorizedProvider,
				ExpectError: regexp.MustCompile("Client Error"),
			},
			// Update and Read testing
			{
				Config: testAccNotificationProwlResourceConfig("resourceProwlTest", 2),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("sonarr_notification_prowl.test", "priority", "2"),
				),
			},
			// ImportState testing
			{
				ResourceName:            "sonarr_notification_prowl.test",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"api_key"},
			},
			// Delete testing automatically occurs in TestCase
		},
	})
}

func testAccNotificationProwlResourceConfig(name string, priority int) string {
	return fmt.Sprintf(`
	resource "sonarr_notification_prowl" "test" {
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
		priority = %d
	}`, name, priority)
}

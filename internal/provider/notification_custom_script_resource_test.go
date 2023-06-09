package provider

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccNotificationCustomScriptResource(t *testing.T) {
	t.Parallel()

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Unauthorized Create
			{
				Config:      testAccNotificationCustomScriptResourceConfig("resourceScriptTest", "false") + testUnauthorizedProvider,
				ExpectError: regexp.MustCompile("Client Error"),
			},
			// Create and Read testing
			{
				Config: testAccNotificationCustomScriptResourceConfig("resourceScriptTest", "false"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("sonarr_notification_custom_script.test", "on_upgrade", "false"),
					resource.TestCheckResourceAttrSet("sonarr_notification_custom_script.test", "id"),
				),
			},
			// Unauthorized Read
			{
				Config:      testAccNotificationCustomScriptResourceConfig("resourceScriptTest", "false") + testUnauthorizedProvider,
				ExpectError: regexp.MustCompile("Client Error"),
			},
			// Update and Read testing
			{
				Config: testAccNotificationCustomScriptResourceConfig("resourceScriptTest", "true"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("sonarr_notification_custom_script.test", "on_upgrade", "true"),
				),
			},
			// ImportState testing
			{
				ResourceName:      "sonarr_notification_custom_script.test",
				ImportState:       true,
				ImportStateVerify: true,
			},
			// Delete testing automatically occurs in TestCase
		},
	})
}

func testAccNotificationCustomScriptResourceConfig(name, upgrade string) string {
	return fmt.Sprintf(`
	resource "sonarr_notification_custom_script" "test" {
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
	  
		path = "/scripts/test.sh"
	}`, upgrade, name)
}

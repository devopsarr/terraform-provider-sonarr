package provider

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccNotificationResource(t *testing.T) {
	t.Parallel()

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read testing
			{
				Config: testAccNotificationResourceConfig("resourceTest", "false"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("sonarr_notification.test", "on_upgrade", "false"),
					resource.TestCheckResourceAttrSet("sonarr_notification.test", "id"),
				),
			},
			// Update and Read testing
			{
				Config: testAccNotificationResourceConfig("resourceTest", "true"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("sonarr_notification.test", "on_upgrade", "true"),
				),
			},
			// ImportState testing
			{
				ResourceName:      "sonarr_notification.test",
				ImportState:       true,
				ImportStateVerify: true,
			},
			// Delete testing automatically occurs in TestCase
		},
	})
}

func testAccNotificationResourceConfig(name, upgrade string) string {
	return fmt.Sprintf(`
	resource "sonarr_notification" "test" {
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
	  
		implementation  = "CustomScript"
		config_contract = "CustomScriptSettings"
	  
		path = "/scripts/test.sh"
	}`, upgrade, name)
}

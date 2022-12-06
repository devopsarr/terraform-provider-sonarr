package provider

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccNotificationEmailResource(t *testing.T) {
	t.Parallel()

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read testing
			{
				Config: testAccNotificationEmailResourceConfig("resourceEmailTest", "test@email.com"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("sonarr_notification_email.test", "from", "test@email.com"),
					resource.TestCheckResourceAttrSet("sonarr_notification_email.test", "id"),
				),
			},
			// Update and Read testing
			{
				Config: testAccNotificationEmailResourceConfig("resourceEmailTest", "test123@email.com"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("sonarr_notification_email.test", "from", "test123@email.com"),
				),
			},
			// ImportState testing
			{
				ResourceName:      "sonarr_notification_email.test",
				ImportState:       true,
				ImportStateVerify: true,
			},
			// Delete testing automatically occurs in TestCase
		},
	})
}

func testAccNotificationEmailResourceConfig(name, from string) string {
	return fmt.Sprintf(`
	resource "sonarr_notification_email" "test" {
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
	  
		server = "http://email-server.net"
		port = 587
		from = "%s"
		to = ["test@test.com", "test1@test.com"]
	}`, name, from)
}

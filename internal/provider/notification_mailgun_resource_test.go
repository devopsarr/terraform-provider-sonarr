package provider

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccNotificationMailgunResource(t *testing.T) {
	t.Parallel()

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Unauthorized Create
			{
				Config:      testAccNotificationMailgunResourceConfig("resourceMailgunTest", "test@mailgun.com") + testUnauthorizedProvider,
				ExpectError: regexp.MustCompile("Client Error"),
			},
			// Create and Read testing
			{
				Config: testAccNotificationMailgunResourceConfig("resourceMailgunTest", "test@mailgun.com"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("sonarr_notification_mailgun.test", "from", "test@mailgun.com"),
					resource.TestCheckResourceAttrSet("sonarr_notification_mailgun.test", "id"),
				),
			},
			// Unauthorized Read
			{
				Config:      testAccNotificationMailgunResourceConfig("resourceMailgunTest", "test@mailgun.com") + testUnauthorizedProvider,
				ExpectError: regexp.MustCompile("Client Error"),
			},
			// Update and Read testing
			{
				Config: testAccNotificationMailgunResourceConfig("resourceMailgunTest", "test123@mailgun.com"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("sonarr_notification_mailgun.test", "from", "test123@mailgun.com"),
				),
			},
			// ImportState testing
			{
				ResourceName:            "sonarr_notification_mailgun.test",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"api_key"},
			},
			// Delete testing automatically occurs in TestCase
		},
	})
}

func testAccNotificationMailgunResourceConfig(name, from string) string {
	return fmt.Sprintf(`
	resource "sonarr_notification_mailgun" "test" {
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
		
		api_key = "APIkey"
		from = "%s"
		recipients = ["test@test.com", "test1@test.com"]
	}`, name, from)
}

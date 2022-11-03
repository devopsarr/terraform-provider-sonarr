package provider

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccDelayProfileResource(t *testing.T) {
	t.Parallel()

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read testing
			{
				Config: testAccDelayProfileResourceConfig("usenet"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("sonarr_delay_profile.test", "preferred_protocol", "usenet"),
					resource.TestCheckResourceAttrSet("sonarr_delay_profile.test", "id"),
				),
			},
			// Update and Read testing
			{
				Config: testAccDelayProfileResourceConfig("torrent"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("sonarr_delay_profile.test", "preferred_protocol", "torrent"),
				),
			},
			// ImportState testing
			{
				ResourceName:      "sonarr_delay_profile.test",
				ImportState:       true,
				ImportStateVerify: true,
			},
			// Delete testing automatically occurs in TestCase
		},
	})
}

func testAccDelayProfileResourceConfig(protocol string) string {
	return fmt.Sprintf(`
	resource "sonarr_tag" "test" {
		label = "delay_profile_resource"
	}

	resource "sonarr_delay_profile" "test" {
		enable_usenet = true
		enable_torrent = true
		bypass_if_highest_quality = true
		order = 100
		usenet_delay = 0
		torrent_delay = 0
		preferred_protocol= "%s"
		tags = [sonarr_tag.test.id]
	}`, protocol)
}

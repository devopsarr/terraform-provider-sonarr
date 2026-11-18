package provider

import (
	"context"
	"fmt"
	"net/http"
	"regexp"
	"testing"

	"github.com/devopsarr/sonarr-go/sonarr"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccDelayProfileResource(t *testing.T) {
	t.Parallel()

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Unauthorized Create
			{
				Config:      testAccDelayProfileResourceConfig("usenet", "0") + testUnauthorizedProvider,
				ExpectError: regexp.MustCompile("Client Error"),
			},
			// Create and Read testing
			{
				Config: testAccTagResourceConfig("test", "delay-profile-resource") + testAccDelayProfileResourceConfig("usenet", "sonarr_tag.test.id"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("sonarr_delay_profile.test", "preferred_protocol", "usenet"),
					resource.TestCheckResourceAttrSet("sonarr_delay_profile.test", "id"),
				),
			},
			// Unauthorized Read
			{
				Config:      testAccDelayProfileResourceConfig("usenet", "0") + testUnauthorizedProvider,
				ExpectError: regexp.MustCompile("Client Error"),
			},
			// Update and Read testing
			{
				Config: testAccTagResourceConfig("test", "delay-profile-resource") + testAccDelayProfileResourceConfig("torrent", "sonarr_tag.test.id"),
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

func testAccDelayProfileResourceConfig(protocol, tag string) string {
	return fmt.Sprintf(`
	resource "sonarr_delay_profile" "test" {
		enable_usenet = true
		enable_torrent = true
		bypass_if_highest_quality = true
		bypass_if_above_custom_format_score = true
		order = 100
		usenet_delay = 0
		torrent_delay = 0
		preferred_protocol= "%s"
		tags = [%s]
	}`, protocol, tag)
}

//nolint:paralleltest // deletes an object outside Terraform, a parallel test could otherwise take over its freed ID
func TestAccDelayProfileResourceDisappears(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create, then delete outside Terraform: Read removes it from state instead of failing the plan
			{
				Config: testAccTagResourceConfig("test", "delay-profile-resource") + testAccDelayProfileResourceConfig("torrent", "sonarr_tag.test.id"),
				Check: testAccCheckResourceDisappears("sonarr_delay_profile.test", func(client *sonarr.APIClient, id int32) (*http.Response, error) {
					return client.DelayProfileAPI.DeleteDelayProfile(context.TODO(), id).Execute()
				}),
				ExpectNonEmptyPlan: true,
			},
			// Create again after the deletion outside Terraform
			{
				Config: testAccTagResourceConfig("test", "delay-profile-resource") + testAccDelayProfileResourceConfig("torrent", "sonarr_tag.test.id"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet("sonarr_delay_profile.test", "id"),
				),
			},
		},
	})
}

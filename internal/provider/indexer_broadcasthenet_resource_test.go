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

func TestAccIndexerBroadcastheNetResource(t *testing.T) {
	t.Parallel()

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Unauthorized Create
			{
				Config:      testAccIndexerBroadcastheNetResourceConfig("broadcasthenetResourceTest", 1) + testUnauthorizedProvider,
				ExpectError: regexp.MustCompile("Client Error"),
			},
			// Create and Read testing
			{
				Config: testAccIndexerBroadcastheNetResourceConfig("broadcasthenetResourceTest", 1),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("sonarr_indexer_broadcasthenet.test", "seed_time", "1"),
					resource.TestCheckResourceAttrSet("sonarr_indexer_broadcasthenet.test", "id"),
				),
			},
			// Unauthorized Read
			{
				Config:      testAccIndexerBroadcastheNetResourceConfig("broadcasthenetResourceTest", 1) + testUnauthorizedProvider,
				ExpectError: regexp.MustCompile("Client Error"),
			},
			// Update and Read testing
			{
				Config: testAccIndexerBroadcastheNetResourceConfig("broadcasthenetResourceTest", 2),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("sonarr_indexer_broadcasthenet.test", "seed_time", "2"),
				),
			},
			// ImportState testing
			{
				ResourceName:            "sonarr_indexer_broadcasthenet.test",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"api_key"},
			},
			// Delete testing automatically occurs in TestCase
		},
	})
}

func testAccIndexerBroadcastheNetResourceConfig(name string, sTime int) string {
	return fmt.Sprintf(`
	resource "sonarr_indexer_broadcasthenet" "test" {
		enable_automatic_search = false
		name = "%s"
		base_url = "https://api.broadcasthe.net/"
		api_key = "API_key"
		minimum_seeders = 1
		season_pack_seed_time = 1
		seed_time = %d
		seed_ratio = 0.5
	}`, name, sTime)
}

//nolint:paralleltest // deletes an object outside Terraform, a parallel test could otherwise take over its freed ID
func TestAccIndexerBroadcastheNetResourceDisappears(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create, then delete outside Terraform: Read removes it from state instead of failing the plan
			{
				Config: testAccIndexerBroadcastheNetResourceConfig("broadcasthenetResourceTest", 2),
				Check: testAccCheckResourceDisappears("sonarr_indexer_broadcasthenet.test", func(client *sonarr.APIClient, id int32) (*http.Response, error) {
					return client.IndexerAPI.DeleteIndexer(context.TODO(), id).Execute()
				}),
				ExpectNonEmptyPlan: true,
			},
			// Create again after the deletion outside Terraform
			{
				Config: testAccIndexerBroadcastheNetResourceConfig("broadcasthenetResourceTest", 2),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet("sonarr_indexer_broadcasthenet.test", "id"),
				),
			},
		},
	})
}

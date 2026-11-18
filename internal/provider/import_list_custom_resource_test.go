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

func TestAccImportListCustomResource(t *testing.T) {
	t.Parallel()

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Unauthorized Create
			{
				Config:      testAccImportListCustomResourceConfig("resourceCustomTest", "false") + testUnauthorizedProvider,
				ExpectError: regexp.MustCompile("Client Error"),
			},
			// Create and Read testing
			{
				PreConfig: rootFolderDSInit,
				Config:    testAccImportListCustomResourceConfig("resourceCustomTest", "false"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("sonarr_import_list_custom.test", "season_folder", "false"),
					resource.TestCheckResourceAttrSet("sonarr_import_list_custom.test", "id"),
				),
			},
			// Unauthorized Read
			{
				Config:      testAccImportListCustomResourceConfig("resourceCustomTest", "false") + testUnauthorizedProvider,
				ExpectError: regexp.MustCompile("Client Error"),
			},
			// Update and Read testing
			{
				Config: testAccImportListCustomResourceConfig("resourceCustomTest", "true"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("sonarr_import_list_custom.test", "season_folder", "true"),
				),
			},
			// ImportState testing
			{
				ResourceName:      "sonarr_import_list_custom.test",
				ImportState:       true,
				ImportStateVerify: true,
			},
			// Delete testing automatically occurs in TestCase
		},
	})
}

func testAccImportListCustomResourceConfig(name, folder string) string {
	return fmt.Sprintf(`

	resource "sonarr_import_list_custom" "test" {
		enable_automatic_add = false
		season_folder = %s
		should_monitor = "all"
		series_type = "standard"
		root_folder_path = "/config"
		quality_profile_id = 1
		name = "%s"
		base_url = "localhost"
		tags = []
	}`, folder, name)
}

//nolint:paralleltest // deletes an object outside Terraform, a parallel test could otherwise take over its freed ID
func TestAccImportListCustomResourceDisappears(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create, then delete outside Terraform: Read removes it from state instead of failing the plan
			{
				PreConfig: rootFolderDSInit,
				Config:    testAccImportListCustomResourceConfig("resourceCustomTest", "true"),
				Check: testAccCheckResourceDisappears("sonarr_import_list_custom.test", func(client *sonarr.APIClient, id int32) (*http.Response, error) {
					return client.ImportListAPI.DeleteImportList(context.TODO(), id).Execute()
				}),
				ExpectNonEmptyPlan: true,
			},
			// Create again after the deletion outside Terraform
			{
				Config: testAccImportListCustomResourceConfig("resourceCustomTest", "true"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet("sonarr_import_list_custom.test", "id"),
				),
			},
		},
	})
}

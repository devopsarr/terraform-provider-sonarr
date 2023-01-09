package provider

import (
	"context"
	"testing"

	"github.com/devopsarr/sonarr-go/sonarr"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccRootFolderDataSource(t *testing.T) {
	t.Parallel()

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Read testing
			{
				PreConfig: rootFolderDSInit,
				Config:    testAccRootFolderDataSourceConfig,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.sonarr_root_folder.test", "id"),
					resource.TestCheckResourceAttr("data.sonarr_root_folder.test", "path", "/defaults")),
			},
		},
	})
}

const testAccRootFolderDataSourceConfig = `
data "sonarr_root_folder" "test" {
	path = "/defaults"
}
`

func rootFolderDSInit() {
	// ensure a /defaults root path is configured
	client := testAccAPIClient()
	folder := sonarr.NewRootFolderResource()
	folder.SetPath("/defaults")
	_, _, _ = client.RootFolderApi.CreateRootFolder(context.TODO()).RootFolderResource(*folder).Execute()
}

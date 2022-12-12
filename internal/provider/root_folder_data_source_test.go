package provider

import (
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"golift.io/starr"
	"golift.io/starr/sonarr"
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
	client := *sonarr.New(starr.New(os.Getenv("SONARR_API_KEY"), os.Getenv("SONARR_URL"), 0))
	_, _ = client.AddRootFolder(&sonarr.RootFolder{Path: "/defaults"})
}

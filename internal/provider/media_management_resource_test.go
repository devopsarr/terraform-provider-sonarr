package provider

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccMediaManagementResource(t *testing.T) {
	t.Parallel()

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read testing
			{
				Config: testAccMediaManagementResourceConfig("none"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("sonarr_media_management.test", "file_date", "none"),
					resource.TestCheckResourceAttrSet("sonarr_media_management.test", "id"),
				),
			},
			// Update and Read testing
			{
				Config: testAccMediaManagementResourceConfig("localAirDate"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("sonarr_media_management.test", "file_date", "localAirDate"),
				),
			},
			// ImportState testing
			{
				ResourceName:      "sonarr_media_management.test",
				ImportState:       true,
				ImportStateVerify: true,
			},
			// Delete testing automatically occurs in TestCase
		},
	})
}

func testAccMediaManagementResourceConfig(date string) string {
	return fmt.Sprintf(`
	resource "sonarr_media_management" "test" {
		unmonitor_previous_episodes = true
		hardlinks_copy              = true
		create_empty_folders        = true
		delete_empty_folders        = true
		enable_media_info           = true
		import_extra_files          = true
		set_permissions             = true
		skip_free_space_check       = true
		minimum_free_space          = 100
		recycle_bin_days            = 7
		chmod_folder                = "755"
		chown_group                 = "arrs"
		download_propers_repacks    = "preferAndUpgrade"
		episode_title_required      = "always"
		extra_file_extensions       = "srt,info"
		file_date                   = "%s"
		recycle_bin_path            = "/config/MediaCover"
		rescan_after_refresh        = "always"
	}`, date)
}

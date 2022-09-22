resource "sonarr_download_client_config" "example" {
  enable_completed_download_handling = true
  auto_redownload_failed             = false
}
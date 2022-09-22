resource "sonarr_remote_path_mapping" "example" {
  host        = "%s"
  remote_path = "%s"
  local_path  = "/config/"
}

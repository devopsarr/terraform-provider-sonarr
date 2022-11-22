resource "sonarr_remote_path_mapping" "example" {
  host        = "www.transmission.com"
  remote_path = "/download/"
  local_path  = "/transmission-download/"
}
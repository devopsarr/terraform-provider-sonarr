---
# generated by https://github.com/hashicorp/terraform-plugin-docs
page_title: "sonarr_download_client_deluge Resource - terraform-provider-sonarr"
subcategory: "Download Clients"
description: |-
  Download Client Deluge resource.
  For more information refer to Download Client https://wiki.servarr.com/sonarr/settings#download-clients and Deluge https://wiki.servarr.com/sonarr/supported#deluge.
---

# sonarr_download_client_deluge (Resource)

<!-- subcategory:Download Clients -->
Download Client Deluge resource.
For more information refer to [Download Client](https://wiki.servarr.com/sonarr/settings#download-clients) and [Deluge](https://wiki.servarr.com/sonarr/supported#deluge).

## Example Usage

```terraform
resource "sonarr_download_client_deluge" "example" {
  enable   = true
  priority = 1
  name     = "Example"
  host     = "deluge"
  url_base = "/deluge/"
  port     = 9091
}
```

<!-- schema generated by tfplugindocs -->
## Schema

### Required

- `name` (String) Download Client name.
- `password` (String, Sensitive) Password.

### Optional

- `add_paused` (Boolean) Add paused flag.
- `enable` (Boolean) Enable flag.
- `host` (String) host.
- `older_tv_priority` (Number) Older TV priority. `0` Last, `1` First.
- `port` (Number) Port.
- `priority` (Number) Priority.
- `recent_tv_priority` (Number) Recent TV priority. `0` Last, `1` First.
- `remove_completed_downloads` (Boolean) Remove completed downloads flag.
- `remove_failed_downloads` (Boolean) Remove failed downloads flag.
- `tags` (Set of Number) List of associated tags.
- `tv_category` (String) TV category.
- `tv_imported_category` (String) TV imported category.
- `url_base` (String) Base URL.
- `use_ssl` (Boolean) Use SSL flag.

### Read-Only

- `id` (Number) Download Client ID.

## Import

Import is supported using the following syntax:

```shell
# import using the API/UI ID
terraform import sonarr_download_client_deluge.example 1
```

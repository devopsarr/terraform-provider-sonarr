---
# generated by https://github.com/hashicorp/terraform-plugin-docs
page_title: "sonarr_download_client_transmission Resource - terraform-provider-sonarr"
subcategory: "Download Clients"
description: |-
  Download Client Transmission resource.
  For more information refer to Download Client https://wiki.servarr.com/sonarr/settings#download-clients and Transmission https://wiki.servarr.com/sonarr/supported#transmission.
---

# sonarr_download_client_transmission (Resource)

<!-- subcategory:Download Clients -->
Download Client Transmission resource.
For more information refer to [Download Client](https://wiki.servarr.com/sonarr/settings#download-clients) and [Transmission](https://wiki.servarr.com/sonarr/supported#transmission).

## Example Usage

```terraform
resource "sonarr_download_client_transmission" "example" {
  enable   = true
  priority = 1
  name     = "Example"
  host     = "transmission"
  url_base = "/transmission/"
  port     = 9091
}
```

<!-- schema generated by tfplugindocs -->
## Schema

### Required

- `name` (String) Download Client name.

### Optional

- `add_paused` (Boolean) Add paused flag.
- `enable` (Boolean) Enable flag.
- `host` (String) host.
- `older_tv_priority` (Number) Older TV priority. `0` Last, `1` First.
- `password` (String, Sensitive) Password.
- `port` (Number) Port.
- `priority` (Number) Priority.
- `recent_tv_priority` (Number) Recent TV priority. `0` Last, `1` First.
- `remove_completed_downloads` (Boolean) Remove completed downloads flag.
- `remove_failed_downloads` (Boolean) Remove failed downloads flag.
- `tags` (Set of Number) List of associated tags.
- `tv_category` (String) TV category.
- `tv_directory` (String) TV directory.
- `url_base` (String) Base URL.
- `use_ssl` (Boolean) Use SSL flag.
- `username` (String) Username.

### Read-Only

- `id` (Number) Download Client ID.

## Import

Import is supported using the following syntax:

```shell
# import using the API/UI ID
terraform import sonarr_download_client_transmission.example 1
```

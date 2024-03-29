---
# generated by https://github.com/hashicorp/terraform-plugin-docs
page_title: "sonarr_import_list_custom Resource - terraform-provider-sonarr"
subcategory: "Import Lists"
description: |-
  ImportList Custom resource.
  For more information refer to Import List https://wiki.servarr.com/sonarr/settings#import-lists and Custom https://wiki.servarr.com/sonarr/supported#customimport.
---

# sonarr_import_list_custom (Resource)

<!-- subcategory:Import Lists -->
ImportList Custom resource.
For more information refer to [Import List](https://wiki.servarr.com/sonarr/settings#import-lists) and [Custom](https://wiki.servarr.com/sonarr/supported#customimport).

## Example Usage

```terraform
resource "sonarr_import_list_custom" "example" {
  enable_automatic_add = true
  season_folder        = true
  should_monitor       = "all"
  series_type          = "standard"
  root_folder_path     = sonarr_root_folder.example.path
  quality_profile_id   = 1
  name                 = "Example"
  base_url             = "localhost:8080"
}
```

<!-- schema generated by tfplugindocs -->
## Schema

### Required

- `base_url` (String) Base URL.
- `enable_automatic_add` (Boolean) Enable automatic add flag.
- `name` (String) Import List name.
- `quality_profile_id` (Number) Quality profile ID.
- `root_folder_path` (String) Root folder path.
- `season_folder` (Boolean) Season folder flag.
- `series_type` (String) Series type.
- `should_monitor` (String) Should monitor.

### Optional

- `tags` (Set of Number) List of associated tags.

### Read-Only

- `id` (Number) Import List ID.

## Import

Import is supported using the following syntax:

```shell
# import using the API/UI ID
terraform import sonarr_import_list_custom.example 1
```

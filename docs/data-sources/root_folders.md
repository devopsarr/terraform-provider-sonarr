---
# generated by https://github.com/hashicorp/terraform-plugin-docs
page_title: "sonarr_root_folders Data Source - terraform-provider-sonarr"
subcategory: "Media Management"
description: |-
  List all available Root Folders ../resources/root_folder.
---

# sonarr_root_folders (Data Source)

<!-- subcategory:Media Management -->
List all available [Root Folders](../resources/root_folder).

## Example Usage

```terraform
data "sonarr_root_folders" "example" {
}
```

<!-- schema generated by tfplugindocs -->
## Schema

### Read-Only

- `id` (String) The ID of this resource.
- `root_folders` (Attributes Set) Root Folder list. (see [below for nested schema](#nestedatt--root_folders))

<a id="nestedatt--root_folders"></a>
### Nested Schema for `root_folders`

Read-Only:

- `accessible` (Boolean) Access flag.
- `id` (Number) Root Folder ID.
- `path` (String) Root Folder absolute path.
- `unmapped_folders` (Attributes Set) List of folders with no associated series. (see [below for nested schema](#nestedatt--root_folders--unmapped_folders))

<a id="nestedatt--root_folders--unmapped_folders"></a>
### Nested Schema for `root_folders.unmapped_folders`

Read-Only:

- `name` (String) Name of unmapped folder.
- `path` (String) Path of unmapped folder.

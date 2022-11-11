---
# generated by https://github.com/hashicorp/terraform-plugin-docs
page_title: "sonarr_release_profile Resource - terraform-provider-sonarr"
subcategory: "Profiles"
description: |-
  Release Profile resource.
  For more information refer to Release Profiles https://wiki.servarr.com/sonarr/settings#release-profiles documentation.
---

# sonarr_release_profile (Resource)

<!-- subcategory:Profiles -->Release Profile resource.
For more information refer to [Release Profiles](https://wiki.servarr.com/sonarr/settings#release-profiles) documentation.

## Example Usage

```terraform
resource "sonarr_release_profile" "example" {
  enabled    = true
  name       = "Example"
  required   = ["proper"]
  ignored    = ["repack"]
  indexer_id = 1
}
```

<!-- schema generated by tfplugindocs -->
## Schema

### Optional

- `enabled` (Boolean) Enabled
- `ignored` (Set of String) Ignored terms. At least one of `required` and `ignored` must be set.
- `indexer_id` (Number) Indexer ID. Set `0` for all.
- `name` (String) Release profile name.
- `required` (Set of String) Required terms. At least one of `required` and `ignored` must be set.
- `tags` (Set of Number) List of associated tags.

### Read-Only

- `id` (Number) Release Profile ID.

## Import

Import is supported using the following syntax:

```shell
# import using the API/UI ID
terraform import sonarr_release_profile.example 10
```
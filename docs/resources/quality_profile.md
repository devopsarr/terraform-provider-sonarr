---
# generated by https://github.com/hashicorp/terraform-plugin-docs
page_title: "sonarr_quality_profile Resource - terraform-provider-sonarr"
subcategory: ""
description: |-
  QualityProfile resource
---

# sonarr_quality_profile (Resource)

QualityProfile resource

## Example Usage

```terraform
resource "sonarr_quality_profile" "example" {
  name            = "example-4k"
  upgrade_allowed = true
  cutoff          = 1100

  quality_groups = [
    {
      id   = 1100
      name = "4k"
      qualities = [
        {
          id         = 18
          name       = "WEBDL-2160p"
          source     = "web"
          resolution = 2160
        },
        {
          id         = 19
          name       = "Bluray-2160p"
          source     = "bluray"
          resolution = 2160
        }
      ]
    }
  ]
}
```

<!-- schema generated by tfplugindocs -->
## Schema

### Required

- `name` (String) Name
- `quality_groups` (Attributes Set) Quality groups (see [below for nested schema](#nestedatt--quality_groups))

### Optional

- `cutoff` (Number) Quality ID to which cutoff
- `upgrade_allowed` (Boolean) Upgrade allowed flag

### Read-Only

- `id` (Number) ID of qualityprofile

<a id="nestedatt--quality_groups"></a>
### Nested Schema for `quality_groups`

Required:

- `qualities` (Attributes Set) Qualities in group (see [below for nested schema](#nestedatt--quality_groups--qualities))

Optional:

- `id` (Number) ID of quality group
- `name` (String) Name of quality group

<a id="nestedatt--quality_groups--qualities"></a>
### Nested Schema for `quality_groups.qualities`

Optional:

- `id` (Number) ID of quality group
- `name` (String) Name of quality group
- `resolution` (Number) Resolution
- `source` (String) Source

## Import

Import is supported using the following syntax:

```shell
# import using the API/UI ID
terraform import sonarr_quality_profile.example 10
```

---
# generated by https://github.com/hashicorp/terraform-plugin-docs
page_title: "sonarr_quality_profile Data Source - terraform-provider-sonarr"
subcategory: "Profiles"
description: |-
  Single Quality Profile ../resources/quality_profile.
---

# sonarr_quality_profile (Data Source)

<!-- subcategory:Profiles -->
Single [Quality Profile](../resources/quality_profile).

## Example Usage

```terraform
data "sonarr_quality_profile" "example" {
  name = "HD"
}
```

<!-- schema generated by tfplugindocs -->
## Schema

### Required

- `name` (String) Quality Profile Name.

### Read-Only

- `cutoff` (Number) Quality ID to which cutoff.
- `cutoff_format_score` (Number) Cutoff format score.
- `format_items` (Attributes Set) Quality groups. (see [below for nested schema](#nestedatt--format_items))
- `id` (Number) Quality Profile ID.
- `min_format_score` (Number) Min format score.
- `min_upgrade_format_score` (Number) Min upgrade format score.
- `quality_groups` (Attributes List) Quality groups. (see [below for nested schema](#nestedatt--quality_groups))
- `upgrade_allowed` (Boolean) Upgrade allowed flag.

<a id="nestedatt--format_items"></a>
### Nested Schema for `format_items`

Read-Only:

- `format` (Number) Format.
- `name` (String) Name.
- `score` (Number) Score.


<a id="nestedatt--quality_groups"></a>
### Nested Schema for `quality_groups`

Read-Only:

- `id` (Number) Quality group ID.
- `name` (String) Quality group name.
- `qualities` (Attributes List) Qualities in group. (see [below for nested schema](#nestedatt--quality_groups--qualities))

<a id="nestedatt--quality_groups--qualities"></a>
### Nested Schema for `quality_groups.qualities`

Read-Only:

- `id` (Number) Quality ID.
- `name` (String) Quality name.
- `resolution` (Number) Resolution.
- `source` (String) Source.

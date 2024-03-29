---
# generated by https://github.com/hashicorp/terraform-plugin-docs
page_title: "sonarr_delay_profiles Data Source - terraform-provider-sonarr"
subcategory: "Profiles"
description: |-
  List all available Delay Profiles ../resources/delay_profile.
---

# sonarr_delay_profiles (Data Source)

<!-- subcategory:Profiles -->
List all available [Delay Profiles](../resources/delay_profile).

## Example Usage

```terraform
data "sonarr_delay_profiles" "example" {
}
```

<!-- schema generated by tfplugindocs -->
## Schema

### Read-Only

- `delay_profiles` (Attributes Set) Delay Profile list. (see [below for nested schema](#nestedatt--delay_profiles))
- `id` (String) The ID of this resource.

<a id="nestedatt--delay_profiles"></a>
### Nested Schema for `delay_profiles`

Read-Only:

- `bypass_if_above_custom_format_score` (Boolean) Bypass for higher custom format score flag.
- `bypass_if_highest_quality` (Boolean) Bypass for highest quality Flag.
- `enable_torrent` (Boolean) Torrent allowed Flag.
- `enable_usenet` (Boolean) Usenet allowed Flag.
- `id` (Number) Delay Profile ID.
- `minimum_custom_format_score` (Number) Minimum custom format score.
- `order` (Number) Order.
- `preferred_protocol` (String) Preferred protocol.
- `tags` (Set of Number) List of associated tags.
- `torrent_delay` (Number) Torrent Delay.
- `usenet_delay` (Number) Usenet delay.

---
# generated by https://github.com/hashicorp/terraform-plugin-docs
page_title: "sonarr_metadata Data Source - terraform-provider-sonarr"
subcategory: "Metadata"
description: |-
  Single Metadata ../resources/metadata.
---

# sonarr_metadata (Data Source)

<!-- subcategory:Metadata -->Single [Metadata](../resources/metadata).

## Example Usage

```terraform
data "sonarr_metadata" "example" {
  name = "Example"
}
```

<!-- schema generated by tfplugindocs -->
## Schema

### Required

- `name` (String) Metadata name.

### Optional

- `episode_images` (Boolean) Episode images flag.

### Read-Only

- `config_contract` (String) Metadata configuration template.
- `enable` (Boolean) Enable flag.
- `episode_metadata` (Boolean) Episode metadata flag.
- `id` (Number) Metadata ID.
- `implementation` (String) Metadata implementation name.
- `season_images` (Boolean) Season images flag.
- `series_images` (Boolean) Series images flag.
- `series_metadata` (Boolean) Series metafata flag.
- `series_metadata_url` (Boolean) Series metadata URL flag.
- `tags` (Set of Number) List of associated tags.


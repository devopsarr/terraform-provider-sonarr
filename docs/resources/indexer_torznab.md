---
# generated by https://github.com/hashicorp/terraform-plugin-docs
page_title: "sonarr_indexer_torznab Resource - terraform-provider-sonarr"
subcategory: "Indexers"
description: |-
  Indexer Torznab resource.
  For more information refer to Indexer https://wiki.servarr.com/sonarr/settings#indexers and Torznab https://wiki.servarr.com/sonarr/supported#torznab.
---

# sonarr_indexer_torznab (Resource)

<!-- subcategory:Indexers -->
Indexer Torznab resource.
For more information refer to [Indexer](https://wiki.servarr.com/sonarr/settings#indexers) and [Torznab](https://wiki.servarr.com/sonarr/supported#torznab).

## Example Usage

```terraform
resource "sonarr_indexer_torznab" "example" {
  enable_automatic_search = true
  name                    = "Example"
  base_url                = "https://feed.animetosho.org"
  api_path                = "/nabapi"
  anime_categories        = [5070]
  minimum_seeders         = 1
}
```

<!-- schema generated by tfplugindocs -->
## Schema

### Required

- `base_url` (String) Base URL.
- `name` (String) IndexerTorznab name.

### Optional

- `additional_parameters` (String) Additional parameters.
- `anime_categories` (Set of Number) Anime categories list.
- `anime_standard_format_search` (Boolean) Search anime in standard format.
- `api_key` (String, Sensitive) API key.
- `api_path` (String) API path.
- `categories` (Set of Number) Categories list.
- `download_client_id` (Number) Download client ID.
- `enable_automatic_search` (Boolean) Enable automatic search flag.
- `enable_interactive_search` (Boolean) Enable interactive search flag.
- `enable_rss` (Boolean) Enable RSS flag.
- `minimum_seeders` (Number) Minimum seeders.
- `priority` (Number) Priority.
- `season_pack_seed_time` (Number) Season seed time.
- `seed_ratio` (Number) Seed ratio.
- `seed_time` (Number) Seed time.
- `tags` (Set of Number) List of associated tags.

### Read-Only

- `id` (Number) IndexerTorznab ID.

## Import

Import is supported using the following syntax:

```shell
# import using the API/UI ID
terraform import sonarr_indexer_torznab.example 1
```

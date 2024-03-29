---
# generated by https://github.com/hashicorp/terraform-plugin-docs
page_title: "sonarr_indexer_hdbits Resource - terraform-provider-sonarr"
subcategory: "Indexers"
description: |-
  Indexer HDBits resource.
  For more information refer to Indexer https://wiki.servarr.com/sonarr/settings#indexers and HDBits https://wiki.servarr.com/sonarr/supported#hdbits.
---

# sonarr_indexer_hdbits (Resource)

<!-- subcategory:Indexers -->
Indexer HDBits resource.
For more information refer to [Indexer](https://wiki.servarr.com/sonarr/settings#indexers) and [HDBits](https://wiki.servarr.com/sonarr/supported#hdbits).

## Example Usage

```terraform
resource "sonarr_indexer_hdbits" "example" {
  enable_automatic_search = true
  name                    = "Example"
  base_url                = "https://hdbits.org"
  username                = "User"
  api_key                 = "APIKey"
  minimum_seeders         = 1
}
```

<!-- schema generated by tfplugindocs -->
## Schema

### Required

- `api_key` (String, Sensitive) API key.
- `name` (String) IndexerHdbits name.
- `username` (String) Username.

### Optional

- `base_url` (String) Base URL.
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

- `id` (Number) IndexerHdbits ID.

## Import

Import is supported using the following syntax:

```shell
# import using the API/UI ID
terraform import sonarr_indexer_hdbits.example 1
```

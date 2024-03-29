---
# generated by https://github.com/hashicorp/terraform-plugin-docs
page_title: "sonarr_download_clients Data Source - terraform-provider-sonarr"
subcategory: "Download Clients"
description: |-
  List all available Download Clients ../resources/download_client.
---

# sonarr_download_clients (Data Source)

<!-- subcategory:Download Clients -->
List all available [Download Clients](../resources/download_client).

## Example Usage

```terraform
data "sonarr_download_clients" "example" {
}
```

<!-- schema generated by tfplugindocs -->
## Schema

### Read-Only

- `download_clients` (Attributes Set) Download Client list. (see [below for nested schema](#nestedatt--download_clients))
- `id` (String) The ID of this resource.

<a id="nestedatt--download_clients"></a>
### Nested Schema for `download_clients`

Read-Only:

- `add_paused` (Boolean) Add paused flag.
- `add_stopped` (Boolean) Add stopped flag.
- `additional_tags` (Set of Number) Additional tags, `0` TitleSlug, `1` Quality, `2` Language, `3` ReleaseGroup, `4` Year, `5` Indexer, `6` Network.
- `api_key` (String, Sensitive) API key.
- `category` (String) Category.
- `config_contract` (String) DownloadClient configuration template.
- `destination` (String) Destination.
- `enable` (Boolean) Enable flag.
- `field_tags` (Set of String) Field tags.
- `first_and_last` (Boolean) First and last flag.
- `host` (String) host.
- `id` (Number) Download Client ID.
- `implementation` (String) DownloadClient implementation name.
- `initial_state` (Number) Initial state. `0` Start, `1` ForceStart, `2` Pause.
- `intial_state` (Number) Initial state, with Stop support. `0` Start, `1` ForceStart, `2` Pause, `3` Stop.
- `magnet_file_extension` (String) Magnet file extension.
- `name` (String) Download Client name.
- `nzb_folder` (String) NZB folder.
- `older_tv_priority` (Number) Older TV priority. `0` Last, `1` First.
- `password` (String, Sensitive) Password.
- `port` (Number) Port.
- `post_import_tags` (Set of String) Post import tags.
- `priority` (Number) Priority.
- `protocol` (String) Protocol. Valid values are 'usenet' and 'torrent'.
- `read_only` (Boolean) Read only flag.
- `recent_tv_priority` (Number) Recent TV priority. `0` Last, `1` First.
- `remove_completed_downloads` (Boolean) Remove completed downloads flag.
- `remove_failed_downloads` (Boolean) Remove failed downloads flag.
- `rpc_path` (String) RPC path.
- `save_magnet_files` (Boolean) Save magnet files flag.
- `secret_token` (String, Sensitive) Secret token.
- `sequential_order` (Boolean) Sequential order flag.
- `start_on_add` (Boolean) Start on add flag.
- `strm_folder` (String) STRM folder.
- `tags` (Set of Number) List of associated tags.
- `torrent_folder` (String) Torrent folder.
- `tv_category` (String) TV category.
- `tv_directory` (String) TV directory.
- `tv_imported_category` (String) TV imported category.
- `url_base` (String) Base URL.
- `use_ssl` (Boolean) Use SSL flag.
- `username` (String) Username.
- `watch_folder` (String) Watch folder flag.

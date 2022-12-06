---
# generated by https://github.com/hashicorp/terraform-plugin-docs
page_title: "sonarr_notification_plex Resource - terraform-provider-sonarr"
subcategory: "Notifications"
description: |-
  Notification Plex resource.
  For more information refer to Notification https://wiki.servarr.com/sonarr/settings#connect and Plex https://wiki.servarr.com/sonarr/supported#plexserver.
---

# sonarr_notification_plex (Resource)

<!-- subcategory:Notifications -->Notification Plex resource.
For more information refer to [Notification](https://wiki.servarr.com/sonarr/settings#connect) and [Plex](https://wiki.servarr.com/sonarr/supported#plexserver).

## Example Usage

```terraform
resource "sonarr_notification_plex" "example" {
  on_download                        = true
  on_upgrade                         = true
  on_rename                          = false
  on_series_delete                   = false
  on_episode_file_delete             = false
  on_episode_file_delete_for_upgrade = true

  include_health_warnings = false
  name                    = "Example"

  host       = "plex.lcl"
  port       = 32400
  auth_token = "AuthTOKEN"
}
```

<!-- schema generated by tfplugindocs -->
## Schema

### Required

- `auth_token` (String, Sensitive) Auth Token.
- `host` (String) Host.
- `include_health_warnings` (Boolean) Include health warnings.
- `name` (String) NotificationPlex name.
- `on_download` (Boolean) On download flag.
- `on_episode_file_delete` (Boolean) On episode file delete flag.
- `on_episode_file_delete_for_upgrade` (Boolean) On episode file delete for upgrade flag.
- `on_rename` (Boolean) On rename flag.
- `on_series_delete` (Boolean) On series delete flag.
- `on_upgrade` (Boolean) On upgrade flag.

### Optional

- `port` (Number) Port.
- `tags` (Set of Number) List of associated tags.
- `update_library` (Boolean) Update library flag.
- `use_ssl` (Boolean) Use SSL flag.

### Read-Only

- `id` (Number) Notification ID.

## Import

Import is supported using the following syntax:

```shell
# import using the API/UI ID
terraform import sonarr_notification_plex.example 1
```
---
# generated by https://github.com/hashicorp/terraform-plugin-docs
page_title: "sonarr_auto_tag_condition_genres Data Source - terraform-provider-sonarr"
subcategory: "Tags"
description: |-
  Auto Tag Condition Genres data source.
  For more intagion refer to Auto Tag Conditions https://wiki.servarr.com/sonarr/settings#conditions.
---

# sonarr_auto_tag_condition_genres (Data Source)

<!-- subcategory:Tags -->
 Auto Tag Condition Genres data source.
For more intagion refer to [Auto Tag Conditions](https://wiki.servarr.com/sonarr/settings#conditions).



<!-- schema generated by tfplugindocs -->
## Schema

### Required

- `name` (String) Specification name.
- `negate` (Boolean) Negate flag.
- `required` (Boolean) Computed flag.
- `value` (String) Genres. Space separated list of genres.

### Read-Only

- `id` (Number) Auto tag condition series type ID.
- `implementation` (String) Implementation.

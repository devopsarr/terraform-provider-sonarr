---
# generated by https://github.com/hashicorp/terraform-plugin-docs
page_title: "sonarr_tags Data Source - terraform-provider-sonarr"
subcategory: "Tags"
description: |-
  List all available Tags ../resources/tag.
---

# sonarr_tags (Data Source)

<!-- subcategory:Tags -->
List all available [Tags](../resources/tag).

## Example Usage

```terraform
data "sonarr_tags" "example" {
}
```

<!-- schema generated by tfplugindocs -->
## Schema

### Read-Only

- `id` (String) The ID of this resource.
- `tags` (Attributes Set) Tag list. (see [below for nested schema](#nestedatt--tags))

<a id="nestedatt--tags"></a>
### Nested Schema for `tags`

Read-Only:

- `id` (Number) Tag ID.
- `label` (String) Tag label.

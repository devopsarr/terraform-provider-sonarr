---
# generated by https://github.com/hashicorp/terraform-plugin-docs
page_title: "sonarr_tags Data Source - terraform-provider-sonarr"
subcategory: ""
description: |-
  List all available tags
---

# sonarr_tags (Data Source)

List all available tags

## Example Usage

```terraform
data "sonarr_tags" "example" {
}
```

<!-- schema generated by tfplugindocs -->
## Schema

### Read-Only

- `id` (String) The ID of this resource.
- `tags` (Attributes Set) List of tags (see [below for nested schema](#nestedatt--tags))

<a id="nestedatt--tags"></a>
### Nested Schema for `tags`

Read-Only:

- `id` (Number) ID of tag
- `label` (String) Actual tag



---
# generated by https://github.com/hashicorp/terraform-plugin-docs
page_title: "sonarr_tag Data Source - terraform-provider-sonarr"
subcategory: "Tags"
description: |-
  Single Tag ../resources/tag.
---

# sonarr_tag (Data Source)

<!-- subcategory:Tags -->
Single [Tag](../resources/tag).

## Example Usage

```terraform
data "sonarr_tag" "example" {
  label = "example"
}
```

<!-- schema generated by tfplugindocs -->
## Schema

### Required

- `label` (String) Tag label.

### Read-Only

- `id` (Number) Tag ID.

---
# generated by https://github.com/hashicorp/terraform-plugin-docs
page_title: "sonarr_languages Data Source - terraform-provider-sonarr"
subcategory: "Languages"
description: |-
  List all available Languages ../data-sources/language.
---

# sonarr_languages (Data Source)

<!-- subcategory:Languages -->
List all available [Languages](../data-sources/language).

## Example Usage

```terraform
data "sonarr_languages" "example" {
}
```

<!-- schema generated by tfplugindocs -->
## Schema

### Read-Only

- `id` (String) The ID of this resource.
- `languages` (Attributes Set) Language list. (see [below for nested schema](#nestedatt--languages))

<a id="nestedatt--languages"></a>
### Nested Schema for `languages`

Read-Only:

- `id` (Number) Language ID.
- `name` (String) Language.
- `name_lower` (String) Language in lowercase.

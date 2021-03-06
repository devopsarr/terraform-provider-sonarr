---
# generated by https://github.com/hashicorp/terraform-plugin-docs
page_title: "sonarr_language_profiles Data Source - terraform-provider-sonarr"
subcategory: ""
description: |-
  List all available languageprofiles
---

# sonarr_language_profiles (Data Source)

List all available languageprofiles

## Example Usage

```terraform
data "sonarr_language_profiles" "example" {
}
```

<!-- schema generated by tfplugindocs -->
## Schema

### Read-Only

- `id` (String) The ID of this resource.
- `language_profiles` (Attributes Set) List of languageprofiles (see [below for nested schema](#nestedatt--language_profiles))

<a id="nestedatt--language_profiles"></a>
### Nested Schema for `language_profiles`

Read-Only:

- `cutoff_language` (String) Cutoff Language
- `id` (Number) ID of languageprofile
- `languages` (Set of String) list of languages in profile
- `name` (String) Name of languageprofile
- `upgrade_allowed` (Boolean) Upgrade allowed Flag



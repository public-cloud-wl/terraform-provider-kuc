---
# generated by https://github.com/hashicorp/terraform-plugin-docs
page_title: "kuc_user Resource - kuc"
subcategory: ""
description: |-
  User resource
---

# kuc_user (Resource)

User resource

## Example Usage

```terraform
resource "kuc_user" "example" {
  username = "a123456"
}
```

<!-- schema generated by tfplugindocs -->
## Schema

### Required

- `username` (String) User name

### Read-Only

- `id` (String) User identifier
---
# generated by https://github.com/hashicorp/terraform-plugin-docs
page_title: "snyk_organization Data Source - terraform-provider-snyk"
subcategory: ""
description: |-
  
---

# snyk_organization (Data Source)



## Example Usage

```terraform
data "snyk_organization" "example" {
  id = "ORG_ID_HERE"
}
```

<!-- schema generated by tfplugindocs -->
## Schema

### Read-Only

- `id` (String) The ID of this resource.
- `name` (String) Organization name
- `slug` (String) Organization short name
- `url` (String) Organization web console login



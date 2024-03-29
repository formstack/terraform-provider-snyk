---
# generated by https://github.com/hashicorp/terraform-plugin-docs
page_title: "snyk_integration Resource - terraform-provider-snyk"
subcategory: ""
description: |-
  
---

# snyk_integration (Resource)



## Example Usage

```terraform
resource "snyk_integration" "example_integration" {
  organization = snyk_organization.example.id
  type         = "bitbucket-cloud"
  credentials {
    username = "username"
    password = "password" # Make sure your backend is encrypted - this is stored in plaintext!
  }
}
```

<!-- schema generated by tfplugindocs -->
## Schema

### Required

- `credentials` (Block List, Min: 1, Max: 1) (see [below for nested schema](#nestedblock--credentials))
- `organization` (String)
- `type` (String)

### Read-Only

- `id` (String) The ID of this resource.

<a id="nestedblock--credentials"></a>
### Nested Schema for `credentials`

Optional:

- `password` (String, Sensitive)
- `region` (String)
- `registry_base` (String)
- `role_arn` (String)
- `token` (String, Sensitive)
- `url` (String)
- `username` (String)



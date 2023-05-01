# Building and using locally
- run `go install` from the root of the repository
- there should now be a binary at `$GOPATH/bin`
- put the binary in `/usr/local/share/terraform/plugins/registry.terraform.io/formstack/snyk/<version number>/` (you'll need to create the folders)
- in your `.terraformrc` file put in a config that looks like this
```
provider_installation {
  filesystem_mirror {
    path    = "/usr/local/share/terraform/plugins"
    include = ["formstack/snyk"]
  }

direct {
  exclude = ["formstack/snyk"]
}
```

- in your terraform project reference the provider like this (you can also refernce a specific version as you normally would if you have several installed locally)
```
terraform {
  required_providers {
    snyk = {
      source = "formstack/snyk"
    }
}
```

# Building with your local changes
The terraform provider files import the API from github directly at `github.com/formstack/terraform-provider-snyk/snyk/api` so if you are adding something to the API your local API changes would need to be merged in before your local terraform files could see and use them.

In `go.mod` add a line at the bottom similar to the following to direct go to build using your local files
```
replace github.com/formstack/terraform-provider-snyk => /Users/jay/formstack/terraform-provider-snyk
```

- I'm currently using Go 1.17 and everything seems to be working fine



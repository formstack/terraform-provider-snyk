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
- I'm currently using Go 1.1.7 and everything seems to be working fine


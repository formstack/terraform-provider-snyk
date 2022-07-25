resource "snyk_project" "example" {
  organization     = data.snyk_organization.example.id
  repository_owner = "repo-owner"
  repository_name  = "repo-name"
  branch           = "main"
  integration      = data.snyk_integration.example.id
}

package snyk

import (
	"context"

	"github.com/formstack/terraform-provider-snyk/snyk/api"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func init() {
	// Set descriptions to support markdown syntax, this will be used in document generation
	// and the language server.
	schema.DescriptionKind = schema.StringMarkdown
}

func Provider(version string) func() *schema.Provider {
	return func() *schema.Provider {
		p := &schema.Provider{
			Schema: map[string]*schema.Schema{
				"group_id": {
					Type:        schema.TypeString,
					Required:    true,
					DefaultFunc: schema.EnvDefaultFunc("SNYK_API_GROUP", nil),
				},
				"api_key": {
					Type:        schema.TypeString,
					Required:    true,
					Sensitive:   true,
					DefaultFunc: schema.EnvDefaultFunc("SNYK_API_KEY", nil),
				},
			},
			ResourcesMap: map[string]*schema.Resource{
				"snyk_organization": resourceOrganization(),
				"snyk_integration":  resourceIntegration(),
				"snyk_project":      resourceProject(),
				"snyk_target":       resourceProject(),
			},
			DataSourcesMap: map[string]*schema.Resource{
				"snyk_organization":  dataSourceOrganization(),
				"snyk_project":       dataSourceProject(),
				"snyk_project_owner": dataSourceProjectOwner(),
				"snyk_integration":   dataSourceIntegration(),
			},
		}

		p.ConfigureContextFunc = configure(version, p)

		return p
	}
}

func configure(version string, p *schema.Provider) func(context.Context, *schema.ResourceData) (interface{}, diag.Diagnostics) {
	return func(ctx context.Context, d *schema.ResourceData) (interface{}, diag.Diagnostics) {
		var diags diag.Diagnostics

		config := api.SnykOptions{
			GroupId:   d.Get("group_id").(string),
			ApiKey:    d.Get("api_key").(string),
			UserAgent: p.UserAgent("terraform-provider-snyk", version),
		}

		return config, diags
	}
}

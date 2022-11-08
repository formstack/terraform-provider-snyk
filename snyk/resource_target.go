package snyk

import (
	"context"

	"github.com/formstack/terraform-provider-snyk/snyk/api"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceTarget() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceTargetCreate,
		Schema: map[string]*schema.Schema{
			"organization": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"owner": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"project_name": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"branch": {
				Type:     schema.TypeString,
				Required: true,
			},
			"integration": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"repository_owner": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"repository_name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
		},
	}
}

func resourceTargetCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	var err error

	so := m.(api.SnykOptions)

	organization := d.Get("organization").(string)
	repository_owner := d.Get("repository_owner").(string)
	branch := d.Get("branch").(string)
	repository_name := d.Get("repository_name").(string)
	integration := d.Get("integration").(string)

	if err != nil {
		return diag.FromErr(err)
	}

	var project *api.Target

	project, err = api.ImportProject(so, organization, integration, repository_owner, repository_name, branch)

	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(project.Id)
	d.Set("organization", organization)
	d.Set("project_name", project.Name)
	d.Set("branch", project.Branch)

	//setCredentialState(project.TargetCredentials, d)

	return diags
}

package snyk

import (
	"context"
	"log"

	"github.com/formstack/terraform-provider-snyk/snyk/api"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceProject() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceProjectCreate,
		ReadContext:   resourceProjectRead,
		UpdateContext: resourceProjectUpdate,
		DeleteContext: resourceProjectDelete,
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

func resourceProjectCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
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

	project, err = api.ImportTarget(so, organization, integration, repository_owner, repository_name, branch)

	if err != nil {
		return diag.FromErr(err)
	}
	log.Println("final_p_name ", project.Name)
	d.SetId(project.Id)
	d.Set("organization", organization)
	d.Set("project_name", project.Name)
	d.Set("branch", project.Branch)

	//setCredentialState(project.TargetCredentials, d)

	return diags
}

func resourceProjectRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	so := m.(api.SnykOptions)

	id := d.Id()
	orgId := d.Get("organization").(string)

	project, err := api.GetProjectById(so, orgId, id)

	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(project.Id)
	d.Set("organization", orgId)
	d.Set("branch", project.Branch)

	return diags
}

func resourceProjectUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	so := m.(api.SnykOptions)

	orgId := d.Get("organization").(string)
	intType := d.Get("type").(string)
	credentials, err := getCredentialState(d)

	if err != nil {
		return diag.FromErr(err)
	}

	integration, err := api.UpdateIntegration(so, orgId, intType, credentials)

	if err != nil {
		return diag.FromErr(err)
	}

	setCredentialState(integration.Credentials, d)

	return diags
}

func resourceProjectDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics

	so := m.(api.SnykOptions)

	orgId := d.Get("organization").(string)
	intType := d.Id()

	err := api.DeleteProject(so, orgId, intType)

	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId("")

	return diags
}

package snyk

import (
	"context"

	"github.com/formstack/terraform-provider-snyk/snyk/api"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceProject() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceProjectRead,
		Schema: map[string]*schema.Schema{
			"id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"name": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Project name",
			},
			"origin": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Project origin",
			},
			"branch": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Branch to monitor",
			},
			"owner": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The owner of the project",
			},
		},
	}
}

func dataSourceProjectRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {

	var diags diag.Diagnostics
	so := m.(api.SnykOptions)
	id := d.Get("id").(string)
	orgId := "b652c377-b52f-46de-b7ee-4607868729bd"

	project, err := api.GetProjectById(so, orgId, id)

	if err != nil {
		return diag.FromErr(err)
	}

	d.Set("id", project.Id)
	d.Set("name", project.Name)
	d.Set("origin", project.Origin)
	d.Set("branch", project.Branch)
	d.Set("owner", project.Owner)

	d.SetId(project.Id)

	return diags
}

package snyk

import (
	"context"

	"github.com/formstack/terraform-provider-snyk/snyk/api"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceProjectOwner() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceProjectOwnerRead,
		Schema: map[string]*schema.Schema{
			"id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Project ID",
			},
			"name": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Project name",
			},
			"orgid": {
				Type:        schema.TypeString,
				Description: "ID of the Organization",
				Required:    true,
			},
		},
	}
}

func dataSourceProjectOwnerRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {

	var diags diag.Diagnostics
	so := m.(api.SnykOptions)
	id := d.Get("id").(string)
	orgId := d.Get("orgid").(string)

	project, err := api.GetProjectOwner(so, orgId, id)

	if err != nil {
		return diag.FromErr(err)
	}

	d.Set("id", project.Id)
	d.Set("name", project.Name)

	d.SetId(project.Id)

	return diags
}

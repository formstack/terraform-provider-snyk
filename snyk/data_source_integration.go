package snyk

import (
	"context"

	"github.com/formstack/terraform-provider-snyk/snyk/api"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceIntegration() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceIntegrationRead,
		Schema: map[string]*schema.Schema{
			"id": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Integration ID",
			},
			"type": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Integration type",
			},
			"orgid": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Organization the integtion belongs to.",
			},
		},
	}
}

func dataSourceIntegrationRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	so := m.(api.SnykOptions)

	id := d.Get("orgid").(string)
	intType := d.Get("type").(string)

	integration, err := api.GetIntegration(so, id, intType)

	if err != nil {
		return diag.FromErr(err)
	}
	d.SetId(integration.Id)
	d.Set("organization", integration.OrgId)
	d.Set("type", integration.Type)

	return diags
}

package gotemplate

import (
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/terraform"
)

func Provider() terraform.ResourceProvider {
	return &schema.Provider{
		DataSourcesMap: map[string]*schema.Resource{
			"gotemplate": dataSourceFile(),
		},
	}
}

func dataSourceFile() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceFileRead,

		Schema: map[string]*schema.Schema{
			"template": &schema.Schema{
				Type:        schema.TypeString,
				Optional:    true,
				Description: "path to go template file",
			},
			"data": &schema.Schema{
				Type:         schema.TypeString,
				Optional:     true,
				Default:      "",
				Description:  "variables to substitute",
				ValidateFunc: nil,
			},
			"rendered": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Description: "rendered template",
			},
		},
	}
}

package kubernetes

import (
	"github.com/hashicorp/terraform/helper/schema"
)

func resourceKubectlServerVersion() *schema.Resource {
	return &schema.Resource{
		Create: dataSourceKubectlServerVersionRead,
		Read:   dataSourceKubectlServerVersionRead,
		Delete: resourceKubectlServerVersionDelete,
		Schema: map[string]*schema.Schema{
			"version": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
			"major": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
			"minor": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
			"patch": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
			"git_version": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
			"git_commit": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
			"build_date": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
			"platform": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func resourceKubectlServerVersionDelete(d *schema.ResourceData, meta interface{}) error {
	d.SetId("")
	return nil
}

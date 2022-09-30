package kubernetes

import (
	"context"
	"crypto/sha256"
	"fmt"

	"github.com/gavinbunney/terraform-provider-kubectl/yaml"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func reourceKubectlFileDocuments() *schema.Resource {
	return &schema.Resource{
		Description: "The resource to parse a multi-document yaml file into a map of manifests.",

		CreateContext: resourceKubectlFileDocumentsCreate,
		ReadContext:   resourceKubectlFileDocumentsRead,
		DeleteContext: resourceKubectlFileDocumentsDelete,

		Schema: map[string]*schema.Schema{
			"content": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"documents": {
				Type:     schema.TypeList,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Computed: true,
			},
			"manifests": {
				Type: schema.TypeMap,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Computed: true,
			},
		},
	}
}

func resourceKubectlFileDocumentsCreate(_ context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	var diags diag.Diagnostics

	content := d.Get("content").(string)
	documents, err := yaml.SplitMultiDocumentYAML(content)
	if err != nil {
		return append(diags, diag.FromErr(err)...)
	}

	manifests := make(map[string]string, 0)
	for _, doc := range documents {
		manifest, err := yaml.ParseYAML(doc)
		if err != nil {
			return append(diags, diag.Errorf("failed to parse yaml as a kubernetes yaml manifest: %v", err)...)
		}

		parsed, err := manifest.AsYAML()
		if err != nil {
			return append(diags, diag.Errorf("failed to parse convert manifest to yaml: %v", err)...)
		}

		if _, exists := manifests[manifest.GetSelfLink()]; exists {
			return append(diags, diag.Errorf("duplicate manifest found with id: %v", manifest.GetSelfLink())...)
		}

		manifests[manifest.GetSelfLink()] = parsed
	}

	if err := d.Set("documents", documents); err != nil {
		return append(diags, diag.Errorf("error setting documents: %s", err)...)
	}

	if err := d.Set("manifests", manifests); err != nil {
		return append(diags, diag.Errorf("error setting manifests: %s", err)...)
	}

	d.SetId(fmt.Sprintf("%x", sha256.Sum256([]byte(content))))

	return nil
}

func resourceKubectlFileDocumentsRead(_ context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	return nil
}

func resourceKubectlFileDocumentsDelete(_ context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	d.SetId("")
	return nil
}

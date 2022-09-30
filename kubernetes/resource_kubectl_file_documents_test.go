package kubernetes

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccKubectlResourceFileDocuments_single(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:  func() {},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccKubernetesResourceFileDocumentsConfig_basic(1),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("kubectl_file_documents.test", "documents.#", "1"),
					resource.TestCheckResourceAttr("kubectl_file_documents.test", "documents.0", "kind: Service1"),
					resource.TestCheckResourceAttr("kubectl_file_documents.test", "manifests.%", "1"),
					resource.TestCheckResourceAttr("kubectl_file_documents.test", "manifests./apis/service1s", "kind: Service1\n"),
				),
			},
		},
	})
}

func TestAccKubectlResourceFileDocuments_basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:  func() {},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccKubernetesResourceFileDocumentsConfig_basic(2),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("kubectl_file_documents.test", "documents.#", "2"),
					resource.TestCheckResourceAttr("kubectl_file_documents.test", "documents.0", "kind: Service1"),
					resource.TestCheckResourceAttr("kubectl_file_documents.test", "documents.1", "kind: Service2"),
					resource.TestCheckResourceAttr("kubectl_file_documents.test", "manifests.%", "2"),
					resource.TestCheckResourceAttr("kubectl_file_documents.test", "manifests./apis/service1s", "kind: Service1\n"),
					resource.TestCheckResourceAttr("kubectl_file_documents.test", "manifests./apis/service2s", "kind: Service2\n"),
				),
			},
		},
	})
}

func TestAccKubectlResourceFileDocuments_basicMultipleEmpty(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:  func() {},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: `
resource "kubectl_file_documents" "test" {
	content = <<YAML
kind: Service1
---
# just a comment
---
kind: Service2
---
YAML
}
`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("kubectl_file_documents.test", "documents.#", "2"),
					resource.TestCheckResourceAttr("kubectl_file_documents.test", "documents.0", "kind: Service1"),
					resource.TestCheckResourceAttr("kubectl_file_documents.test", "documents.1", "kind: Service2"),
					resource.TestCheckResourceAttr("kubectl_file_documents.test", "manifests.%", "2"),
					resource.TestCheckResourceAttr("kubectl_file_documents.test", "manifests./apis/service1s", "kind: Service1\n"),
					resource.TestCheckResourceAttr("kubectl_file_documents.test", "manifests./apis/service2s", "kind: Service2\n"),
				),
			},
		},
	})
}

func testAccKubernetesResourceFileDocumentsConfig_basic(docs int) string {
	var content = ""
	for i := 1; i <= docs; i++ {
		content += fmt.Sprintf("\nkind: Service%v\n---", i)
	}

	return fmt.Sprintf(`
resource "kubectl_file_documents" "test" {
	content = <<YAML
%s
YAML
}
`, content)
}

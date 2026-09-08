package rules

import (
	"github.com/hashicorp/hcl/v2"
	"github.com/terraform-linters/tflint-plugin-sdk/hclext"
	"github.com/terraform-linters/tflint-plugin-sdk/tflint"
)

func s3BackendBlocks(runner tflint.Runner) (hclext.Blocks, error) {
	content, err := runner.GetModuleContent(&hclext.BodySchema{
		Blocks: []hclext.BlockSchema{
			{
				Type: "terraform",
				Body: &hclext.BodySchema{
					Blocks: []hclext.BlockSchema{
						{
							Type:       "backend",
							LabelNames: []string{"type"},
							Body: &hclext.BodySchema{
								Attributes: []hclext.AttributeSchema{
									{Name: "bucket", Required: false},
									{Name: "use_lockfile", Required: false},
								},
							},
						},
					},
				},
			},
		},
	}, nil)
	if err != nil {
		return nil, err
	}

	var blocks hclext.Blocks
	for _, terraform := range content.Blocks {
		for _, backend := range terraform.Body.Blocks {
			if len(backend.Labels) > 0 && backend.Labels[0] == "s3" {
				blocks = append(blocks, backend)
			}
		}
	}
	return blocks, nil
}

func s3BackendBlockIsEmpty(runner tflint.Runner, target *hclext.Block) (bool, error) {
	file, err := runner.GetFile(target.DefRange.Filename)
	if err != nil {
		return false, err
	}

	terraformContent, _, _ := file.Body.PartialContent(&hcl.BodySchema{
		Blocks: []hcl.BlockHeaderSchema{{Type: "terraform"}},
	})
	for _, terraform := range terraformContent.Blocks {
		backendContent, _, _ := terraform.Body.PartialContent(&hcl.BodySchema{
			Blocks: []hcl.BlockHeaderSchema{{
				Type:       "backend",
				LabelNames: []string{"type"},
			}},
		})
		for _, backend := range backendContent.Blocks {
			if backend.DefRange != target.DefRange {
				continue
			}

			attributes, diagnostics := backend.Body.JustAttributes()
			return len(attributes) == 0 && !diagnostics.HasErrors(), nil
		}
	}

	return false, nil
}

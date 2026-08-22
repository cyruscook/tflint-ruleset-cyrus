package rules

import (
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

package rules

import (
	"github.com/terraform-linters/tflint-plugin-sdk/hclext"
	"github.com/terraform-linters/tflint-plugin-sdk/tflint"
)

const terraformVariableNullableMessage = "Variables must explicitly set nullable"

type TerraformVariableNullableRule struct {
	tflint.DefaultRule
}

func NewTerraformVariableNullableRule() *TerraformVariableNullableRule {
	return &TerraformVariableNullableRule{}
}

func (r *TerraformVariableNullableRule) Name() string {
	return "terraform_variable_nullable"
}

func (r *TerraformVariableNullableRule) Enabled() bool {
	return true
}

func (r *TerraformVariableNullableRule) Severity() tflint.Severity {
	return tflint.ERROR
}

func (r *TerraformVariableNullableRule) Check(runner tflint.Runner) error {
	content, err := runner.GetModuleContent(&hclext.BodySchema{
		Blocks: []hclext.BlockSchema{
			{
				Type:       "variable",
				LabelNames: []string{"name"},
				Body: &hclext.BodySchema{
					Attributes: []hclext.AttributeSchema{
						{Name: "nullable", Required: false},
					},
				},
			},
		},
	}, nil)
	if err != nil {
		return err
	}

	for _, variable := range content.Blocks {
		if _, ok := variable.Body.Attributes["nullable"]; ok {
			continue
		}
		if err := runner.EmitIssue(r, terraformVariableNullableMessage, variable.DefRange); err != nil {
			return err
		}
	}
	return nil
}

var _ tflint.Rule = (*TerraformVariableNullableRule)(nil)

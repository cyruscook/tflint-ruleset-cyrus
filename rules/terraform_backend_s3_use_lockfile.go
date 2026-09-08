package rules

import (
	"github.com/terraform-linters/tflint-plugin-sdk/tflint"
)

const terraformBackendS3UseLockfileMessage = "S3 backends must set use_lockfile to true"

type TerraformBackendS3UseLockfileRule struct {
	tflint.DefaultRule
}

func NewTerraformBackendS3UseLockfileRule() *TerraformBackendS3UseLockfileRule {
	return &TerraformBackendS3UseLockfileRule{}
}

func (r *TerraformBackendS3UseLockfileRule) Name() string {
	return "terraform_backend_s3_use_lockfile"
}

func (r *TerraformBackendS3UseLockfileRule) Enabled() bool {
	return true
}

func (r *TerraformBackendS3UseLockfileRule) Severity() tflint.Severity {
	return tflint.ERROR
}

func (r *TerraformBackendS3UseLockfileRule) Check(runner tflint.Runner) error {
	blocks, err := s3BackendBlocks(runner)
	if err != nil {
		return err
	}

	for _, block := range blocks {
		empty, err := s3BackendBlockIsEmpty(runner, block)
		if err != nil {
			return err
		}
		if empty {
			continue
		}

		attr, ok := block.Body.Attributes["use_lockfile"]
		if !ok {
			if err := runner.EmitIssue(r, terraformBackendS3UseLockfileMessage, block.DefRange); err != nil {
				return err
			}
			continue
		}

		value, known, err := evaluatedBool(runner, attr.Expr)
		if err != nil {
			return err
		}
		if !known || !value {
			if err := runner.EmitIssue(r, terraformBackendS3UseLockfileMessage, attr.Expr.Range()); err != nil {
				return err
			}
		}
	}
	return nil
}

var _ tflint.Rule = (*TerraformBackendS3UseLockfileRule)(nil)

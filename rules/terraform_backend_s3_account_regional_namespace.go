package rules

import (
	"regexp"

	"github.com/terraform-linters/tflint-plugin-sdk/tflint"
)

const terraformBackendS3AccountRegionalNamespaceMessage = "S3 backend buckets must use the account-regional namespace"

var accountRegionalBucketPattern = regexp.MustCompile(`^.+-[0-9]{12}-[a-z0-9]+(-[a-z0-9]+)*-[0-9]+-an$`)

type TerraformBackendS3AccountRegionalNamespaceRule struct {
	tflint.DefaultRule
}

func NewTerraformBackendS3AccountRegionalNamespaceRule() *TerraformBackendS3AccountRegionalNamespaceRule {
	return &TerraformBackendS3AccountRegionalNamespaceRule{}
}

func (r *TerraformBackendS3AccountRegionalNamespaceRule) Name() string {
	return "terraform_backend_s3_account_regional_namespace"
}

func (r *TerraformBackendS3AccountRegionalNamespaceRule) Enabled() bool {
	return true
}

func (r *TerraformBackendS3AccountRegionalNamespaceRule) Severity() tflint.Severity {
	return tflint.ERROR
}

func (r *TerraformBackendS3AccountRegionalNamespaceRule) Check(runner tflint.Runner) error {
	blocks, err := s3BackendBlocks(runner)
	if err != nil {
		return err
	}

	for _, block := range blocks {
		attr, ok := block.Body.Attributes["bucket"]
		if !ok {
			if err := runner.EmitIssue(r, terraformBackendS3AccountRegionalNamespaceMessage, block.DefRange); err != nil {
				return err
			}
			continue
		}

		value, known, err := evaluatedString(runner, attr.Expr)
		if err != nil {
			return err
		}
		if !known || !accountRegionalBucketPattern.MatchString(value) {
			if err := runner.EmitIssue(r, terraformBackendS3AccountRegionalNamespaceMessage, attr.Expr.Range()); err != nil {
				return err
			}
		}
	}
	return nil
}

var _ tflint.Rule = (*TerraformBackendS3AccountRegionalNamespaceRule)(nil)

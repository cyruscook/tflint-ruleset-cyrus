package rules

import (
	"regexp"

	"github.com/hashicorp/hcl/v2"
	"github.com/terraform-linters/tflint-plugin-sdk/tflint"
)

const awsNoHardcodedServicePrincipalMessage = "AWS service principals must use data.aws_service_principal instead of a hardcoded hostname"

var awsHardcodedServicePrincipalPattern = regexp.MustCompile(`(^|["'[:space:]])[a-z0-9-]+(\.[a-z0-9-]+)*\.amazonaws\.com(\.cn)?($|["'[:space:]])`)

type AwsNoHardcodedServicePrincipalRule struct {
	tflint.DefaultRule
}

func NewAwsNoHardcodedServicePrincipalRule() *AwsNoHardcodedServicePrincipalRule {
	return &AwsNoHardcodedServicePrincipalRule{}
}

func (r *AwsNoHardcodedServicePrincipalRule) Name() string {
	return "aws_no_hardcoded_service_principal"
}

func (r *AwsNoHardcodedServicePrincipalRule) Enabled() bool {
	return true
}

func (r *AwsNoHardcodedServicePrincipalRule) Severity() tflint.Severity {
	return tflint.ERROR
}

func (r *AwsNoHardcodedServicePrincipalRule) Check(runner tflint.Runner) error {
	return walkLiteralStrings(runner, func(value string, rng hcl.Range) error {
		if awsHardcodedServicePrincipalPattern.MatchString(value) {
			return runner.EmitIssue(r, awsNoHardcodedServicePrincipalMessage, rng)
		}
		return nil
	})
}

var _ tflint.Rule = (*AwsNoHardcodedServicePrincipalRule)(nil)

package rules

import (
	"regexp"

	"github.com/hashicorp/hcl/v2"
	"github.com/terraform-linters/tflint-plugin-sdk/tflint"
)

const awsNoHardcodedPartitionMessage = "AWS ARN partitions must use data.aws_partition instead of a hardcoded partition"

var awsHardcodedPartitionPattern = regexp.MustCompile(`arn:aws(-[a-z0-9-]+)?:`)

type AwsNoHardcodedPartitionRule struct {
	tflint.DefaultRule
}

func NewAwsNoHardcodedPartitionRule() *AwsNoHardcodedPartitionRule {
	return &AwsNoHardcodedPartitionRule{}
}

func (r *AwsNoHardcodedPartitionRule) Name() string {
	return "aws_no_hardcoded_partition"
}

func (r *AwsNoHardcodedPartitionRule) Enabled() bool {
	return true
}

func (r *AwsNoHardcodedPartitionRule) Severity() tflint.Severity {
	return tflint.ERROR
}

func (r *AwsNoHardcodedPartitionRule) Check(runner tflint.Runner) error {
	return walkLiteralStrings(runner, func(value string, rng hcl.Range) error {
		if awsHardcodedPartitionPattern.MatchString(value) {
			return runner.EmitIssue(r, awsNoHardcodedPartitionMessage, rng)
		}
		return nil
	})
}

var _ tflint.Rule = (*AwsNoHardcodedPartitionRule)(nil)

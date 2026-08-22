package rules

import (
	"github.com/terraform-linters/tflint-plugin-sdk/hclext"
	"github.com/terraform-linters/tflint-plugin-sdk/tflint"
)

const awsS3BucketAccountRegionalNamespaceMessage = "S3 buckets must set bucket_namespace to \"account-regional\""

type AwsS3BucketAccountRegionalNamespaceRule struct {
	tflint.DefaultRule
}

func NewAwsS3BucketAccountRegionalNamespaceRule() *AwsS3BucketAccountRegionalNamespaceRule {
	return &AwsS3BucketAccountRegionalNamespaceRule{}
}

func (r *AwsS3BucketAccountRegionalNamespaceRule) Name() string {
	return "aws_s3_bucket_account_regional_namespace"
}

func (r *AwsS3BucketAccountRegionalNamespaceRule) Enabled() bool {
	return true
}

func (r *AwsS3BucketAccountRegionalNamespaceRule) Severity() tflint.Severity {
	return tflint.ERROR
}

func (r *AwsS3BucketAccountRegionalNamespaceRule) Check(runner tflint.Runner) error {
	content, err := runner.GetResourceContent("aws_s3_bucket", &hclext.BodySchema{
		Attributes: []hclext.AttributeSchema{{Name: "bucket_namespace", Required: false}},
	}, nil)
	if err != nil {
		return err
	}

	for _, resource := range content.Blocks {
		attr, ok := resource.Body.Attributes["bucket_namespace"]
		if !ok {
			if err := runner.EmitIssue(r, awsS3BucketAccountRegionalNamespaceMessage, resource.DefRange); err != nil {
				return err
			}
			continue
		}

		value, known, err := evaluatedString(runner, attr.Expr)
		if err != nil {
			return err
		}
		if !known || value != "account-regional" {
			if err := runner.EmitIssue(r, awsS3BucketAccountRegionalNamespaceMessage, attr.Expr.Range()); err != nil {
				return err
			}
		}
	}
	return nil
}

var _ tflint.Rule = (*AwsS3BucketAccountRegionalNamespaceRule)(nil)

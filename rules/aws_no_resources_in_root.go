package rules

import (
	"strings"

	"github.com/terraform-linters/tflint-plugin-sdk/hclext"
	"github.com/terraform-linters/tflint-plugin-sdk/tflint"
)

const awsNoResourcesInRootMessage = "AWS resources must be created in modules, not in the root module"

type AwsNoResourcesInRootRule struct {
	tflint.DefaultRule
}

func NewAwsNoResourcesInRootRule() *AwsNoResourcesInRootRule {
	return &AwsNoResourcesInRootRule{}
}

func (r *AwsNoResourcesInRootRule) Name() string {
	return "aws_no_resources_in_root"
}

func (r *AwsNoResourcesInRootRule) Enabled() bool {
	return true
}

func (r *AwsNoResourcesInRootRule) Severity() tflint.Severity {
	return tflint.ERROR
}

func (r *AwsNoResourcesInRootRule) Check(runner tflint.Runner) error {
	modulePath, err := runner.GetModulePath()
	if err != nil {
		return err
	}
	if !modulePath.IsRoot() {
		return nil
	}

	content, err := runner.GetModuleContent(&hclext.BodySchema{
		Blocks: []hclext.BlockSchema{
			{
				Type:       "resource",
				LabelNames: []string{"type", "name"},
				Body:       &hclext.BodySchema{},
			},
		},
	}, &tflint.GetModuleContentOption{ExpandMode: tflint.ExpandModeNone})
	if err != nil {
		return err
	}

	for _, resource := range content.Blocks {
		if !strings.HasPrefix(resource.Labels[0], "aws_") {
			continue
		}
		if err := runner.EmitIssue(r, awsNoResourcesInRootMessage, resource.DefRange); err != nil {
			return err
		}
	}
	return nil
}

var _ tflint.Rule = (*AwsNoResourcesInRootRule)(nil)

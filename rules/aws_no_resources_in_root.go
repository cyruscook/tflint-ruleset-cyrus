package rules

import (
	"github.com/terraform-linters/tflint-plugin-sdk/hclext"
	"github.com/terraform-linters/tflint-plugin-sdk/tflint"
	"path"
	"path/filepath"
	"strings"
)

const awsNoResourcesInRootMessage = "AWS resources must be created in modules, not in the root module"

type AwsNoResourcesInRootRule struct {
	tflint.DefaultRule
}

type awsNoResourcesInRootRuleConfig struct {
	RootModulePaths []string `hclext:"root_module_paths,optional"`
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

	config := &awsNoResourcesInRootRuleConfig{}
	if err := runner.DecodeRuleConfig(r.Name(), config); err != nil {
		return err
	}

	var originalwd string
	if len(config.RootModulePaths) > 0 {
		originalwd, err = runner.GetOriginalwd()
		if err != nil {
			return err
		}
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
		if len(config.RootModulePaths) > 0 {
			isRoot, err := matchesRootModulePath(originalwd, resource.DefRange.Filename, config.RootModulePaths)
			if err != nil {
				return err
			}
			if !isRoot {
				continue
			}
		}
		if err := runner.EmitIssue(r, awsNoResourcesInRootMessage, resource.DefRange); err != nil {
			return err
		}
	}
	return nil
}

func matchesRootModulePath(originalwd, filename string, patterns []string) (bool, error) {
	if !filepath.IsAbs(filename) {
		filename = filepath.Join(originalwd, filename)
	}
	modulePath, err := filepath.Rel(originalwd, filepath.Dir(filename))
	if err != nil {
		return false, err
	}
	modulePath = filepath.ToSlash(modulePath)

	for _, pattern := range patterns {
		matches, err := path.Match(filepath.ToSlash(pattern), modulePath)
		if err != nil {
			return false, err
		}
		if matches {
			return true, nil
		}
	}
	return false, nil
}

var _ tflint.Rule = (*AwsNoResourcesInRootRule)(nil)

package rules

import (
	"testing"

	"github.com/hashicorp/hcl/v2"
	"github.com/terraform-linters/tflint-plugin-sdk/helper"
	"github.com/terraform-linters/tflint-plugin-sdk/terraform/addrs"
	"github.com/terraform-linters/tflint-plugin-sdk/tflint"
)

type childModuleRunner struct {
	*helper.Runner
}

func (r *childModuleRunner) GetModulePath() (addrs.Module, error) {
	return addrs.Module{"child"}, nil
}

func TestAwsNoResourcesInRoot(t *testing.T) {
	tests := []struct {
		name  string
		src   string
		child bool
		want  helper.Issues
	}{
		{
			name: "root AWS resource",
			src:  `resource "aws_s3_bucket" "example" {}`,
			want: helper.Issues{&helper.Issue{Rule: NewAwsNoResourcesInRootRule(), Message: awsNoResourcesInRootMessage, Range: hcl.Range{Filename: "main.tf", Start: hcl.Pos{Line: 1, Column: 1}, End: hcl.Pos{Line: 1, Column: 35}}}},
		},
		{
			name: "multiple root AWS resources",
			src: `resource "aws_s3_bucket" "first" {}
resource "aws_s3_bucket" "second" {}`,
			want: helper.Issues{
				&helper.Issue{Rule: NewAwsNoResourcesInRootRule(), Message: awsNoResourcesInRootMessage, Range: hcl.Range{Filename: "main.tf", Start: hcl.Pos{Line: 1, Column: 1}, End: hcl.Pos{Line: 1, Column: 33}}},
				&helper.Issue{Rule: NewAwsNoResourcesInRootRule(), Message: awsNoResourcesInRootMessage, Range: hcl.Range{Filename: "main.tf", Start: hcl.Pos{Line: 2, Column: 1}, End: hcl.Pos{Line: 2, Column: 34}}},
			},
		},
		{
			name: "data non-AWS and module blocks ignored",
			src: `data "aws_caller_identity" "current" {}
resource "random_id" "example" {}
module "example" {}`,
			want: helper.Issues{},
		},
		{
			name: "count reports one declaration issue",
			src: `resource "aws_s3_bucket" "example" {
  count = 1
}`,
			want: helper.Issues{&helper.Issue{Rule: NewAwsNoResourcesInRootRule(), Message: awsNoResourcesInRootMessage, Range: hcl.Range{Filename: "main.tf", Start: hcl.Pos{Line: 1, Column: 1}, End: hcl.Pos{Line: 1, Column: 35}}}},
		},
		{
			name:  "child module AWS resource ignored",
			src:   `resource "aws_s3_bucket" "example" {}`,
			child: true,
			want:  helper.Issues{},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			baseRunner := helper.TestRunner(t, map[string]string{"main.tf": test.src})
			var runner tflint.Runner = baseRunner
			if test.child {
				runner = &childModuleRunner{Runner: baseRunner}
			}
			if err := NewAwsNoResourcesInRootRule().Check(runner); err != nil {
				t.Fatal(err)
			}
			helper.AssertIssues(t, test.want, baseRunner.Issues)
		})
	}
}

package rules

import (
	"testing"

	"github.com/hashicorp/hcl/v2"
	"github.com/terraform-linters/tflint-plugin-sdk/helper"
)

func TestAwsS3BucketAccountRegionalNamespace(t *testing.T) {
	tests := []struct {
		name string
		src  string
		want helper.Issues
	}{
		{
			name: "account regional passes",
			src: `resource "aws_s3_bucket" "state" {
  bucket_namespace = "account-regional"
}`,
			want: helper.Issues{},
		},
		{
			name: "missing fails",
			src:  `resource "aws_s3_bucket" "state" {}`,
			want: helper.Issues{&helper.Issue{Rule: NewAwsS3BucketAccountRegionalNamespaceRule(), Message: awsS3BucketAccountRegionalNamespaceMessage, Range: hcl.Range{Filename: "main.tf", Start: hcl.Pos{Line: 1, Column: 1}, End: hcl.Pos{Line: 1, Column: 33}}}},
		},
		{
			name: "global fails",
			src: `resource "aws_s3_bucket" "state" {
  bucket_namespace = "global"
}`,
			want: helper.Issues{&helper.Issue{Rule: NewAwsS3BucketAccountRegionalNamespaceRule(), Message: awsS3BucketAccountRegionalNamespaceMessage, Range: hcl.Range{Filename: "main.tf", Start: hcl.Pos{Line: 2, Column: 22}, End: hcl.Pos{Line: 2, Column: 30}}}},
		},
		{
			name: "unknown fails",
			src: `variable "namespace" {}
resource "aws_s3_bucket" "state" {
  bucket_namespace = var.namespace
}`,
			want: helper.Issues{&helper.Issue{Rule: NewAwsS3BucketAccountRegionalNamespaceRule(), Message: awsS3BucketAccountRegionalNamespaceMessage, Range: hcl.Range{Filename: "main.tf", Start: hcl.Pos{Line: 3, Column: 22}, End: hcl.Pos{Line: 3, Column: 35}}}},
		},
		{
			name: "unrelated resource ignored",
			src:  `resource "aws_s3_object" "object" {}`,
			want: helper.Issues{},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			runner := helper.TestRunner(t, map[string]string{"main.tf": test.src})
			if err := NewAwsS3BucketAccountRegionalNamespaceRule().Check(runner); err != nil {
				t.Fatal(err)
			}
			helper.AssertIssues(t, test.want, runner.Issues)
		})
	}
}

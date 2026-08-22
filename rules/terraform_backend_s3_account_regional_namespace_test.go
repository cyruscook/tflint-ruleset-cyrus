package rules

import (
	"testing"

	"github.com/hashicorp/hcl/v2"
	"github.com/terraform-linters/tflint-plugin-sdk/helper"
)

func TestTerraformBackendS3AccountRegionalNamespace(t *testing.T) {
	tests := []struct {
		name string
		src  string
		want helper.Issues
	}{
		{
			name: "standard region passes",
			src: `terraform {
  backend "s3" {
    bucket = "state-123456789012-us-east-1-an"
  }
}`,
			want: helper.Issues{},
		},
		{
			name: "govcloud region passes",
			src: `terraform {
  backend "s3" {
    bucket = "state-123456789012-us-gov-west-1-an"
  }
}`,
			want: helper.Issues{},
		},
		{
			name: "global name fails",
			src: `terraform {
  backend "s3" {
    bucket = "state-global"
  }
}`,
			want: helper.Issues{&helper.Issue{Rule: NewTerraformBackendS3AccountRegionalNamespaceRule(), Message: terraformBackendS3AccountRegionalNamespaceMessage, Range: hcl.Range{Filename: "main.tf", Start: hcl.Pos{Line: 3, Column: 14}, End: hcl.Pos{Line: 3, Column: 28}}}},
		},
		{
			name: "missing bucket fails",
			src: `terraform {
  backend "s3" {}
}`,
			want: helper.Issues{&helper.Issue{Rule: NewTerraformBackendS3AccountRegionalNamespaceRule(), Message: terraformBackendS3AccountRegionalNamespaceMessage, Range: hcl.Range{Filename: "main.tf", Start: hcl.Pos{Line: 2, Column: 3}, End: hcl.Pos{Line: 2, Column: 15}}}},
		},
		{
			name: "unknown bucket fails and non s3 ignored",
			src: `variable "bucket" {}
terraform {
  backend "s3" {
    bucket = var.bucket
  }
  backend "local" {}
}`,
			want: helper.Issues{&helper.Issue{Rule: NewTerraformBackendS3AccountRegionalNamespaceRule(), Message: terraformBackendS3AccountRegionalNamespaceMessage, Range: hcl.Range{Filename: "main.tf", Start: hcl.Pos{Line: 4, Column: 14}, End: hcl.Pos{Line: 4, Column: 24}}}},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			runner := helper.TestRunner(t, map[string]string{"main.tf": test.src})
			if err := NewTerraformBackendS3AccountRegionalNamespaceRule().Check(runner); err != nil {
				t.Fatal(err)
			}
			helper.AssertIssues(t, test.want, runner.Issues)
		})
	}
}

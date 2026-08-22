package rules

import (
	"testing"

	"github.com/hashicorp/hcl/v2"
	"github.com/terraform-linters/tflint-plugin-sdk/helper"
)

func TestTerraformBackendS3UseLockfile(t *testing.T) {
	tests := []struct {
		name string
		src  string
		want helper.Issues
	}{
		{
			name: "true passes",
			src: `terraform {
  backend "s3" {
    use_lockfile = true
  }
}`,
			want: helper.Issues{},
		},
		{
			name: "false fails",
			src: `terraform {
  backend "s3" {
    use_lockfile = false
  }
}`,
			want: helper.Issues{&helper.Issue{Rule: NewTerraformBackendS3UseLockfileRule(), Message: terraformBackendS3UseLockfileMessage, Range: hcl.Range{Filename: "main.tf", Start: hcl.Pos{Line: 3, Column: 20}, End: hcl.Pos{Line: 3, Column: 25}}}},
		},
		{
			name: "missing fails",
			src: `terraform {
  backend "s3" {}
}`,
			want: helper.Issues{&helper.Issue{Rule: NewTerraformBackendS3UseLockfileRule(), Message: terraformBackendS3UseLockfileMessage, Range: hcl.Range{Filename: "main.tf", Start: hcl.Pos{Line: 2, Column: 3}, End: hcl.Pos{Line: 2, Column: 15}}}},
		},
		{
			name: "non boolean fails",
			src: `terraform {
  backend "s3" {
    use_lockfile = "true"
  }
}`,
			want: helper.Issues{&helper.Issue{Rule: NewTerraformBackendS3UseLockfileRule(), Message: terraformBackendS3UseLockfileMessage, Range: hcl.Range{Filename: "main.tf", Start: hcl.Pos{Line: 3, Column: 20}, End: hcl.Pos{Line: 3, Column: 26}}}},
		},
		{
			name: "unknown fails and non s3 ignored",
			src: `variable "lock" {}
terraform {
  backend "s3" {
    use_lockfile = var.lock
  }
  backend "local" {}
}`,
			want: helper.Issues{&helper.Issue{Rule: NewTerraformBackendS3UseLockfileRule(), Message: terraformBackendS3UseLockfileMessage, Range: hcl.Range{Filename: "main.tf", Start: hcl.Pos{Line: 4, Column: 20}, End: hcl.Pos{Line: 4, Column: 28}}}},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			runner := helper.TestRunner(t, map[string]string{"main.tf": test.src})
			if err := NewTerraformBackendS3UseLockfileRule().Check(runner); err != nil {
				t.Fatal(err)
			}
			helper.AssertIssues(t, test.want, runner.Issues)
		})
	}
}

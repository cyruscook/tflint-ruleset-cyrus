package rules

import (
	"testing"

	"github.com/hashicorp/hcl/v2"
	"github.com/terraform-linters/tflint-plugin-sdk/helper"
)

func TestTerraformVariableNullable(t *testing.T) {
	tests := []struct {
		name  string
		files map[string]string
		want  helper.Issues
	}{
		{
			name: "true passes",
			files: map[string]string{"main.tf": `variable "enabled" {
  nullable = true
}`},
			want: helper.Issues{},
		},
		{
			name: "false passes",
			files: map[string]string{"main.tf": `variable "enabled" {
  nullable = false
}`},
			want: helper.Issues{},
		},
		{
			name:  "missing fails",
			files: map[string]string{"main.tf": `variable "name" {}`},
			want:  helper.Issues{&helper.Issue{Rule: NewTerraformVariableNullableRule(), Message: terraformVariableNullableMessage, Range: hcl.Range{Filename: "main.tf", Start: hcl.Pos{Line: 1, Column: 1}, End: hcl.Pos{Line: 1, Column: 16}}}},
		},
		{
			name: "mixed declarations report only omissions",
			files: map[string]string{"main.tf": `variable "a" {
  nullable = true
}
variable "b" {}
variable "c" {
  nullable = false
}
variable "d" {}`},
			want: helper.Issues{
				&helper.Issue{Rule: NewTerraformVariableNullableRule(), Message: terraformVariableNullableMessage, Range: hcl.Range{Filename: "main.tf", Start: hcl.Pos{Line: 4, Column: 1}, End: hcl.Pos{Line: 4, Column: 13}}},
				&helper.Issue{Rule: NewTerraformVariableNullableRule(), Message: terraformVariableNullableMessage, Range: hcl.Range{Filename: "main.tf", Start: hcl.Pos{Line: 8, Column: 1}, End: hcl.Pos{Line: 8, Column: 13}}},
			},
		},
		{
			name:  "JSON present passes",
			files: map[string]string{"main.tf.json": `{"variable":{"name":{"nullable":true}}}`},
			want:  helper.Issues{},
		},
		{
			name:  "JSON missing fails",
			files: map[string]string{"main.tf.json": `{"variable":{"name":{}}}`},
			want:  helper.Issues{&helper.Issue{Rule: NewTerraformVariableNullableRule(), Message: terraformVariableNullableMessage, Range: hcl.Range{Filename: "main.tf.json", Start: hcl.Pos{Line: 1, Column: 21}, End: hcl.Pos{Line: 1, Column: 22}}}},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			runner := helper.TestRunner(t, test.files)
			if err := NewTerraformVariableNullableRule().Check(runner); err != nil {
				t.Fatal(err)
			}
			helper.AssertIssues(t, test.want, runner.Issues)
		})
	}
}

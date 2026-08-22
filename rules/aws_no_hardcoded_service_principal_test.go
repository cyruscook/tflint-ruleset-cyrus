package rules

import (
	"testing"

	"github.com/hashicorp/hcl/v2"
	"github.com/terraform-linters/tflint-plugin-sdk/helper"
)

func TestAwsNoHardcodedServicePrincipal(t *testing.T) {
	tests := []struct {
		name  string
		files map[string]string
		want  helper.Issues
	}{
		{
			name:  "direct principal",
			files: map[string]string{"main.tf": `locals { principal = "lambda.amazonaws.com" }`},
			want:  helper.Issues{&helper.Issue{Rule: NewAwsNoHardcodedServicePrincipalRule(), Message: awsNoHardcodedServicePrincipalMessage, Range: hcl.Range{Filename: "main.tf", Start: hcl.Pos{Line: 1, Column: 23}, End: hcl.Pos{Line: 1, Column: 43}}}},
		},
		{
			name:  "regional and multi label principal",
			files: map[string]string{"main.tf": `locals { principal = "service.us-east-1.amazonaws.com" }`},
			want:  helper.Issues{&helper.Issue{Rule: NewAwsNoHardcodedServicePrincipalRule(), Message: awsNoHardcodedServicePrincipalMessage, Range: hcl.Range{Filename: "main.tf", Start: hcl.Pos{Line: 1, Column: 23}, End: hcl.Pos{Line: 1, Column: 54}}}},
		},
		{
			name:  "china principal",
			files: map[string]string{"main.tf": `locals { principal = "service.amazonaws.com.cn" }`},
			want:  helper.Issues{&helper.Issue{Rule: NewAwsNoHardcodedServicePrincipalRule(), Message: awsNoHardcodedServicePrincipalMessage, Range: hcl.Range{Filename: "main.tf", Start: hcl.Pos{Line: 1, Column: 23}, End: hcl.Pos{Line: 1, Column: 47}}}},
		},
		{
			name: "policy heredoc",
			files: map[string]string{"main.tf": `locals {
  policy = <<JSON
{"Principal":"lambda.amazonaws.com"}
JSON
}`},
			want: helper.Issues{&helper.Issue{Rule: NewAwsNoHardcodedServicePrincipalRule(), Message: awsNoHardcodedServicePrincipalMessage, Range: hcl.Range{Filename: "main.tf", Start: hcl.Pos{Line: 3, Column: 1}, End: hcl.Pos{Line: 4, Column: 1}}}},
		},
		{
			name:  "data principal passes",
			files: map[string]string{"main.tf": `locals { principal = data.aws_service_principal.lambda.name }`},
			want:  helper.Issues{},
		},
		{
			name: "comments unrelated and URL pass",
			files: map[string]string{"main.tf": `# lambda.amazonaws.com
locals {
  example = "example.com"
  endpoint = "https://service.amazonaws.com/path"
}`},
			want: helper.Issues{},
		},
		{
			name:  "JSON constant string",
			files: map[string]string{"main.tf.json": `{"principal":"lambda.amazonaws.com"}`},
			want:  helper.Issues{&helper.Issue{Rule: NewAwsNoHardcodedServicePrincipalRule(), Message: awsNoHardcodedServicePrincipalMessage, Range: hcl.Range{Filename: "main.tf.json", Start: hcl.Pos{Line: 1, Column: 14}, End: hcl.Pos{Line: 1, Column: 36}}}},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			runner := helper.TestRunner(t, test.files)
			if err := NewAwsNoHardcodedServicePrincipalRule().Check(runner); err != nil {
				t.Fatal(err)
			}
			helper.AssertIssues(t, test.want, runner.Issues)
		})
	}
}

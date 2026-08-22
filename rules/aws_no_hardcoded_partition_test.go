package rules

import (
	"testing"

	"github.com/hashicorp/hcl/v2"
	"github.com/terraform-linters/tflint-plugin-sdk/helper"
)

func TestAwsNoHardcodedPartition(t *testing.T) {
	tests := []struct {
		name  string
		files map[string]string
		want  helper.Issues
	}{
		{
			name:  "direct commercial ARN",
			files: map[string]string{"main.tf": `locals { arn = "arn:aws:s3:::example" }`},
			want:  helper.Issues{&helper.Issue{Rule: NewAwsNoHardcodedPartitionRule(), Message: awsNoHardcodedPartitionMessage, Range: hcl.Range{Filename: "main.tf", Start: hcl.Pos{Line: 1, Column: 17}, End: hcl.Pos{Line: 1, Column: 37}}}},
		},
		{
			name:  "govcloud ARN",
			files: map[string]string{"main.tf": `locals { arn = "arn:aws-us-gov:s3:::example" }`},
			want:  helper.Issues{&helper.Issue{Rule: NewAwsNoHardcodedPartitionRule(), Message: awsNoHardcodedPartitionMessage, Range: hcl.Range{Filename: "main.tf", Start: hcl.Pos{Line: 1, Column: 17}, End: hcl.Pos{Line: 1, Column: 44}}}},
		},
		{
			name: "policy heredoc",
			files: map[string]string{"main.tf": `locals {
  policy = <<JSON
{"Resource":"arn:aws:s3:::example"}
JSON
}`},
			want: helper.Issues{&helper.Issue{Rule: NewAwsNoHardcodedPartitionRule(), Message: awsNoHardcodedPartitionMessage, Range: hcl.Range{Filename: "main.tf", Start: hcl.Pos{Line: 3, Column: 1}, End: hcl.Pos{Line: 4, Column: 1}}}},
		},
		{
			name:  "template literal segment",
			files: map[string]string{"main.tf": `locals { arn = "arn:aws:${var.rest}" }`},
			want:  helper.Issues{&helper.Issue{Rule: NewAwsNoHardcodedPartitionRule(), Message: awsNoHardcodedPartitionMessage, Range: hcl.Range{Filename: "main.tf", Start: hcl.Pos{Line: 1, Column: 17}, End: hcl.Pos{Line: 1, Column: 25}}}},
		},
		{
			name:  "dynamic partition passes",
			files: map[string]string{"main.tf": `locals { arn = "arn:${data.aws_partition.current.partition}:s3:::example" }`},
			want:  helper.Issues{},
		},
		{
			name: "comments and non ARN pass",
			files: map[string]string{"main.tf": `# arn:aws:s3:::comment
locals { text = "not an ARN" }`},
			want: helper.Issues{},
		},
		{
			name:  "JSON constant string",
			files: map[string]string{"main.tf.json": `{"arn":"arn:aws:s3:::example"}`},
			want:  helper.Issues{&helper.Issue{Rule: NewAwsNoHardcodedPartitionRule(), Message: awsNoHardcodedPartitionMessage, Range: hcl.Range{Filename: "main.tf.json", Start: hcl.Pos{Line: 1, Column: 8}, End: hcl.Pos{Line: 1, Column: 30}}}},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			runner := helper.TestRunner(t, test.files)
			if err := NewAwsNoHardcodedPartitionRule().Check(runner); err != nil {
				t.Fatal(err)
			}
			helper.AssertIssues(t, test.want, runner.Issues)
		})
	}
}

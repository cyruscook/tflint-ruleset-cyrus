package main

import (
	"github.com/terraform-linters/tflint-plugin-sdk/plugin"
	"github.com/terraform-linters/tflint-plugin-sdk/tflint"
	"github.com/terraform-linters/tflint-ruleset-cyrus/rules"
)

func main() {
	plugin.Serve(&plugin.ServeOpts{
		RuleSet: &tflint.BuiltinRuleSet{
			Name:    "cyrus",
			Version: version,
			Rules: []tflint.Rule{
				rules.NewAwsNoHardcodedPartitionRule(),
				rules.NewAwsNoHardcodedServicePrincipalRule(),
				rules.NewAwsS3BucketAccountRegionalNamespaceRule(),
				rules.NewTerraformBackendS3AccountRegionalNamespaceRule(),
				rules.NewTerraformBackendS3UseLockfileRule(),
				rules.NewTerraformVariableNullableRule(),
			},
		},
	})
}

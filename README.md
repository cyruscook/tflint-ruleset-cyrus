# TFLint Ruleset Cyrus

[![Build Status](https://github.com/terraform-linters/tflint-ruleset-cyrus/actions/workflows/build.yml/badge.svg?branch=main)](https://github.com/terraform-linters/tflint-ruleset-cyrus/actions)

A custom [TFLint](https://github.com/terraform-linters/tflint) ruleset plugin built with the [TFLint Plugin SDK](https://github.com/terraform-linters/tflint-plugin-sdk).

The repository starts with reference rules covering resource attributes, nested blocks, rule configuration, and non-resource blocks. Replace or extend the rules in `rules/` as the ruleset evolves.

## Requirements

- TFLint v0.46+
- Go v1.26+

## Installation

After publishing a GitHub release that follows TFLint's plugin asset naming requirements, install the ruleset with `tflint --init`:

```hcl
plugin "cyrus" {
  enabled = true
  version = "0.1.0"
  source  = "github.com/terraform-linters/tflint-ruleset-cyrus"
}
```

For local development, build and install the plugin directly:

```console
$ make install
```

Then enable it in `.tflint.hcl`:

```hcl
plugin "cyrus" {
  enabled = true
}
```

## Rules

| Name | Description | Severity | Enabled |
| --- | --- | --- | --- |
| `aws_instance_example_type` | Example rule for accessing and evaluating top-level attributes | ERROR | Yes |
| `aws_s3_bucket_example_lifecycle_rule` | Example rule for accessing top-level and nested blocks and attributes | ERROR | Yes |
| `google_compute_ssl_policy` | Example rule with a custom rule configuration | WARNING | No |
| `terraform_backend_type` | Example rule for accessing non-resource blocks | ERROR | Yes |

The `google_compute_ssl_policy` rule requires configuration in `.tflint.hcl`:

```hcl
rule "google_compute_ssl_policy" {
  enabled          = true
  allowed_versions = ["TLS_1_2"]
}
```

## Building and testing

```console
$ make test
$ make build
```

The resulting binary is `tflint-ruleset-cyrus`. `make install` places it in `~/.tflint.d/plugins`.

## Releasing

Tag a release as `vMAJOR.MINOR.PATCH`. The included GoReleaser configuration produces platform archives, `checksums.txt`, and checksum attestation/signature inputs expected by TFLint's automatic plugin installation.

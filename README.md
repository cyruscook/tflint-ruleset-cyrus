# TFLint Ruleset Cyrus

Custom opinionated Terraform ruleset.

## Requirements

- TFLint v0.46+
- Go v1.26+

## Installation

Install the ruleset with `tflint --init`:

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

| Name                                              | Description                                                                                  | Severity | Enabled |
| ------------------------------------------------- | -------------------------------------------------------------------------------------------- | -------- | ------- |
| `aws_no_hardcoded_partition`                      | AWS ARNs must use `data.aws_partition` instead of a hardcoded partition                      | ERROR    | Yes     |
| `aws_no_hardcoded_service_principal`              | AWS service principals must use `data.aws_service_principal` instead of a hardcoded hostname | ERROR    | Yes     |
| `aws_no_resources_in_root`                        | AWS resources must be created in modules, not in the root module                             | ERROR    | Yes     |
| `aws_s3_bucket_account_regional_namespace`        | `aws_s3_bucket` resources must set `bucket_namespace = "account-regional"`                   | ERROR    | Yes     |
| `terraform_backend_s3_account_regional_namespace` | S3 backend buckets must use the `-<12-digit-account>-<region>-an` namespace suffix           | ERROR    | Yes     |
| `terraform_backend_s3_use_lockfile`               | S3 backends must set `use_lockfile = true`                                                   | ERROR    | Yes     |

## Building and testing

```console
$ make test
$ make build
```

The resulting binary is `tflint-ruleset-cyrus`. `make install` places it in `~/.tflint.d/plugins`.

## Releasing

Tag a release as `vMAJOR.MINOR.PATCH`. The included GoReleaser configuration produces platform archives, `checksums.txt`, and checksum attestation/signature inputs expected by TFLint's automatic plugin installation.

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
  version = "0.6.0" # x-release-please-version
  source  = "github.com/cyruscook/tflint-ruleset-cyrus"
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
| `aws_s3_bucket_account_regional_namespace`        | `aws_s3_bucket` resources must set `bucket_namespace = "account-regional"`                   | ERROR    | Yes     |
| `terraform_backend_s3_account_regional_namespace` | S3 backend buckets must use the `-<12-digit-account>-<region>-an` namespace suffix; empty backend blocks are allowed for external configuration | ERROR    | Yes     |
| `terraform_backend_s3_use_lockfile`               | S3 backends must set `use_lockfile = true`; empty backend blocks are allowed for external configuration                            | ERROR    | Yes     |
| `terraform_variable_nullable`                     | Variables must explicitly set the nullable attribute                                      | ERROR    | Yes     |

## Building and testing

```console
$ make test
$ make build
```

The resulting binary is `tflint-ruleset-cyrus`. `make install` places it in `~/.tflint.d/plugins`.

## Releasing

[Release Please](https://github.com/googleapis/release-please-action) maintains a
release PR with `CHANGELOG.md`, the version in `version.go`, the installation
example above, and `.release-please-manifest.json`. It starts from the existing
v0.6.0 release; it does not recreate historical releases.

1. Use conventional PR titles and squash-merge so the title becomes the commit
   on `main`: `fix: report the correct plugin version` produces a patch release;
   `feat: add a rule` produces a minor release. While below 1.0, `feat!:` or a
   `BREAKING CHANGE:` footer produces a minor release. Ordinary `chore:`, `ci:`,
   and `docs:` changes do not independently request a release. Plain Dependabot
   `Bump ...` titles are not conventional commits: use `fix(deps): ...` for a
   dependency update that should produce a plugin release, or `ci(deps): ...`
   for workflow-only updates.
2. Review the generated release PR's version, changelog, and validation. Keep its
   generated title and release labels. Merge it when ready to publish; there is
   no auto-merge. That merge is the publication decision.
3. The `release-please` workflow creates the tag and GitHub release, then directly
   calls the reusable `release` workflow for that exact tag. GoReleaser uploads
   the platform archives, checksums, and GPG signature, retaining Release
   Please's notes. The checksum attestation runs as before. Release Please's
   `GITHUB_TOKEN` tag does not need to trigger a second workflow.

The workflow uses the built-in `GITHUB_TOKEN` with write permissions scoped to
release management, and the existing `GPG_PRIVATE_KEY` and `PASSPHRASE` secrets
for publishing. No personal token or GitHub App is required. The repository must
already allow GitHub Actions to create pull requests. If policy disables that,
release PR creation will fail; this change does not alter repository settings.

Bot-created PR workflows may require a maintainer to select **Approve workflows
to run**. Independently, `validate-release-pr` runs tests, builds, and the version
injection check on the exact generated PR head with read-only permissions and
no publishing secrets. These reusable-workflow results belong to the originating
`main` run, not the PR head's required checks. If branch protection requires PR
checks, approve the pending PR workflows and wait for those checks before merging.
No workflow approvals or protection changes are automated.

### Manual tags and recovery

A human-pushed `vMAJOR.MINOR.PATCH` tag still invokes the same publisher. Use this
only as an intentional alternative to a Release Please release, and synchronize
`version.go`, the README example, changelog, and manifest on `main` when bypassing
the release PR process. Do not manually tag an unmerged release PR. GoReleaser
embeds the tag's version through linker flags, even on this manual path; ordinary
`make build` uses the source version maintained by Release Please.

Publishing is serialized per tag. A release with any existing assets is skipped
to prevent overwrites or duplicate uploads. A newly created Release Please
release without assets is filled by GoReleaser, rather than creating another
release. Inspect the workflow result as well as the release page: the GitHub
release can exist before asset publication finishes. For a failure before any
assets upload, rerun the failed publishing job (rerunning Release Please alone
will not emit `release_created` again). After partial uploads or an attestation
failure, inspect and recover manually; reruns do not overwrite existing assets.

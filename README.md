# gcp_tf_iam_binding_validator

This is dead simple but useful command line tool which checks duplicated role in [google_project_iam_binding](https://registry.terraform.io/providers/hashicorp/google/latest/docs/resources/google_project_iam) in your terraform files.

## Why need check?

In [google_project_iam | Resources | hashicorp/google | Terraform Registry](https://registry.terraform.io/providers/hashicorp/google/latest/docs/resources/google_project_iam), it says:

> google_project_iam_binding: Authoritative for a given role. Updates the IAM policy to grant a role to a list of members. Other roles within the IAM policy for the project are preserved.

This means if you apply a `google_project_iam_binding` which contain the same role with existing one, the first one will be **replaced** with the new one. This is much surprising or even dangerous if you are not aware of the documentation.

(If you are not aware how it is dangerous, see [google cloud platform - Terraform google_project_iam_binding deletes GCP compute engine default service account from IAM principals - Stack Overflow](https://stackoverflow.com/questions/70703088/terraform-google-project-iam-binding-deletes-gcp-compute-engine-default-service).)

Even if you know the behavior, sometimes there can be a lot of .tf files in your workspace. In such cases, we want to make sure there are no `google_project_iam_binding` which already exists. This tool just does that.

## Installation

Just run:

```shell
go install github.com/hidetatz/gcp_tf_iam_binding_validator/cmd/gcp_tf_iam_binding_validator@latest
```

This tool is much intended to be used in your CI workflow.
This is an example for GitHub actions users:

```yaml
name: Check gcp_tf_iam_binding_validator

on:
  push:
    branches:
      - main
  pull_request:

jobs:
  test_gcp_tf_iam_binding_validator:
    runs-on: ubuntu-20.04
    steps:
      - uses: actions/checkout@v2
      - uses: actions/setup-go@v2
        with:
          go-version: 1.18
      - run: go install github.com/hidetatz/gcp_tf_iam_binding_validator/cmd/gcp_tf_iam_binding_validator@latest
      - run: |
          gcp_tf_iam_binding_validator -dir your_terraform_directory
```

If there are duplications, it will be shown in the standard output then the process exits with 1. Otherwise 0.

## Usage

Pass the directory which contains terraform (.tf) files.

```shell
gcp_tf_iam_binding_validator -dir ./test/1
```

Note that gcp_tf_iam_binding_validator does not check the GCP project in the google_project_iam_binding definition. This means you should make sure every terraform files in your passing directory are for the same GCP project.

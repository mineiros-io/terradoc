section {
  content = <<END
[<img src="https://raw.githubusercontent.com/mineiros-io/brand/3bffd30e8bdbbde32c143e2650b2faa55f1df3ea/mineiros-primary-logo.svg" width="400"/>][homepage]

[![Terraform Version][badge-terraform]][releases-terraform]
[![Google Provider Version][badge-tf-gcp]][releases-google-provider]
[![Join Slack][badge-slack]][slack]
END

  section {
    title = "terraform-google-secret-manager-iam"
    content = <<END
A [Terraform](https://www.terraform.io) module to create a [Google Secret Manager IAM](https://cloud.google.com/secret-manager/docs/access-control) on [Google Cloud Services (GCP)](https://cloud.google.com/).

**_This module supports Terraform version 1
and is compatible with the Terraform Google Provider version 3._**

This module is part of our Infrastructure as Code (IaC) framework
that enables our users and customers to easily deploy and manage reusable,
secure, and production-grade cloud infrastructure.

- [Module Features](#module-features)
- [Getting Started](#getting-started)
- [Module Argument Reference](#module-argument-reference)
  - [Top-level Arguments](#top-level-arguments)
    - [Module Configuration](#module-configuration)
    - [Main Resource Configuration](#main-resource-configuration)
    - [Extended Resource Configuration](#extended-resource-configuration)
- [Module Attributes Reference](#module-attributes-reference)
- [External Documentation](#external-documentation)
  - [Google Documentation](#google-documentation)
  - [Terraform Google Provider Documentation](#terraform-google-provider-documentation)
- [Module Versioning](#module-versioning)
  - [Backwards compatibility in `0.0.z` and `0.y.z` version](#backwards-compatibility-in-00z-and-0yz-version)
- [About Mineiros](#about-mineiros)
- [Reporting Issues](#reporting-issues)
- [Contributing](#contributing)
- [Makefile Targets](#makefile-targets)
- [License](#license)
END

    section {
      title = "Module Features"
      content = <<END
This module implements the following terraform resources:

- `google_secret_manager_secret_iam_binding`
- `google_secret_manager_secret_iam_member`
- `google_secret_manager_secret_iam_policy`
END
    }

    section {
      title = "Getting Started"
      content = <<END
Most basic usage just setting required arguments:

```hcl
module "terraform-google-secret-manager-iam" {
  source = "github.com/mineiros-io/terraform-google-secret-manager-iam.git?ref=v0.1.0"

  secret_id = google_secret_manager_secret.secret-basic.secret_id
  role      = "roles/secretmanager.secretAccessor"
  members   = ["user:admin@example.com"]
}
```
END
    }

    section {
      title = "Module Argument Reference"
      content = "See [variables.tf] and [examples/] for details and use-cases."

      section {
        title = "Top-level Arguments"

        section {
          title = "Module Configuration"

          variable "module_enabled" {
            type = bool
            description = "Specifies whether resources in the module will be created."
            default = true
          }

          variable "module_depends_on" {
            type = list(any)
            readme_type = "list(dependencies)"
            description = "A list of dependencies. Any object can be _assigned_ to this list to define a hidden external dependency."
            readme_example = <<END
module_depends_on = [
  google_network.network
]
END
          }
        }

        section {
          title = "Main Resource Configuration"

          variable "secret_id" {
            required = true
            type = string
            description = "The id of the secret."
          }

          variable "members" {
            type = set(string)
            default = []
            description = <<END
Identities that will be granted the privilege in role. Each entry can have one of the following values:
  - `allUsers`: A special identifier that represents anyone who is on the internet; with or without a Google account.
  - `allAuthenticatedUsers`: A special identifier that represents anyone who is authenticated with a Google account or a service account.
  - `user:{emailid}`: An email address that represents a specific Google account. For example, alice@gmail.com or joe@example.com.
  - `serviceAccount:{emailid}`: An email address that represents a service account. For example, my-other-app@appspot.gserviceaccount.com.
  - `group:{emailid}`: An email address that represents a Google group. For example, admins@example.com.
  - `domain:{domain}`: A G Suite domain (primary, instead of alias) name that represents all the users of that domain. For example, google.com or example.com.
END
          }

          variable "role" {
            type = string
            description = "The role that should be applied. Note that custom roles must be of the format `[projects|organizations]/{parent-name}/roles/{role-name}`."
          }

          variable "project" {
            type = string
            description = "The resource name of the project the policy is attached to. Its format is `projects/{project_id}`."
          }

          variable "authoritative" {
            description = "Whether to exclusively set (authoritative mode) or add (non-authoritative/additive mode) members to the role."
            type = bool
            default = true
          }

          variable "policy_bindings" {
            type = list(any)
            readme_type = "list(policy_bindings)"
            description = "A list of IAM policy bindings."
            readme_example = <<END
policy_bindings = [{
  role    = "roles/secretmanager.secretAccessor"
  members = ["user:member@example.com"]
}]
END

            attribute "role" {
              description = "The role that should be applied."
              required = true
              type = string
            }

            attribute "members" {
              required = true
              type = string
              default = "var.members"
              description = "Identities that will be granted the privilege in `role`."
            }

            attribute "condition" {
              type = any
              readme_type = "object(condition)"
              description = "An IAM Condition for a given binding."
              readme_example = <<END
condition = {
  expression = "request.time < timestamp(\"2022-01-01T00:00:00Z\")"
  title      = "expires_after_2021_12_31"
}
END

              attribute "expression" {
                type = string
                required = true
                description = "Textual representation of an expression in Common Expression Language syntax."
              }

              attribute "title" {
                type = string
                required = true
                description = "A title for the expression, i.e. a short string describing its purpose."
              }

              attribute "description" {
                type = string
                description = "An optional description of the expression. This is a longer text which describes the expression, e.g. when hovered over it in a UI."
              }
            }
          }
        }
      }
    }
  }

  section {
    title = "External Documentation"

    section {
      title = "Google Documentation"
      content = <<END
- Secret Manager: <https://cloud.google.com/secret-manager/docs>
- Secret Manager Access Control: <https://cloud.google.com/secret-manager/docs/access-control>
END

    }

    section {
      title = "Terraform Google Provider Documentation"
      content = <<END
- <https://registry.terraform.io/providers/hashicorp/google/latest/docs/resources/secret_manager_secret>
- <https://registry.terraform.io/providers/hashicorp/google/latest/docs/resources/secret_manager_secret_iam>
END
    }
  }

  section {
    title = "Module Versioning"
    content = <<END
This Module follows the principles of [Semantic Versioning (SemVer)].

Given a version number `MAJOR.MINOR.PATCH`, we increment the:

1. `MAJOR` version when we make incompatible changes,
2. `MINOR` version when we add functionality in a backwards compatible manner, and
3. `PATCH` version when we make backwards compatible bug fixes.
END

    section {
      title = "Backwards compatibility in `0.0.z` and `0.y.z` version"
      content = <<END
- Backwards compatibility in versions `0.0.z` is **not guaranteed** when `z` is increased. (Initial development)
- Backwards compatibility in versions `0.y.z` is **not guaranteed** when `y` is increased. (Pre-release)
END
    }
  }

  section {
    title = "About Mineiros"
    content = <<END
[Mineiros][homepage] is a remote-first company headquartered in Berlin, Germany
that solves development, automation and security challenges in cloud infrastructure.

Our vision is to massively reduce time and overhead for teams to manage and
deploy production-grade and secure cloud infrastructure.

We offer commercial support for all of our modules and encourage you to reach out
if you have any questions or need help. Feel free to email us at [hello@mineiros.io] or join our
[Community Slack channel][slack].
END
  }

  section {
    title = "Reporting Issues"
    content = "We use GitHub [Issues] to track community reported issues and missing features."
  }

  section {
    title = "Contributing"
    content = <<END
Contributions are always encouraged and welcome! For the process of accepting changes, we use
[Pull Requests]. If you'd like more information, please see our [Contribution Guidelines].
END
  }

  section {
    title = "Makefile Targets"
    content = <<END
This repository comes with a handy [Makefile].
Run `make help` to see details on each available target.
END
  }

  section {
    title = "License"
    content = <<END
[![license][badge-license]][apache20]

This module is licensed under the Apache License Version 2.0, January 2004.
Please see [LICENSE] for full details.

Copyright &copy; 2020-2021 [Mineiros GmbH][homepage]
END
  }
}

references {
  ref "homepage" {
    value = "https://mineiros.io/?ref=terraform-google-secret-manager-iam"
  }

  ref "hello@mineiros.io" {
    value = "mailto:hello@mineiros.io"
  }

  ref "badge-build" {
    value = "https://github.com/mineiros-io/terraform-google-secret-manager-iam/workflows/Tests/badge.svg"
  }

  ref "badge-semver" {
    value = "https://img.shields.io/github/v/tag/mineiros-io/terraform-google-secret-manager-iam.svg?label=latest&sort=semver"
  }

  ref "badge-license" {
    value = "https://img.shields.io/badge/license-Apache%202.0-brightgreen.svg"
  }

  ref "badge-terraform" {
    value = "https://img.shields.io/badge/Terraform-1.x-623CE4.svg?logo=terraform"
  }

  ref "badge-slack" {
    value = "https://img.shields.io/badge/slack-@mineiros--community-f32752.svg?logo=slack"
  }

  ref "build-status" {
    value = "https://github.com/mineiros-io/terraform-google-secret-manager-iam/actions"
  }

  ref "releases-github" {
    value = "https://github.com/mineiros-io/erraform-google-secret-manager-iam/releases"
  }

  ref "releases-terraform" {
    value = "https://github.com/hashicorp/terraform/releases"
  }

  ref "badge-tf-gcp" {
    value = "https://img.shields.io/badge/google-3.x-1A73E8.svg?logo=terraform"
  }

  ref "releases-google-provider" {
    value = "https://github.com/terraform-providers/terraform-provider-google/releases"
  }

  ref "apache20" {
    value = "https://opensource.org/licenses/Apache-2.0"
  }

  ref "slack" {
    value = "https://mineiros.io/slack"
  }

  ref "terraform" {
    value = "https://www.terraform.io"
  }

  ref "gcp" {
    value = "https://cloud.google.com/"
  }

  ref "semantic versioning (semver)" {
    value = "https://semver.org/"
  }

  ref "variables.tf" {
    value = "https://github.com/mineiros-io/terraform-google-secret-manager-iam/blob/main/variables.tf"
  }

  ref "examples/" {
    value = "https://github.com/mineiros-io/terraform-google-secret-manager-iam/blob/main/examples"
  }

  ref "issues" {
    value = "https://github.com/mineiros-io/terraform-google-secret-manager-iam/issues"
  }

  ref "license" {
    value = "https://github.com/mineiros-io/terraform-google-secret-manager-iam/blob/main/LICENSE"
  }

  ref "makefile" {
    value = "https://github.com/mineiros-io/terraform-google-secret-manager-iam/blob/main/Makefile"
  }

  ref "pull requests" {
    value = "https://github.com/mineiros-io/terraform-google-secret-manager-iam/pulls"
  }

  ref "contribution guidelines" {
    value = "https://github.com/mineiros-io/terraform-google-secret-manager-iam/blob/main/CONTRIBUTING.md"
  }
}

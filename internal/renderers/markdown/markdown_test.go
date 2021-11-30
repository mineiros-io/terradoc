package markdown_test

import (
	"bytes"
	"strings"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/mineiros-io/terradoc/internal/entities"
	"github.com/mineiros-io/terradoc/internal/renderers/markdown"
	"github.com/mineiros-io/terradoc/internal/types"
	"github.com/mineiros-io/terradoc/test"
)

func TestRender(t *testing.T) {
	definition := entities.Definition{
		Header: entities.Header{
			Image: "https://raw.githubusercontent.com/mineiros-io/brand/3bffd30e8bdbbde32c143e2650b2faa55f1df3ea/mineiros-primary-logo.svg",
			URL:   "https://mineiros.io/?ref=terraform-google-secret-manager-iam",
			Badges: []entities.Badge{
				{
					Image: "https://img.shields.io/badge/Terraform-1.x-623CE4.svg?logo=terraform",
					Text:  "Terraform Version",
					Name:  "badge-terraform",
					URL:   "https://github.com/hashicorp/terraform/releases",
				},
				{
					Image: "https://img.shields.io/badge/google-3.x-1A73E8.svg?logo=terraform",
					Text:  "Google Provider Version",
					Name:  "badge-tf-gcp",
					URL:   "https://github.com/terraform-providers/terraform-provider-google/releases",
				},
				{
					Image: "https://img.shields.io/badge/slack-@mineiros--community-f32752.svg?logo=slack",
					Text:  "Join Slack",
					Name:  "badge-slack",
					URL:   "https://mineiros.io/slack",
				},
			},
		},
		Sections: []entities.Section{
			{
				Level: 1,
				Title: "terraform-google-secret-manager-iam",
				Content: `A [Terraform](https://www.terraform.io) module to create a [Google Secret Manager IAM](https://cloud.google.com/secret-manager/docs/access-control) on [Google Cloud Services (GCP)](https://cloud.google.com/).

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
  - [Backwards compatibility in ` + "`0.0.z` and `0.y.z`" + ` version](#backwards-compatibility-in-00z-and-0yz-version)
- [About Mineiros](#about-mineiros)
- [Reporting Issues](#reporting-issues)
- [Contributing](#contributing)
- [Makefile Targets](#makefile-targets)
- [License](#license)`,
				SubSections: []entities.Section{
					{
						Level: 2,
						Title: "Module Features",
						Content: `This module implements the following terraform resources:

- ` + "`google_secret_manager_secret_iam_binding`" + `
- ` + "`google_secret_manager_secret_iam_member`" + `
- ` + "`google_secret_manager_secret_iam_policy`",
					},
					{
						Level: 2,
						Title: "Getting Started",
						Content: `Most basic usage just setting required arguments:

` + "```hcl" + `
module "terraform-google-secret-manager-iam" {
  source = "github.com/mineiros-io/terraform-google-secret-manager-iam.git?ref=v0.1.0"

  secret_id = google_secret_manager_secret.secret-basic.secret_id
  role      = "roles/secretmanager.secretAccessor"
  members   = ["user:admin@example.com"]
}
` + "```",
					},
					{
						Level:   2,
						Title:   "Module Argument Reference",
						Content: "See [variables.tf] and [examples/] for details and use-cases.",
						SubSections: []entities.Section{
							{
								Level: 3,
								Title: "Top-level Arguments",
								SubSections: []entities.Section{
									{
										Title: "Module Configuration",
										Level: 4,
										Variables: []entities.Variable{
											{
												Name: "module_enabled",
												Type: entities.Type{
													TerraformType: entities.TerraformType{
														Type: types.TerraformBool,
													},
												},
												Description: "Specifies whether resources in the module will be created.",
												Default:     []byte("true"),
											},
											{
												Name: "module_depends_on",
												Type: entities.Type{
													ReadmeType: "list(dependencies)",
													TerraformType: entities.TerraformType{
														Type:       types.TerraformList,
														NestedType: types.TerraformAny,
													},
												},
												Description: "A list of dependencies. Any object can be _assigned_ to this list to define a hidden external dependency.",
												ReadmeExample: `module_depends_on = [
  google_network.network
]`,
											},
										},
									},
									{
										Level: 4,
										Title: "Main Resource Configuration",
										Variables: []entities.Variable{
											{
												Name:     "secret_id",
												Required: true,
												Type: entities.Type{
													TerraformType: entities.TerraformType{
														Type: types.TerraformString,
													},
												},
												Description: "The id of the secret.",
											},
											{
												Name: "members",
												Type: entities.Type{
													TerraformType: entities.TerraformType{
														Type:       types.TerraformSet,
														NestedType: types.TerraformString,
													},
												},
												Default: []byte("[]"),
												Description: `Identities that will be granted the privilege in role. Each entry can have one of the following values:
  - ` + "`allUsers`" + `: A special identifier that represents anyone who is on the internet; with or without a Google account.
  - ` + "`allAuthenticatedUsers`" + `: A special identifier that represents anyone who is authenticated with a Google account or a service account.
  - ` + "`user:{emailid}`" + `: An email address that represents a specific Google account. For example, alice@gmail.com or joe@example.com.
  - ` + "`serviceAccount:{emailid}`" + `: An email address that represents a service account. For example, my-other-app@appspot.gserviceaccount.com.
  - ` + "`group:{emailid}`" + `: An email address that represents a Google group. For example, admins@example.com.
  - ` + "`domain:{domain}`" + `: A G Suite domain (primary, instead of alias) name that represents all the users of that domain. For example, google.com or example.com.
  - ` + "`projectOwner:projectid`" + `: Owners of the given project. For example, ` + "`projectOwner:my-example-project`" + `
  - ` + "`projectEditor:projectid`" + `: Editors of the given project. For example, ` + "`projectEditor:my-example-project`" + `
  - ` + "`projectViewer:projectid`" + `: Viewers of the given project. For example, ` + "`projectViewer:my-example-project`",
											},
											{
												Name: "role",
												Type: entities.Type{
													TerraformType: entities.TerraformType{
														Type: types.TerraformString,
													},
												},
												Description: "The role that should be applied. Note that custom roles must be of the format `[projects|organizations]/{parent-name}/roles/{role-name}`.",
											},
											{
												Name: "project",
												Type: entities.Type{
													TerraformType: entities.TerraformType{
														Type: types.TerraformString,
													},
												},
												Description: "The ID of the project in which the resource belongs. If it is not provided, the project will be parsed from the identifier of the parent resource. If no project is provided in the parent identifier and no project is specified, the provider project is used.",
											},
											{
												Name: "authoritative",
												Type: entities.Type{
													TerraformType: entities.TerraformType{
														Type: types.TerraformBool,
													},
												},
												Default:     []byte("true"),
												Description: "Whether to exclusively set (authoritative mode) or add (non-authoritative/additive mode) members to the role.",
											},
											{
												Name: "policy_bindings",
												Type: entities.Type{
													ReadmeType: "list(policy_bindings)",
													TerraformType: entities.TerraformType{
														Type:       types.TerraformList,
														NestedType: types.TerraformAny,
													},
												},
												Description: "A list of IAM policy bindings.",
												ReadmeExample: `policy_bindings = [{
  role    = "roles/secretmanager.secretAccessor"
  members = ["user:member@example.com"]
}]`,
												Attributes: []entities.Attribute{
													{
														Level:       1,
														Name:        "role",
														Description: "The role that should be applied.",
														Required:    true,
														Type: entities.Type{
															TerraformType: entities.TerraformType{
																Type: types.TerraformString,
															},
														},
													},
													{
														Level:       1,
														Name:        "members",
														Description: "Identities that will be granted the privilege in `role`.",
														Type: entities.Type{
															TerraformType: entities.TerraformType{
																Type:       types.TerraformSet,
																NestedType: types.TerraformString,
															},
														},
														Default: []byte("var.members"),
													},
													{
														Level:       1,
														Name:        "condition",
														Description: "An IAM Condition for a given binding.",
														Type: entities.Type{
															ReadmeType: "object(condition)",
															TerraformType: entities.TerraformType{
																Type: types.TerraformAny,
															},
														},
														ReadmeExample: `condition = {
  expression = "request.time < timestamp(\"2022-01-01T00:00:00Z\")"
  title      = "expires_after_2021_12_31"
}`,
														Attributes: []entities.Attribute{
															{
																Level:       2,
																Name:        "expression",
																Description: "Textual representation of an expression in Common Expression Language syntax.",
																Required:    true,
																Type: entities.Type{
																	TerraformType: entities.TerraformType{
																		Type: types.TerraformString,
																	},
																},
															},
															{
																Level:       2,
																Name:        "title",
																Description: "A title for the expression, i.e. a short string describing its purpose.",
																Required:    true,
																Type: entities.Type{
																	TerraformType: entities.TerraformType{
																		Type: types.TerraformString,
																	},
																},
															},
															{
																Level:       2,
																Name:        "description",
																Description: "An optional description of the expression. This is a longer text which describes the expression, e.g. when hovered over it in a UI.",
																Type: entities.Type{
																	TerraformType: entities.TerraformType{
																		Type: types.TerraformString,
																	},
																},
															},
														},
													},
												},
											},
										},
									},
									{
										Level: 4,
										Title: "Extended Resource Configuration",
									},
								},
							},
						},
					},
				},
			},
			{
				Level: 2,
				Title: "Module Attributes Reference",
				Content: `The following attributes are exported in the outputs of the module:

- **` + "`module_enabled`" + `**

  Whether this module is enabled.

- **` + "`iam`" + `**

  All attributes of the created ` + "`iam_binding` or `iam_member` or `iam_policy`" + ` resource according to the mode.`,
			},
			{
				Level: 2,
				Title: "External Documentation",
				SubSections: []entities.Section{
					{
						Level: 3,
						Title: "Google Documentation",
						Content: `- Secret Manager: <https://cloud.google.com/secret-manager/docs>
- Secret Manager Access Control: <https://cloud.google.com/secret-manager/docs/access-control>`,
					},
					{
						Level: 3,
						Title: "Terraform Google Provider Documentation",
						Content: `- <https://registry.terraform.io/providers/hashicorp/google/latest/docs/resources/secret_manager_secret>
- <https://registry.terraform.io/providers/hashicorp/google/latest/docs/resources/secret_manager_secret_iam>`,
					},
				},
			},
			{
				Level: 2,
				Title: "Module Versioning",
				Content: `This Module follows the principles of [Semantic Versioning (SemVer)].

Given a version number ` + "`MAJOR.MINOR.PATCH`" + `, we increment the:

1. ` + "`MAJOR`" + ` version when we make incompatible changes,
2. ` + "`MINOR`" + ` version when we add functionality in a backwards compatible manner, and
3. ` + "`PATCH`" + ` version when we make backwards compatible bug fixes.`,
				SubSections: []entities.Section{
					{
						Level: 3,
						Title: "Backwards compatibility in `0.0.z` and `0.y.z` version",
						Content: `- Backwards compatibility in versions ` + "`0.0.z`" + ` is **not guaranteed** when ` + "`z`" + ` is increased. (Initial development)
- Backwards compatibility in versions ` + "`0.y.z`" + ` is **not guaranteed** when ` + "`y`" + ` is increased. (Pre-release)`,
					},
				},
			},
			{
				Level: 2,
				Title: "About Mineiros",
				Content: `[Mineiros][homepage] is a remote-first company headquartered in Berlin, Germany
that solves development, automation and security challenges in cloud infrastructure.

Our vision is to massively reduce time and overhead for teams to manage and
deploy production-grade and secure cloud infrastructure.

We offer commercial support for all of our modules and encourage you to reach out
if you have any questions or need help. Feel free to email us at [hello@mineiros.io] or join our
[Community Slack channel][slack].`,
			},
			{
				Level:   2,
				Title:   "Reporting Issues",
				Content: "We use GitHub [Issues] to track community reported issues and missing features.",
			},
			{
				Level: 2,
				Title: "Contributing",
				Content: `Contributions are always encouraged and welcome! For the process of accepting changes, we use
[Pull Requests]. If you'd like more information, please see our [Contribution Guidelines].`,
			},
			{
				Level: 2,
				Title: "Makefile Targets",
				Content: `This repository comes with a handy [Makefile].
Run ` + "`make help`" + ` to see details on each available target.`,
			},
			{
				Level: 2,
				Title: "License",
				Content: `[![license][badge-license]][apache20]

This module is licensed under the Apache License Version 2.0, January 2004.
Please see [LICENSE] for full details.

Copyright &copy; 2020-2021 [Mineiros GmbH][homepage]`,
			},
		},
		References: []entities.Reference{
			{
				Name:  "homepage",
				Value: "https://mineiros.io/?ref=terraform-google-secret-manager-iam",
			},
			{
				Name:  "hello@mineiros.io",
				Value: "mailto:hello@mineiros.io",
			},
			{
				Name:  "badge-build",
				Value: "https://github.com/mineiros-io/terraform-google-secret-manager-iam/workflows/Tests/badge.svg",
			},
			{
				Name:  "badge-semver",
				Value: "https://img.shields.io/github/v/tag/mineiros-io/terraform-google-secret-manager-iam.svg?label=latest&sort=semver",
			},
			{
				Name:  "badge-license",
				Value: "https://img.shields.io/badge/license-Apache%202.0-brightgreen.svg",
			},
			{
				Name:  "badge-terraform",
				Value: "https://img.shields.io/badge/Terraform-1.x-623CE4.svg?logo=terraform",
			},
			{
				Name:  "badge-slack",
				Value: "https://img.shields.io/badge/slack-@mineiros--community-f32752.svg?logo=slack",
			},
			{
				Name:  "build-status",
				Value: "https://github.com/mineiros-io/terraform-google-secret-manager-iam/actions",
			},
			{
				Name:  "releases-github",
				Value: "https://github.com/mineiros-io/erraform-google-secret-manager-iam/releases",
			},
			{
				Name:  "releases-terraform",
				Value: "https://github.com/hashicorp/terraform/releases",
			},
			{
				Name:  "badge-tf-gcp",
				Value: "https://img.shields.io/badge/google-3.x-1A73E8.svg?logo=terraform",
			},
			{
				Name:  "releases-google-provider",
				Value: "https://github.com/terraform-providers/terraform-provider-google/releases",
			},
			{
				Name:  "apache20",
				Value: "https://opensource.org/licenses/Apache-2.0",
			},
			{
				Name:  "slack",
				Value: "https://mineiros.io/slack",
			},
			{
				Name:  "terraform",
				Value: "https://www.terraform.io",
			},
			{
				Name:  "gcp",
				Value: "https://cloud.google.com/",
			},
			{
				Name:  "semantic versioning (semver)",
				Value: "https://semver.org/",
			},
			{
				Name:  "variables.tf",
				Value: "https://github.com/mineiros-io/terraform-google-secret-manager-iam/blob/main/variables.tf",
			},
			{
				Name:  "examples/",
				Value: "https://github.com/mineiros-io/terraform-google-secret-manager-iam/blob/main/examples",
			},
			{
				Name:  "issues",
				Value: "https://github.com/mineiros-io/terraform-google-secret-manager-iam/issues",
			},
			{
				Name:  "license",
				Value: "https://github.com/mineiros-io/terraform-google-secret-manager-iam/blob/main/LICENSE",
			},
			{
				Name:  "makefile",
				Value: "https://github.com/mineiros-io/terraform-google-secret-manager-iam/blob/main/Makefile",
			},
			{
				Name:  "pull requests",
				Value: "https://github.com/mineiros-io/terraform-google-secret-manager-iam/pulls",
			},
			{
				Name:  "contribution guidelines",
				Value: "https://github.com/mineiros-io/terraform-google-secret-manager-iam/blob/main/CONTRIBUTING.md",
			},
		},
	}

	buf := new(bytes.Buffer)
	err := markdown.Render(buf, definition)
	if err != nil {
		t.Errorf("Expected no error but got %q instead", err)
	}

	got := strings.TrimSpace(buf.String())

	wantContent := test.ReadFixture(t, "golden-readme.md")
	want := string(bytes.TrimSpace(wantContent))

	if diff := cmp.Diff(got, want); diff != "" {
		t.Logf("\n\nWANT:\n%q\n\nGOT:\n%q\n", want, got)
		t.Errorf("Expected golden file to match result (-want +got):\n%s", diff)
	}
}

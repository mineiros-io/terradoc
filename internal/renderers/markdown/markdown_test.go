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
		Sections: []entities.Section{
			{
				Title:       "Module Argument Reference",
				Description: "See [variables.tf] and [examples/] for details and use-cases.",
				Level:       1,
				SubSections: []entities.Section{
					{
						Title: "Top-level Arguments",
						Level: 2,
						SubSections: []entities.Section{
							{
								Title: "Module Configuration",
								Level: 3,
								Variables: []entities.Variable{
									{
										Name: "module_enabled",
										Type: entities.Type{
											TerraformType: entities.TerraformType{Type: types.TerraformBool},
										},
										Description: "Specifies whether resources in the module will be created.",
										Default:     []byte("true"),
									},
									{
										Name: "module_depends_on",
										Type: entities.Type{
											TerraformType: entities.TerraformType{Type: types.TerraformList},
											ReadmeType:    "list(dependencies)",
										},
										Description: "A list of dependencies. Any object can be _assigned_ to this list to define a hidden external dependency.",
										// BUG: this is not correct - google_network.network should not be a string
										// but it fails right now
										ReadmeExample: `module_depends_on = ["google_network.network"]`,
									},
								},
							},
						},
					},
					{
						Level: 3,
						Title: "Main Resource Configuration",
						Variables: []entities.Variable{
							{
								Name: "secret_id",
								Type: entities.Type{
									TerraformType: entities.TerraformType{Type: types.TerraformString},
								},
								Description: "The id of the secret.",
								Required:    true,
							},
							{
								Name: "members",
								Type: entities.Type{
									TerraformType: entities.TerraformType{Type: types.TerraformSet, NestedType: types.TerraformString},
								},
								Description: "Identities that will be granted the privilege in role. Each entry can have one of the following values:\n  - `allUsers`: A special identifier that represents anyone who is on the internet; with or without a Google account.\n  - `allAuthenticatedUsers`: A special identifier that represents anyone who is authenticated with a Google account or a service account.\n  - `user:{emailid}`: An email address that represents a specific Google account. For example, alice@gmail.com or joe@example.com.\n  - `serviceAccount:{emailid}`: An email address that represents a service account. For example, my-other-app@appspot.gserviceaccount.com.\n  - `group:{emailid}`: An email address that represents a Google group. For example, admins@example.com.\n  - `domain:{domain}`: A G Suite domain (primary, instead of alias) name that represents all the users of that domain. For example, google.com or example.com.\n",
								// Description: fmt.Sprintf(`Identities that will be granted the privilege in role. Each entry can have one of the following values:
								// - %s: A special identifier that represents anyone who is on the internet; with or without a Google account.
								// - %s: A special identifier that represents anyone who is authenticated with a Google account or a service account.
								// - %s: An email address that represents a specific Google account. For example, alice@gmail.com or joe@example.com.
								// - %s: An email address that represents a service account. For example, my-other-app@appspot.gserviceaccount.com.
								// - %s: An email address that represents a Google group. For example, admins@example.com.
								// - %s: A G Suite domain (primary, instead of alias) name that represents all the users of that domain. For example, google.com or example.com.`, "`allUsers`", "`allAuthenticatedUsers`", "`user:{emailid}`", "`serviceAccount:{emailid}`", "`group:{emailid}`", "`domain:{domain}`"),
								Default: []byte("[]"),
							},
							{
								Name: "role",
								Type: entities.Type{
									TerraformType: entities.TerraformType{Type: types.TerraformString},
								},
								Description: "The role that should be applied. Note that custom roles must be of the format `[projects|organizations]/{parent-name}/roles/{role-name}`.",
							},
							{
								Name: "project",
								Type: entities.Type{
									TerraformType: entities.TerraformType{Type: types.TerraformString},
								},
								Description: "The resource name of the project the policy is attached to. Its format is `projects/{project_id}`.",
							},
							{
								Name: "authoritative",
								Type: entities.Type{
									TerraformType: entities.TerraformType{Type: types.TerraformBool},
								},
								Description: "Whether to exclusively set (authoritative mode) or add (non-authoritative/additive mode) members to the role.",
								Default:     []byte("true"),
							},
							{
								Name: "policy_bindings",
								Type: entities.Type{
									TerraformType: entities.TerraformType{
										Type:       types.TerraformList,
										NestedType: types.TerraformAny,
									},
									ReadmeType: "list(policy_bindings)",
								},
								Description: "A list of IAM policy bindings.",
								ReadmeExample: `policy_bindings = [{
    members = ["user:member@example.com"]
    role    = "roles/secretmanager.secretAccessor"
  }]`,
								Attributes: []entities.Attribute{
									{
										Level:    1,
										Name:     "role",
										Required: true,
										Type: entities.Type{
											TerraformType: entities.TerraformType{Type: types.TerraformString},
										},
										Description: "The role that should be applied.",
									},
									{
										Level:    1,
										Name:     "members",
										Required: true,
										Type: entities.Type{
											TerraformType: entities.TerraformType{Type: types.TerraformString},
										},
										Description: "Identities that will be granted the privilege in `role`.",
										// BUG: this is not correct - var.members should not be a string
										// but it fails otherwise
										Default: []byte(`"var.members"`),
									},
									{
										Level: 1,
										Name:  "condition",
										Type: entities.Type{
											TerraformType: entities.TerraformType{Type: types.TerraformAny},
											ReadmeType:    "object(condition)",
										},
										Description: "An IAM Condition for a given binding.",
										ReadmeExample: `condition = {
    expression = "request.time < timestamp(\"2022-01-01T00:00:00Z\")"
    title      = "expires_after_2021_12_31"
  }`,
										Attributes: []entities.Attribute{
											{
												Level:    2,
												Name:     "expression",
												Required: true,
												Type: entities.Type{
													TerraformType: entities.TerraformType{Type: types.TerraformString},
												},
												Description: "Textual representation of an expression in Common Expression Language syntax.",
											},
											{
												Level:    2,
												Name:     "title",
												Required: true,
												Type: entities.Type{
													TerraformType: entities.TerraformType{Type: types.TerraformString},
												},
												Description: "A title for the expression, i.e. a short string describing its purpose.",
											},
											{
												Level: 2,
												Name:  "description",
												Type: entities.Type{
													TerraformType: entities.TerraformType{Type: types.TerraformString},
												},
												Description: "An optional description of the expression. This is a longer text which describes the expression, e.g. when hovered over it in a UI.",
											},
										},
									},
								},
							},
						},
					},
				},
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
		t.Errorf("Expected golden file to match result (-want +got):\n%s", diff)
	}
}

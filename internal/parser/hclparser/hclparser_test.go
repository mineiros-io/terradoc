package hclparser_test

import (
	"bytes"
	"encoding/json"
	"strings"
	"testing"

	"github.com/mineiros-io/terradoc/internal/entities"
	"github.com/mineiros-io/terradoc/internal/parser/hclparser"
	"github.com/mineiros-io/terradoc/internal/types"
	"github.com/mineiros-io/terradoc/test"
)

func TestParse(t *testing.T) {
	for _, tt := range []struct {
		desc      string
		inputFile string
		want      entities.Definition
	}{
		{
			desc:      "with a valid input",
			inputFile: "input.tfdoc.hcl",
			want: entities.Definition{
				Header: entities.Header{
					Image: "https://raw.githubusercontent.com/mineiros-io/brand/3bffd30e8bdbbde32c143e2650b2faa55f1df3ea/mineiros-primary-logo.svg",
					URL:   "https://www.mineiros.io",
				},
				Sections: []entities.Section{
					{
						Level: 1,
						Title: "root section",
						Content: `This is the root section content.

Section contents support anything markdown and allow us to make references like this one: [mineiros-website]`,
						SubSections: []entities.Section{
							{
								Level: 2,
								Title: "sections with variables",

								SubSections: []entities.Section{
									{
										Level: 3,
										Title: "example",
										Variables: []entities.Variable{
											{
												Name: "name",
												Type: entities.Type{
													TerraformType: entities.TerraformType{
														Type: types.TerraformString,
													},
												},
												Description: "describes the name of the last person who bothered to change this file",
												Default:     json.RawMessage("nathan"),
											},
										},
									},
									{
										Level:   3,
										Title:   "section of beers",
										Content: "an excuse to mention alcohol",
										Variables: []entities.Variable{
											{
												Name: "beers",
												Type: entities.Type{
													TerraformType: entities.TerraformType{
														Type:       types.TerraformList,
														NestedType: types.TerraformAny,
													},
													ReadmeType: "list(beer)",
												},
												Description:      "a list of beers",
												Default:          json.RawMessage("[]"),
												Required:         true,
												ForcesRecreation: true,
												ReadmeExample:    "",
												Attributes: []entities.Attribute{
													{
														Name: "name",
														Type: entities.Type{
															TerraformType: entities.TerraformType{
																Type: types.TerraformString,
															},
														},
														Description:      "the name of the beer",
														ForcesRecreation: false,
													},
													{
														Name: "type",
														Type: entities.Type{
															TerraformType: entities.TerraformType{
																Type: types.TerraformString,
															},
														},
														Description:      "the type of the beer",
														ForcesRecreation: true,
													},
													{
														Name: "abv",
														Type: entities.Type{
															TerraformType: entities.TerraformType{
																Type: types.TerraformNumber,
															},
														},
														Description:      "beer's alcohol by volume content",
														ForcesRecreation: true,
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
			},
		},
	} {
		t.Run(tt.desc, func(t *testing.T) {
			r := test.OpenFixture(t, tt.inputFile)
			// parsed definition
			definition, err := hclparser.Parse(r, "foo")
			if err != nil {
				t.Fatal(err)
			}

			assertEqualDefinitions(t, tt.want, definition) //
		})
	}
}

func TestParseInvalidContent(t *testing.T) {
	for _, tt := range []struct {
		desc                 string
		wantErrorMsgContains string
		content              string
	}{
		{
			desc:                 "variable block without a label",
			wantErrorMsgContains: "Missing name for variable; All variable blocks must have 1 labels (name).",
			content: `
section {
  title = "test"

  section {
    title = "test"

    variable {
      type        = string
    }
  }
}

`,
		},
		{
			desc:                 "variable block without a type",
			wantErrorMsgContains: "Missing required argument; The argument \"type\" is required, but no definition was found.",
			content: `
section {
  title = "test"

  section {
    title = "test"

    variable "foo" {
      default = []
    }
  }
}

`,
		},
		{
			desc:                 "attribute block without a label",
			wantErrorMsgContains: "Missing name for attribute; All attribute blocks must have 1 labels (name)",
			content: `
section {
  title = "test"

  section {
    title = "test"

    variable "test" {
      type        = string

      attribute {
        type = number
      }
    }
  }
}
`,
		},
		{
			desc:                 "attribute block without a type",
			wantErrorMsgContains: "Missing required argument; The argument \"type\" is required, but no definition was found",
			content: `
section {
  title = "test"

  section {
    title = "test"

    variable "test" {
      type = string

      attribute "bar" {
        default = number
      }
    }
  }
}
`,
		},
	} {
		t.Run(tt.desc, func(t *testing.T) {
			r := bytes.NewBufferString(tt.content)

			_, err := hclparser.Parse(r, "foo-file")
			if err == nil {
				t.Fatal("Expected error but none occurred")
			}

			if !strings.Contains(err.Error(), tt.wantErrorMsgContains) {
				t.Errorf("Expected error message to contain %q but got %q instead", tt.wantErrorMsgContains, err.Error())
			}
		})
	}
}

func assertEqualDefinitions(t *testing.T, want, got entities.Definition) {
	t.Helper()

	assertEqualHeader(t, want.Header, got.Header)
	assertEqualSections(t, want.Sections, got.Sections)
}

func assertEqualHeader(t *testing.T, want, got entities.Header) {
	t.Helper()

	if want.Image != got.Image {
		t.Errorf("Expected header image to be %q but got %q instead", want.Image, got.Image)
	}

	if want.URL != got.URL {
		t.Errorf("Expected header url to be %q but got %q instead", want.URL, got.URL)
	}

	assertEqualBadges(t, want.Badges, got.Badges)

}

func assertEqualBadges(t *testing.T, got, want []entities.Badge) {
	t.Helper()

	if len(want) != len(got) {
		t.Errorf("Expected header url to have %d badges but got %d instead", len(want), len(got))
	}

	// TODO: assert that badges are equivalents
}

func assertEqualSections(t *testing.T, want, got []entities.Section) {
	t.Helper()

	if len(want) != len(got) {
		t.Fatalf("Expected definition to contain %d sections. Found %d instead", len(want), len(got))
	}

	if len(got) == 0 {
		return
	}

	for _, section := range want {
		assertContainsSection(t, got, section)
	}
}

func assertContainsSection(t *testing.T, sectionsList []entities.Section, want entities.Section) {
	t.Helper()

	var found bool
	for _, section := range sectionsList {
		if section.Title == want.Title {
			found = true

			assertSectionEquals(t, want, section)
		}
	}

	if !found {
		t.Errorf("Expected sections list to contain section with title %q. Found none instead", want.Title)
	}
}

func assertSectionEquals(t *testing.T, want, got entities.Section) {
	t.Helper()

	// redundant since we're finding the section by title
	if want.Title != got.Title {
		t.Errorf("Expected section title to be %q. Got %q instead", want.Title, got.Title)
	}

	if want.Content != got.Content {
		t.Errorf("Expected section content to be %q. Got %q instead", want.Content, got.Content)
	}

	if want.Level != got.Level {
		t.Errorf("Expected section level to be %d. Got %d instead", want.Level, got.Level)
	}

	assertEqualVariables(t, want.Variables, got.Variables)
	assertEqualSections(t, want.SubSections, got.SubSections)
}

func assertEqualVariables(t *testing.T, want, got []entities.Variable) {
	t.Helper()

	if len(want) != len(got) {
		t.Fatalf("Expected definition to contain %d variables. Found %d instead", len(want), len(got))
	}

	if len(got) == 0 {
		return
	}

	for _, variable := range want {
		assertContainsVariable(t, got, variable)
	}
}

func assertContainsVariable(t *testing.T, variablesList []entities.Variable, want entities.Variable) {
	t.Helper()

	var found bool
	for _, variable := range variablesList {
		if variable.Name == want.Name {
			found = true

			assertVariableEquals(t, want, variable)
		}
	}

	if !found {
		t.Errorf("Expected variables list to contain %q but didn't find one", want.Name)
	}
}

func assertVariableEquals(t *testing.T, want, got entities.Variable) {
	t.Helper()

	// redundant since we're finding the variable by name
	if want.Name != got.Name {
		t.Errorf("Expected variable name to be %q. Got %q instead", want.Name, got.Name)
	}

	if want.Description != got.Description {
		t.Errorf("Expected variable description to be %q. Got %q instead", want.Description, got.Description)
	}

	if want.Type.TerraformType != got.Type.TerraformType {
		t.Errorf("Expected variable type to be %q. Got %q instead", want.Type.TerraformType, got.Type.TerraformType)
	}

	assertEqualAttributes(t, want.Attributes, got.Attributes)
}

func assertEqualAttributes(t *testing.T, want, got []entities.Attribute) {
	t.Helper()

	if len(want) != len(got) {
		t.Fatalf("Expected definition to contain %d attributes. Found %d instead", len(want), len(got))
	}

	if len(got) == 0 {
		return
	}

	for _, attribute := range want {
		assertContainsAttribute(t, got, attribute)
	}
}

func assertContainsAttribute(t *testing.T, attributesList []entities.Attribute, want entities.Attribute) {
	t.Helper()

	var found bool
	for _, attribute := range attributesList {
		if attribute.Name == want.Name {
			found = true

			assertAttributeEquals(t, want, attribute)
		}
	}

	if !found {
		t.Errorf("Expected attributes list to contain %q but didn't find one", want.Name)
	}
}

func assertAttributeEquals(t *testing.T, want, got entities.Attribute) {
	t.Helper()

	// redundant since we're finding the attribute by name
	if want.Name != got.Name {
		t.Errorf("Expected attribute name to be %q. Got %q instead", want.Name, got.Name)
	}

	if want.Description != got.Description {
		t.Errorf("Expected attribute description to be %q. Got %q instead", want.Description, got.Description)
	}

	if want.Type.TerraformType != got.Type.TerraformType {
		t.Errorf("Expected attribute type to be %q. Got %q instead", want.Type.TerraformType, got.Type.TerraformType)
	}

	assertEqualAttributes(t, want.Attributes, got.Attributes)
}

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
			Image:  "",
			URL:    "",
			Badges: []entities.Badge{},
		},
		Sections: []entities.Section{
			{
				Title:   "Section Title",
				Content: "Section Content",
				Level:   1,
				TOC:     true,
				SubSections: []entities.Section{
					{
						Level: 2,
						Title: "SubSection 1",
						SubSections: []entities.Section{
							{
								Level:   3,
								Title:   "SubSection 1.1",
								Content: "We can add infinite subsections!",
							},
						},
						Variables: []entities.Variable{
							{
								Name: "simple_string",
								Type: entities.Type{
									TerraformType: entities.TerraformType{
										Type: types.TerraformString,
									},
								},
								Description: "A simple string",
							},
						},
					},
					{
						Level: 2,
						Title: "SubSection 2",
						Variables: []entities.Variable{
							{
								Name:    "test_objects",
								Default: []byte("[]"),
								Type: entities.Type{
									ReadmeType: "list(test_object)",
									TerraformType: entities.TerraformType{
										Type:       types.TerraformList,
										NestedType: types.TerraformAny,
									},
								},
								Attributes: []entities.Attribute{
									{
										Level:       1,
										Name:        "name",
										Description: "A string",
										Type: entities.Type{
											TerraformType: entities.TerraformType{
												Type: types.TerraformString,
											},
										},
									},
									{
										Level:       1,
										Name:        "something_complex",
										Description: "Some other object",
										Type: entities.Type{
											TerraformType: entities.TerraformType{
												Type: types.TerraformAny,
											},
										},
										Attributes: []entities.Attribute{
											{
												Level:       2,
												Name:        "nested_string",
												Description: "a nested string",
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
				},
			},
		},
		References: []entities.Reference{},
	}

	buf := new(bytes.Buffer)
	err := markdown.Render(buf, definition)
	if err != nil {
		t.Errorf("Expected no error but got %q instead", err)
	}

	got := strings.TrimSpace(buf.String())

	wantContent := test.ReadFixture(t, "markdown-structure.md")
	want := string(bytes.TrimSpace(wantContent))

	if diff := cmp.Diff(got, want); diff != "" {
		t.Logf("\n\nWANT:\n%q\n\nGOT:\n%q\n", want, got)
		t.Errorf("Expected golden file to match result (-want +got):\n%s", diff)
	}
}

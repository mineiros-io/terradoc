package markdown_test

import (
	"bytes"
	"strings"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/madlambda/spells/assert"
	"github.com/mineiros-io/terradoc/internal/entities"
	"github.com/mineiros-io/terradoc/internal/renderers/markdown"
	"github.com/mineiros-io/terradoc/internal/types"
	"github.com/mineiros-io/terradoc/test"
)

func TestRender(t *testing.T) {
	definition := entities.TFDoc{
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
									TFType: types.TerraformString,
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
									TFType: types.TerraformList,
									Nested: &entities.Type{
										TFType: types.TerraformObject,
										Label:  "test_object",
									},
								},
								Attributes: []entities.Attribute{
									{
										Level:       1,
										Name:        "name",
										Description: "A string",
										Type: entities.Type{
											TFType: types.TerraformString,
										},
									},
									{
										Level:       1,
										Name:        "something_complex",
										Description: "Some other object",
										Type: entities.Type{
											TFType: types.TerraformObject,
											Label:  "nested_object",
										},
										Attributes: []entities.Attribute{
											{
												Level:       2,
												Name:        "nested_string",
												Description: "a nested string",
												Type: entities.Type{
													TFType: types.TerraformString,
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
	assert.NoError(t, err)

	got := strings.TrimSpace(buf.String())

	wantContent := test.ReadFixture(t, "markdown-structure.md")
	want := string(bytes.TrimSpace(wantContent))

	if diff := cmp.Diff(got, want); diff != "" {
		t.Errorf("Expected golden file to match result (-want +got):\n%s", diff)
	}
}

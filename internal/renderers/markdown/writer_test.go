package markdown

import (
	"bytes"
	"fmt"
	"io"
	"strings"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/madlambda/spells/assert"
	"github.com/mineiros-io/terradoc/internal/entities"
	"github.com/mineiros-io/terradoc/internal/types"
)

const (
	lineBreak = "\n"
)

func TestWriteSection(t *testing.T) {
	for _, tt := range []struct {
		desc    string
		section entities.Section
		want    mdSection
	}{
		{
			desc: "with description and 4 levels",
			section: entities.Section{
				Level:   4,
				Title:   "I AM THE TITLE!",
				Content: "Dude, this is a section",
			},
			want: mdSection{
				heading:     "#### I AM THE TITLE!",
				description: "Dude, this is a section",
			},
		},
		{
			desc: "without description and 1 level",
			section: entities.Section{
				Level: 1,
				Title: "section title",
			},
			want: mdSection{
				heading: "# section title",
			},
		},
	} {
		t.Run(tt.desc, func(t *testing.T) {
			buf := &bytes.Buffer{}

			writer := newTestWriter(t, buf)

			err := writer.writeSection(tt.section)
			assert.NoError(t, err)

			assertMarkdownHasSection(t, buf, tt.want)
		})
	}
}

func TestWriteVariable(t *testing.T) {
	for _, tt := range []struct {
		desc     string
		variable entities.Variable
		want     mdVariable
	}{
		{
			desc: "a required string variable with description and default that forces recreation",
			variable: entities.Variable{
				Name: "string_variable",
				Type: entities.Type{
					TFType: types.TerraformString,
				},
				ForcesRecreation: true,
				Required:         true,
				Description:      "i am a variable",
				Default:          []byte(`"default value"`),
			},
			want: mdVariable{
				item:        "- [**`string_variable`**](#var-string_variable): *(**Required** `string`, Forces new resource)*<a name=\"var-string_variable\"></a>",
				description: "i am a variable",
				defaults:    "Default is `\"default value\"`.",
			},
		},
		{
			desc: "an optional number variable with defaults that forces recreation",
			variable: entities.Variable{
				Name: "number_variable",
				Type: entities.Type{
					TFType: types.TerraformNumber,
				},
				ForcesRecreation: true,
				Required:         false,
				Default:          []byte("123"),
			},
			want: mdVariable{
				item:     "- [**`number_variable`**](#var-number_variable): *(Optional `number`, Forces new resource)*<a name=\"var-number_variable\"></a>",
				defaults: "Default is `123`.",
			},
		},
		{
			desc: "a bool variable",
			variable: entities.Variable{
				Name: "bool_variable",
				Type: entities.Type{
					TFType: types.TerraformBool,
				},
				ForcesRecreation: false,
				Required:         false,
			},
			want: mdVariable{
				item: "- [**`bool_variable`**](#var-bool_variable): *(Optional `bool`)*<a name=\"var-bool_variable\"></a>",
			},
		},
		{
			desc: "an object variable with readme example",
			variable: entities.Variable{
				Name: "obj_variable",
				Type: entities.Type{
					TFType: types.TerraformObject,
				},
				ForcesRecreation: true,
				Required:         true,
				ReadmeExample: `obj_variable = {
  a = "foo"
}
`,
			},
			want: mdVariable{
				item: "- [**`obj_variable`**](#var-obj_variable): *(**Required** `object`, Forces new resource)*<a name=\"var-obj_variable\"></a>",
				readmeExample: `obj_variable = {
    a = "foo"
  }
`,
			},
		},
	} {
		t.Run(tt.desc, func(t *testing.T) {
			buf := &bytes.Buffer{}

			writer := newTestWriter(t, buf)

			err := writer.writeVariable(tt.variable)
			assert.NoError(t, err)

			assertMarkdownHasVariable(t, buf, tt.want)
		})
	}
}

func TestWriteAttribute(t *testing.T) {
	for _, tt := range []struct {
		desc string
		attr entities.Attribute
		want mdAttribute
	}{
		{
			desc: "a required string attribute with description that forces recreation",
			attr: entities.Attribute{
				Level:       1,
				Name:        "string_attribute",
				Description: "i am this attribute's description",
				Type: entities.Type{
					TFType: types.TerraformString,
				},
				ForcesRecreation: true,
				Required:         true,
			},
			want: mdAttribute{
				item:        "  - [**`string_attribute`**](#attr-string_attribute-parent_var_name): *(**Required** `string`, Forces new resource)*<a name=\"attr-string_attribute-parent_var_name\"></a>",
				description: "  i am this attribute's description",
			},
		},
		{
			desc: "an optional number attribute that forces recreations",
			attr: entities.Attribute{
				Level: 2,
				Name:  "number_attribute",
				Type: entities.Type{
					TFType: types.TerraformNumber,
				},
				ForcesRecreation: true,
				Required:         false,
			},
			want: mdAttribute{
				item: "    - [**`number_attribute`**](#attr-number_attribute-parent_var_name): *(Optional `number`, Forces new resource)*<a name=\"attr-number_attribute-parent_var_name\"></a>",
			},
		},
		{
			desc: "a bool attribute",
			attr: entities.Attribute{
				Level: 0,
				Name:  "bool_attribute",
				Type: entities.Type{
					TFType: types.TerraformBool,
				},
				ForcesRecreation: false,
				Required:         false,
			},
			want: mdAttribute{
				item: "- [**`bool_attribute`**](#attr-bool_attribute-parent_var_name): *(Optional `bool`)*<a name=\"attr-bool_attribute-parent_var_name\"></a>",
			},
		},
		{
			desc: "an attribute with defautlts",
			attr: entities.Attribute{
				Level: 1,
				Name:  "i_have_a_default",
				Type: entities.Type{
					TFType: types.TerraformNumber,
				},
				Default: []byte("123"),
			},
			want: mdAttribute{
				item:        "  - [**`i_have_a_default`**](#attr-i_have_a_default-parent_var_name): *(Optional `number`)*<a name=\"attr-i_have_a_default-parent_var_name\"></a>",
				description: "  Default is `123`.",
			},
		},
	} {
		t.Run(tt.desc, func(t *testing.T) {
			buf := &bytes.Buffer{}

			writer := newTestWriter(t, buf)
			err := writer.writeAttribute(tt.attr, "parent_var_name")
			assert.NoError(t, err)

			assertMarkdownHasAttribute(t, buf, tt.want)
		})
	}
}

func TestWriteAttributeWithNested(t *testing.T) {
	t.Skip("write rendering tests for nested attributes once we know more about it")
}

// TODO: rewrite all? :D

type mdSection struct {
	heading     string
	description string
}

type mdVariable struct {
	item          string
	description   string
	defaults      string
	readmeExample string
}

type mdAttribute struct {
	item          string
	description   string
	defaults      string
	readmeExample string
}

func assertMarkdownHasSection(t *testing.T, buf *bytes.Buffer, md mdSection) {
	t.Helper()

	want := md.heading + lineBreak

	if md.description != "" {
		want += lineBreak + md.description + lineBreak
	}

	want += lineBreak

	if diff := cmp.Diff(want, buf.String()); diff != "" {
		t.Errorf("Expected section markdown to match (-want +got):\n%s", diff)
	}
}

func assertMarkdownHasVariable(t *testing.T, buf *bytes.Buffer, md mdVariable) {
	t.Helper()

	want := md.item + lineBreak

	if md.description != "" {
		want += fmt.Sprintf("\n  %s\n", md.description)
	}

	if md.defaults != "" {
		want += fmt.Sprintf("\n  %s\n", md.defaults)
	}

	if md.readmeExample != "" {
		// TODO: what's a better way of checking that indentation is right?
		want += fmt.Sprintf("\n  Example:\n\n  ```hcl\n  %s  \n  ```\n", md.readmeExample)
	}

	want += "\n"

	if diff := cmp.Diff(want, buf.String()); diff != "" {
		t.Errorf("Expected variable markdown to match (-want +got):\n%s", diff)
	}
}

func assertMarkdownHasAttribute(t *testing.T, buf *bytes.Buffer, md mdAttribute) {
	t.Helper()

	want := md.item + lineBreak

	if md.description != "" {
		want += fmt.Sprintf("\n  %s\n", md.description)
	}

	if md.defaults != "" {
		want += fmt.Sprintf("\n  %s\n", md.defaults)
	}

	if md.readmeExample != "" {
		want += fmt.Sprintf("\n  Example:\n\n  ```hcl\n  %s\n \n ```\n", md.readmeExample)
	}

	want += lineBreak

	got := strings.TrimLeft(buf.String(), "\n")

	if diff := cmp.Diff(want, got); diff != "" {
		t.Errorf("Expected variable markdown to match (-want +got):\n%s", diff)
	}
}

func newTestWriter(t *testing.T, buf io.Writer) *markdownWriter {
	writer, err := newMarkdownWriter(buf)
	assert.NoError(t, err)

	return writer
}

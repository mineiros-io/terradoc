package hclparser

import (
	"testing"

	"github.com/hashicorp/hcl/v2"
	"github.com/hashicorp/hcl/v2/hclsyntax"
	"github.com/madlambda/spells/assert"
	"github.com/mineiros-io/terradoc/internal/entities"
	"github.com/mineiros-io/terradoc/internal/types"
	"github.com/mineiros-io/terradoc/test"
)

var varTests = []struct {
	expression string
	want       entities.Type
}{
	{
		expression: `list(my_object)`,
		want: entities.Type{
			TFType: types.TerraformList,
			Nested: &entities.Type{
				TFType: types.TerraformObject,
				Label:  "my_object",
			},
		},
	},
	{
		expression: `list(string)`,
		want: entities.Type{
			TFType: types.TerraformList,
			Nested: &entities.Type{
				TFType: types.TerraformString,
			},
		},
	},
	{
		expression: `set(number)`,
		want: entities.Type{
			TFType: types.TerraformSet,
			Nested: &entities.Type{
				TFType: types.TerraformNumber,
			},
		},
	},
	{
		expression: `list(number)`,
		want: entities.Type{
			TFType: types.TerraformList,
			Nested: &entities.Type{
				TFType: types.TerraformNumber,
			},
		},
	},
	{
		expression: `list(another_object)`,
		want: entities.Type{
			TFType: types.TerraformList,
			Nested: &entities.Type{
				TFType: types.TerraformObject,
				Label:  "another_object",
			},
		},
	},
	{
		expression: `set(another_object)`,
		want: entities.Type{
			TFType: types.TerraformSet,
			Nested: &entities.Type{
				TFType: types.TerraformObject,
				Label:  "another_object",
			},
		},
	},
	{
		expression: `object(my_object_name)`,
		want: entities.Type{
			TFType: types.TerraformObject,
			Label:  "my_object_name",
		},
	},
	{
		expression: `map(my_object_name)`,
		want: entities.Type{
			TFType: types.TerraformMap,
			Nested: &entities.Type{
				TFType: types.TerraformObject,
				Label:  "my_object_name",
			},
		},
	},
	{
		expression: `object(another_object_name)`,
		want: entities.Type{
			TFType: types.TerraformObject,
			Label:  "another_object_name",
		},
	},
	{
		expression: `string`,
		want: entities.Type{
			TFType: types.TerraformString,
		},
	},
	{
		expression: `number`,
		want: entities.Type{
			TFType: types.TerraformNumber,
		},
	},
	{
		expression: `bool`,
		want: entities.Type{
			TFType: types.TerraformBool,
		},
	},
}

func TestGetVarTypeFromExpression(t *testing.T) {
	for _, tt := range varTests {
		t.Run(tt.expression, func(t *testing.T) {
			t.Run("when expression is literal", func(t *testing.T) {
				expr, parseDiags := hclsyntax.ParseExpression([]byte(tt.expression), "", hcl.Pos{Line: 1, Column: 1, Byte: 0})
				if parseDiags.HasErrors() {
					t.Errorf("Error parsing expression: %v", parseDiags.Errs())
				}

				got, err := getVarTypeFromExpression(expr)
				assert.NoError(t, err)

				test.AssertEqualTypes(t, tt.want, got)
			})

			t.Run("when expression is a string", func(t *testing.T) {
				got, err := getVarTypeFromString(tt.expression, hcl.Pos{Line: 1, Column: 1, Byte: 0})
				assert.NoError(t, err)

				test.AssertEqualTypes(t, tt.want, got)
			})
		})
	}
}

var outputTests = []struct {
	expression string
	want       entities.Type
}{
	{
		expression: `list(my_object)`,
		want: entities.Type{
			TFType: types.TerraformList,
			Nested: &entities.Type{
				TFType: types.TerraformObject,
				Label:  "my_object",
			},
		},
	},
	{
		expression: `list(string)`,
		want: entities.Type{
			TFType: types.TerraformList,
			Nested: &entities.Type{
				TFType: types.TerraformString,
			},
		},
	},
	{
		expression: `set(number)`,
		want: entities.Type{
			TFType: types.TerraformSet,
			Nested: &entities.Type{
				TFType: types.TerraformNumber,
			},
		},
	},
	{
		expression: `list(number)`,
		want: entities.Type{
			TFType: types.TerraformList,
			Nested: &entities.Type{
				TFType: types.TerraformNumber,
			},
		},
	},
	{
		expression: `list(another_object)`,
		want: entities.Type{
			TFType: types.TerraformList,
			Nested: &entities.Type{
				TFType: types.TerraformObject,
				Label:  "another_object",
			},
		},
	},
	{
		expression: `set(another_object)`,
		want: entities.Type{
			TFType: types.TerraformSet,
			Nested: &entities.Type{
				TFType: types.TerraformObject,
				Label:  "another_object",
			},
		},
	},
	{
		expression: `object(my_object_name)`,
		want: entities.Type{
			TFType: types.TerraformObject,
			Label:  "my_object_name",
		},
	},
	{
		expression: `map(my_object_name)`,
		want: entities.Type{
			TFType: types.TerraformMap,
			Nested: &entities.Type{
				TFType: types.TerraformObject,
				Label:  "my_object_name",
			},
		},
	},
	{
		expression: `object(another_object_name)`,
		want: entities.Type{
			TFType: types.TerraformObject,
			Label:  "another_object_name",
		},
	},
	{
		expression: `string`,
		want: entities.Type{
			TFType: types.TerraformString,
		},
	},
	{
		expression: `number`,
		want: entities.Type{
			TFType: types.TerraformNumber,
		},
	},
	{
		expression: `bool`,
		want: entities.Type{
			TFType: types.TerraformBool,
		},
	},
	{
		expression: `resource(foo_bar_baz)`,
		want: entities.Type{
			TFType: types.TerraformResource,
			Label:  "foo_bar_baz",
		},
	},
}

func TestGetOutputTypeFromExpression(t *testing.T) {
	for _, tt := range outputTests {
		t.Run(tt.expression, func(t *testing.T) {
			expr, parseDiags := hclsyntax.ParseExpression([]byte(tt.expression), "", hcl.Pos{Line: 1, Column: 1, Byte: 0})
			if parseDiags.HasErrors() {
				t.Errorf("Error parsing expression: %v", parseDiags.Errs())
			}

			got, err := getOutputTypeFromExpression(expr)
			assert.NoError(t, err)

			test.AssertEqualTypes(t, tt.want, got)
		})
	}
}

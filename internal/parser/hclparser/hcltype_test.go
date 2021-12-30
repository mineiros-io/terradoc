package hclparser

import (
	"testing"

	"github.com/hashicorp/hcl/v2"
	"github.com/hashicorp/hcl/v2/hclsyntax"
	"github.com/madlambda/spells/assert"
	"github.com/mineiros-io/terradoc/internal/types"
)

var tests = []struct {
	expression          string
	wantType            types.TerraformType
	wantTypeLabel       string
	wantNestedType      types.TerraformType
	wantNestedTypeLabel string
}{
	{
		expression:          `list(my_object)`,
		wantType:            types.TerraformList,
		wantNestedType:      types.TerraformObject,
		wantNestedTypeLabel: "my_object",
	},
	{
		expression:     `list(string)`,
		wantType:       types.TerraformList,
		wantNestedType: types.TerraformString,
	},
	{
		expression:     `set(number)`,
		wantType:       types.TerraformSet,
		wantNestedType: types.TerraformNumber,
	},
	{
		expression:     `list(number)`,
		wantType:       types.TerraformList,
		wantNestedType: types.TerraformNumber,
	},
	{
		expression:          `list(another_object)`,
		wantType:            types.TerraformList,
		wantNestedType:      types.TerraformObject,
		wantNestedTypeLabel: "another_object",
	},
	{
		expression:          `set(another_object)`,
		wantType:            types.TerraformSet,
		wantNestedType:      types.TerraformObject,
		wantNestedTypeLabel: "another_object",
	},
	{
		expression:    `object(my_object_name)`,
		wantType:      types.TerraformObject,
		wantTypeLabel: "my_object_name",
	},
	{
		expression:          `map(my_object_name)`,
		wantType:            types.TerraformMap,
		wantNestedType:      types.TerraformObject,
		wantNestedTypeLabel: "my_object_name",
	},
	{
		expression:    `object(another_object_name)`,
		wantType:      types.TerraformObject,
		wantTypeLabel: "another_object_name",
	},
	{
		expression: `string`,
		wantType:   types.TerraformString,
	},
	{
		expression: `number`,
		wantType:   types.TerraformNumber,
	},
	{
		expression: `bool`,
		wantType:   types.TerraformBool,
	},
}

func TestGetTypeFromExpression(t *testing.T) {
	for _, tt := range tests {
		t.Run(tt.expression, func(t *testing.T) {
			expr, parseDiags := hclsyntax.ParseExpression([]byte(tt.expression), "", hcl.Pos{Line: 1, Column: 1, Byte: 0})
			if parseDiags.HasErrors() {
				t.Errorf("Error parsing expression: %v", parseDiags.Errs())
			}

			got, err := getTypeFromExpression(expr)
			assert.NoError(t, err)

			assert.EqualStrings(t, tt.wantType.String(), got.TFType.String())
			assert.EqualStrings(t, tt.wantTypeLabel, got.TFTypeLabel)
			assert.EqualStrings(t, tt.wantNestedType.String(), got.NestedTFType.String())
			assert.EqualStrings(t, tt.wantNestedTypeLabel, got.NestedTFTypeLabel)
		})
	}
}

func TestGetTypeFromString(t *testing.T) {
	for _, tt := range tests {
		t.Run(tt.expression, func(t *testing.T) {
			got, err := getTypeFromString(tt.expression)
			assert.NoError(t, err)

			assert.EqualStrings(t, tt.wantType.String(), got.TFType.String())
			assert.EqualStrings(t, tt.wantTypeLabel, got.TFTypeLabel)
			assert.EqualStrings(t, tt.wantNestedType.String(), got.NestedTFType.String())
			assert.EqualStrings(t, tt.wantNestedTypeLabel, got.NestedTFTypeLabel)
		})
	}
}

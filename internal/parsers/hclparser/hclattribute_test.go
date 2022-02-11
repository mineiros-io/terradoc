package hclparser

import (
	"encoding/json"
	"fmt"
	"strings"
	"testing"

	"github.com/hashicorp/hcl/v2"
	"github.com/hashicorp/hcl/v2/ext/customdecode"
	"github.com/hashicorp/hcl/v2/hcltest"
	"github.com/madlambda/spells/assert"
	"github.com/mineiros-io/terradoc/internal/types"
	"github.com/zclconf/go-cty/cty"
)

func TestAttributeToString(t *testing.T) {
	attrName := "a-string"

	t.Run("when value is a cty.String", func(t *testing.T) {
		wantString := "test"

		exprValue := cty.StringVal(wantString)

		attr := newMockAttribute(attrName, exprValue)

		res, err := attr.String()
		assert.NoError(t, err)
		assert.EqualStrings(t, wantString, res)
	})

	t.Run("when value is not convertable to string", func(t *testing.T) {
		// test that it doesn't trigger cty's panic calls
		wantErrorMSGContains := fmt.Sprintf("could not convert %q to string", attrName)
		exprValue := customdecode.ExpressionVal(&fakeHCLExpression{})

		attr := newMockAttribute(attrName, exprValue)

		res, err := attr.String()
		assert.Error(t, err)

		if !strings.Contains(err.Error(), wantErrorMSGContains) {
			t.Errorf("Expected error to contain %q but got %q instead", wantErrorMSGContains, err.Error())
		}

		if res != "" {
			t.Errorf("Expected result to be empty. Got %q instead", res)
		}
	})
}

func TestAttributeToBool(t *testing.T) {
	attrName := "a-bool"
	t.Run("when value is a cty.Bool", func(t *testing.T) {
		wantBool := true

		exprValue := cty.BoolVal(wantBool)

		attr := newMockAttribute(attrName, exprValue)

		res, err := attr.Bool()
		assert.NoError(t, err)

		if res != wantBool {
			t.Errorf("Expected returned value to be %t but got %t instead", wantBool, res)
		}
	})

	t.Run("when value is not convertable to bool", func(t *testing.T) {
		// test that it doesn't trigger cty's panic calls
		wantErrorMSGContains := fmt.Sprintf("could not convert %q to bool", attrName)
		exprValue := customdecode.ExpressionVal(&fakeHCLExpression{})

		attr := newMockAttribute(attrName, exprValue)

		res, err := attr.Bool()
		assert.Error(t, err)

		if !strings.Contains(err.Error(), wantErrorMSGContains) {
			t.Errorf("Expected error to contain %q but got %q instead", wantErrorMSGContains, err.Error())
		}

		if res != false {
			t.Errorf("Expected result to be false. Got %t instead", res)
		}
	})
}

func TestAttributeToJSONValue(t *testing.T) {
	for _, tt := range []struct {
		desc  string
		value string
	}{
		{
			desc:  "when value is a list",
			value: `[1, 2, "c", [3, "a", "b"]]`,
		},
		{
			desc:  "when value is a number",
			value: "123",
		},
		{
			desc:  "when value is a string",
			value: `"foo"`,
		},
		{
			desc:  "when value is a map",
			value: `{a=123, b="foo"}`,
		},
	} {
		t.Run(tt.desc, func(t *testing.T) {
			// test that the returned value is not an escaped json string
			expr := hcltest.MockExprLiteral(cty.StringVal(tt.value))
			attr := &HCLAttribute{&hcl.Attribute{Expr: expr}}

			res, err := attr.RawJSON()
			assert.NoError(t, err)

			var strRes string
			err = json.Unmarshal(res, &strRes)

			assert.NoError(t, err)
			assert.EqualStrings(t, tt.value, strRes)
		})
	}
}

func TestAttributeToTypeValidPrimaryType(t *testing.T) {
	for _, tt := range []struct {
		desc                  string
		exprValue             string
		expectedTerraformType types.TerraformType
	}{
		{
			desc:                  "when type is bool",
			exprValue:             "bool",
			expectedTerraformType: types.TerraformBool,
		},
		{
			desc:                  "when type is string",
			exprValue:             "string",
			expectedTerraformType: types.TerraformString,
		},
		{
			desc:                  "when type is number",
			exprValue:             "number",
			expectedTerraformType: types.TerraformNumber,
		},
	} {
		t.Run(tt.desc, func(t *testing.T) {
			attr := newTypeAttribute(tt.exprValue, tt.exprValue)

			res, err := attr.VarType()
			assert.NoError(t, err)

			assert.EqualInts(t, int(tt.expectedTerraformType), int(res.TFType))
			assert.EqualStrings(t, "", res.Label)

			// ensure that type does not have a nested type
			if res.HasNestedType() {
				t.Errorf("type %q should not have a nested type", tt.expectedTerraformType)
			}
		})
	}
}

func TestAttributeToTypeInvalidTypes(t *testing.T) {
	for _, tt := range []struct {
		desc             string
		exprValue        string
		expectedErrorMSG string
	}{
		{
			desc:             "when an invalid primary type is given",
			exprValue:        "foo",
			expectedErrorMSG: "type \"foo\" is invalid",
		},
		{
			desc:             "when type is a list without arguments",
			exprValue:        "list",
			expectedErrorMSG: "type \"list\" needs an argument",
		},
		{
			desc:             "when type is an object without definition",
			exprValue:        "object",
			expectedErrorMSG: "type \"object\" needs an argument",
		},
		{
			desc:             "when type is a tuple without definition",
			exprValue:        "tuple",
			expectedErrorMSG: "type \"tuple\" needs an argument",
		},
		{
			desc:             "when type is a map without definition",
			exprValue:        "map",
			expectedErrorMSG: "type \"map\" needs an argument",
		},
	} {
		t.Run(tt.desc, func(t *testing.T) {
			attr := newTypeAttribute(tt.exprValue, tt.exprValue)

			res, err := attr.VarType()
			assert.Error(t, err)

			if !strings.Contains(err.Error(), tt.expectedErrorMSG) {
				t.Errorf("Expected error to contain %q. Got %q instead", tt.expectedErrorMSG, err.Error())
			}

			assert.EqualInts(t, int(types.TerraformEmptyType), int(res.TFType))
			assert.EqualStrings(t, "", res.Label)

			if res.HasNestedType() {
				t.Error("empty type cannot have a nested type")
			}
		})
	}
}

func TestAttributeToTerraformTypeValidComplexType(t *testing.T) {
	t.Skip("I'm not sure how tf I'll test this")
}

type fakeHCLExpression struct {
	value cty.Value
}

func (expr fakeHCLExpression) Variables() []hcl.Traversal {
	return nil
}

func (expr fakeHCLExpression) Range() hcl.Range {
	return hcl.Range{}
}

func (expr fakeHCLExpression) StartRange() hcl.Range {
	return hcl.Range{}
}

func (expr fakeHCLExpression) Value(_ *hcl.EvalContext) (cty.Value, hcl.Diagnostics) {
	return expr.value, nil
}

func newMockAttribute(name string, returnValue cty.Value) *HCLAttribute {
	fakeExpr := &fakeHCLExpression{
		value: returnValue,
	}
	attr := &hcl.Attribute{Name: name, Expr: fakeExpr}

	return &HCLAttribute{attr}
}

func newTypeAttribute(name, typeStr string) *HCLAttribute {
	expr := hcltest.MockExprVariable(typeStr)
	attr := &hcl.Attribute{Name: name, Expr: expr}

	return &HCLAttribute{attr}
}

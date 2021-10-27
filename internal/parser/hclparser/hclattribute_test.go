package hclparser

import (
	"encoding/json"
	"fmt"
	"strings"
	"testing"

	"github.com/hashicorp/hcl/v2"
	"github.com/hashicorp/hcl/v2/ext/customdecode"
	"github.com/hashicorp/hcl/v2/hcltest"
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
		if err != nil {
			t.Fatalf("Expected no error but got %q instead", err.Error())
		}

		if res != wantString {
			t.Errorf("Expected returned value to be %q but got %q instead", wantString, res)
		}
	})

	t.Run("when value is not convertable to string", func(t *testing.T) {
		// test that it doesn't trigger cty's panic calls
		wantErrorMSGContains := fmt.Sprintf("could not convert %q to string", attrName)
		exprValue := customdecode.ExpressionVal(&fakeHCLExpression{})

		attr := newMockAttribute(attrName, exprValue)

		res, err := attr.String()
		if err == nil {
			t.Fatalf("Expected no error but got none")
		}

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
		if err != nil {
			t.Fatalf("Expected no error but got %q instead", err.Error())
		}

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
		if err == nil {
			t.Fatalf("Expected no error but got none")
		}

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
			attr := &hclAttribute{&hcl.Attribute{Expr: expr}}

			res, err := attr.RawJSON()
			if err != nil {
				t.Fatalf("Expected no error. Got %q instead", err)
			}

			var strRes string
			err = json.Unmarshal(res, &strRes)

			if err != nil {
				t.Fatalf("Expected no error. Got %q instead", err)
			}

			if strRes != tt.value {
				t.Errorf("Expected result to be %q. Got %q instead", tt.value, strRes)
			}
		})
	}
}

func TestAttributeToTerraformTypeValidPrimaryType(t *testing.T) {
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

			res, err := attr.TerraformType()
			if err != nil {
				t.Fatalf("Expected no error. Got %q instead", err)
			}

			if res.Type != tt.expectedTerraformType {
				t.Errorf("Expected type to be %q. Got %q instead", tt.expectedTerraformType.String(), res)
			}
		})
	}
}

func TestAttributeToTerraformTypeInvalidTypes(t *testing.T) {
	for _, tt := range []struct {
		desc             string
		exprValue        string
		expectedErrorMSG string
	}{
		{
			desc:             "when an invalid primary type is given",
			exprValue:        "foo",
			expectedErrorMSG: "The keyword \"foo\" is not a valid type specification",
		},
		{
			desc:             "when type is a list without arguments",
			exprValue:        "list",
			expectedErrorMSG: "The list type constructor requires one argument specifying the element type",
		},
		{
			desc:             "when type is an object without definition",
			exprValue:        "object",
			expectedErrorMSG: "The object type constructor requires one argument specifying the attribute types and values as a map",
		},
		{
			desc:             "when type is a tuple without definition",
			exprValue:        "tuple",
			expectedErrorMSG: "The tuple type constructor requires one argument specifying the element types as a list",
		},
		{
			desc:             "when type is a map without definition",
			exprValue:        "map",
			expectedErrorMSG: "The map type constructor requires one argument specifying the element type",
		},
	} {
		t.Run(tt.desc, func(t *testing.T) {
			attr := newTypeAttribute(tt.exprValue, tt.exprValue)

			res, err := attr.TerraformType()
			if err == nil {
				t.Fatal("Expected  but got none")
			}

			if !strings.Contains(err.Error(), tt.expectedErrorMSG) {
				t.Errorf("Expected error to contain %q. Got %q instead", tt.expectedErrorMSG, err.Error())
			}

			if res.Type != types.TerraformEmptyType {
				t.Errorf("Expected returned type to be %q. Got %q instead", types.TerraformEmptyType, res)
			}
		})
	}
}

func TestAttributeToTerraformTypeValidComplexType(t *testing.T) {
	t.Skip("I'm not sure how tf I'll test this")
}

func TestAttributeToHCL(t *testing.T) {
	t.Run("with valid HCL", func(t *testing.T) {
		// a block like `readme_example = {values = [{key = "value"}]}`
		// gets parsed into a cty object like the following
		objVal := cty.ObjectVal(map[string]cty.Value{
			"values": cty.ListVal([]cty.Value{
				cty.ObjectVal(
					map[string]cty.Value{
						"key": cty.StringVal("value"),
					},
				),
			}),
		})
		// we need to ensure indentation is maintained
		wantVal := `values = [{
    key = "value"
  }]`

		expr := hcltest.MockExprLiteral(objVal)
		attr := &hclAttribute{&hcl.Attribute{Name: "hcl", Expr: expr}}

		res, err := attr.HCLString()
		if err != nil {
			t.Fatalf("Expected no error. Got %q instead", err)
		}

		if res != wantVal {
			t.Errorf("Expected result to be %q. Got %q instead", wantVal, res)
		}
	})
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

func newMockAttribute(name string, returnValue cty.Value) *hclAttribute {
	fakeExpr := &fakeHCLExpression{
		value: returnValue,
	}
	attr := &hcl.Attribute{Name: name, Expr: fakeExpr}

	return &hclAttribute{attr}
}

func newTypeAttribute(name, typeStr string) *hclAttribute {
	expr := hcltest.MockExprVariable(typeStr)
	attr := &hcl.Attribute{Name: name, Expr: expr}

	return &hclAttribute{attr}
}

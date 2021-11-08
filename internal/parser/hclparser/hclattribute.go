package hclparser

import (
	"bytes"
	"encoding/json"
	"fmt"

	"github.com/hashicorp/hcl/v2"
	"github.com/hashicorp/hcl/v2/ext/typeexpr"
	"github.com/hashicorp/hcl/v2/hclwrite"
	"github.com/mineiros-io/terradoc/internal/types"
	"github.com/zclconf/go-cty/cty"
	"github.com/zclconf/go-cty/cty/convert"
	ctyjson "github.com/zclconf/go-cty/cty/json"
)

type hclAttribute struct {
	*hcl.Attribute
}

func (a *hclAttribute) isNil() bool {
	return a.Attribute == nil
}

func (a *hclAttribute) String() (string, error) {
	if a.isNil() {
		return "", nil
	}

	val, diags := a.Expr.Value(nil)
	if diags.HasErrors() {
		return "", fmt.Errorf("Error getting string value for %q: %v", a.Name, diags.Errs())
	}

	// use cty's convert pkg to prevent panic if value is not a string
	strVal, err := convert.Convert(val, cty.String)
	if err != nil {
		return "", fmt.Errorf("Could not convert value of %q to string: %v", a.Name, err)
	}

	return strVal.AsString(), nil
}

func (a *hclAttribute) Bool() (bool, error) {
	if a.isNil() {
		return false, nil
	}

	val, diags := a.Expr.Value(nil)
	if diags.HasErrors() {
		return false, fmt.Errorf("Error getting bool value for %q: %v", a.Name, diags.Errs())
	}

	// use cty's convert pkg to prevent panic if value is not a bool
	boolVal, err := convert.Convert(val, cty.Bool)
	if err != nil {
		return false, fmt.Errorf("Could not convert value of %q to bool: %s", a.Name, err)
	}

	return boolVal.True(), nil
}

func (a *hclAttribute) RawJSON() (json.RawMessage, error) {
	if a.isNil() {
		return nil, nil
	}

	val, diags := a.Expr.Value(nil)
	if diags.HasErrors() {
		return nil, fmt.Errorf("Error getting JSON value for %q: %v", a.Name, diags.Errs())
	}

	// convert cty.Value to SimpleJSONValue to get the correct decoding of its internal value
	jsonVal := ctyjson.SimpleJSONValue{Value: val}

	src, err := jsonVal.MarshalJSON()
	if err != nil {
		return nil, err
	}

	return json.RawMessage(src), nil
}

func (a *hclAttribute) TerraformType() (types.TerraformType, error) {
	if a.isNil() {
		return types.TerraformInvalidType, nil
	}

	val, diags := typeexpr.TypeConstraint(a.Expr)
	if diags.HasErrors() {
		return types.TerraformInvalidType,
			fmt.Errorf("Error converting %q to TerraformType: %v", a.Name, diags.Errs())
	}

	typeName := val.FriendlyName()

	tfType, exists := types.TerraformTypes[typeName]
	if !exists {
		return types.TerraformInvalidType, fmt.Errorf("Could get terraform type for %q", a.Name)
	}

	return tfType, nil
}

func (a *hclAttribute) HCLString() (string, error) {
	if a.isNil() {
		return "", nil
	}

	val, diags := a.Expr.Value(nil)
	if diags.HasErrors() {
		return "", fmt.Errorf("Error getting HCL value for %q: %v", a.Name, diags.Errs())
	}

	tk := hclwrite.TokensForValue(val)
	// As the fetched attribute has the format `{ content }\n`,
	// remove the extra characters to have only the HCL content
	bVal := bytes.Trim(bytes.TrimSpace(tk.Bytes()), "{}\n ")

	return string(bVal), nil
}

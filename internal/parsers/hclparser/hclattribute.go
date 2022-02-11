package hclparser

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/hashicorp/hcl/v2"
	"github.com/hashicorp/hcl/v2/hclwrite"
	"github.com/mineiros-io/terradoc/internal/entities"
	"github.com/zclconf/go-cty/cty"
	"github.com/zclconf/go-cty/cty/convert"
	ctyjson "github.com/zclconf/go-cty/cty/json"
)

type HCLAttribute struct {
	*hcl.Attribute
}

func (a *HCLAttribute) String() (string, error) {
	if a == nil {
		return "", nil
	}

	val, diags := a.Expr.Value(nil)
	if diags.HasErrors() {
		return "", fmt.Errorf("getting string value for %q: %v", a.Name, diags.Errs())
	}

	// use cty's convert pkg to prevent panic if value is not a string
	strVal, err := convert.Convert(val, cty.String)
	if err != nil {
		return "", fmt.Errorf("could not convert %q to string: %v", a.Name, err)
	}

	return strings.TrimSpace(strVal.AsString()), nil
}

func (a *HCLAttribute) Bool() (bool, error) {
	if a == nil {
		return false, nil
	}

	val, diags := a.Expr.Value(nil)
	if diags.HasErrors() {
		return false, fmt.Errorf("fetching bool value for %q: %v", a.Name, diags.Errs())
	}

	// use cty's convert pkg to prevent panic if value is not a bool
	boolVal, err := convert.Convert(val, cty.Bool)
	if err != nil {
		return false, fmt.Errorf("could not convert %q to bool: %s", a.Name, err)
	}

	return boolVal.True(), nil
}

func (a *HCLAttribute) Keyword() string {
	return hcl.ExprAsKeyword(a.Expr)
}

func (a *HCLAttribute) RawJSON() (json.RawMessage, error) {
	if a == nil {
		return nil, nil
	}

	if len(a.Expr.Variables()) > 0 {
		return getRawVariables(a.Expr), nil
	}

	val, diags := a.Expr.Value(nil)
	if diags.HasErrors() {
		return nil, fmt.Errorf("could not fetch JSON value for %q: %v", a.Name, diags.Errs())
	}

	// convert cty.Value to SimpleJSONValue to get the correct decoding of its internal value
	jsonVal := ctyjson.SimpleJSONValue{Value: val}

	src, err := jsonVal.MarshalJSON()
	if err != nil {
		return nil, err
	}

	return json.RawMessage(src), nil
}

func (a *HCLAttribute) VarType() (entities.Type, error) {
	if a == nil {
		return entities.Type{}, nil
	}

	return GetVarTypeFromExpression(a.Expr)
}

func (a *HCLAttribute) VarTypeFromString() (entities.Type, error) {
	if a == nil {
		return entities.Type{}, nil
	}

	val, diags := a.Expr.Value(nil)
	if diags.HasErrors() {
		return entities.Type{}, fmt.Errorf("could not fetch type string value for %q: %v", a.Name, diags.Errs())
	}

	return getVarTypeFromString(val.AsString(), a.Range.Start)
}

func getRawVariables(expr hcl.Expression) json.RawMessage {
	var varValue []byte

	for _, exprVar := range expr.Variables() {
		tk := hclwrite.TokensForTraversal(exprVar)

		varValue = append(varValue, tk.Bytes()...)
	}

	return varValue
}

func (a *HCLAttribute) OutputType() (entities.Type, error) {
	if a == nil {
		return entities.Type{}, nil
	}

	return GetOutputTypeFromExpression(a.Expr)
}

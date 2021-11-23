package hclparser

import (
	"encoding/json"
	"fmt"

	"github.com/hashicorp/hcl/v2"
	"github.com/hashicorp/hcl/v2/ext/typeexpr"
	"github.com/mineiros-io/terradoc/internal/entities"
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
		return "", fmt.Errorf("getting string value for %q: %v", a.Name, diags.Errs())
	}

	// use cty's convert pkg to prevent panic if value is not a string
	strVal, err := convert.Convert(val, cty.String)
	if err != nil {
		return "", fmt.Errorf("could not convert %q to string: %v", a.Name, err)
	}

	return strVal.AsString(), nil
}

func (a *hclAttribute) Bool() (bool, error) {
	if a.isNil() {
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

func (a *hclAttribute) RawJSON() (json.RawMessage, error) {
	if a.isNil() {
		return nil, nil
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

func (a *hclAttribute) TerraformType() (entities.TerraformType, error) {
	if a.isNil() {
		return entities.TerraformType{}, nil
	}

	val, diags := typeexpr.TypeConstraint(a.Expr)
	if diags.HasErrors() {
		return entities.TerraformType{},
			fmt.Errorf("could not convert %q to TerraformType: %v", a.Name, diags.Errs())
	}

	switch {
	case val.IsPrimitiveType():
		return getTerraformPrimitiveType(val)
	case val.IsCollectionType():
		return getTerraformCollectionType(val)
	case val.IsObjectType():
		return getTerraformObjectType(val)
	case val.IsTupleType():
		return getTerraformTupleType(val)
	case val.HasDynamicTypes():
		return entities.TerraformType{Type: types.TerraformAny}, nil
	}

	return entities.TerraformType{}, fmt.Errorf("could not generate TerraformType for attribute %q", a.Name)

}

func getTerraformPrimitiveType(ctyType cty.Type) (entities.TerraformType, error) {
	typeName := typeexpr.TypeString(ctyType)

	tfType, exists := types.TerraformTypes[typeName]
	if !exists {
		return entities.TerraformType{}, fmt.Errorf("could not find TerraformType for cty type %q", typeName)
	}

	return entities.TerraformType{Type: tfType}, nil
}

func getTerraformObjectType(cty.Type) (entities.TerraformType, error) {
	// TODO
	return entities.TerraformType{Type: types.TerraformObject}, nil
}

func getTerraformTupleType(cty.Type) (entities.TerraformType, error) {
	// TODO
	return entities.TerraformType{Type: types.TerraformTuple}, nil
}

func getTerraformCollectionType(ctyType cty.Type) (entities.TerraformType, error) {
	nested := typeexpr.TypeString(ctyType.ElementType())
	nestedType, exists := types.TerraformTypes[nested]
	if !exists {
		return entities.TerraformType{}, fmt.Errorf("could not find TerraformType for %q", nested)
	}

	switch {
	case ctyType.IsListType():
		return entities.TerraformType{NestedType: nestedType, Type: types.TerraformList}, nil
	case ctyType.IsSetType():
		return entities.TerraformType{NestedType: nestedType, Type: types.TerraformSet}, nil
	case ctyType.IsMapType():
		return entities.TerraformType{NestedType: nestedType, Type: types.TerraformMap}, nil
	}

	return entities.TerraformType{},
		fmt.Errorf("could not get type information for %q", typeexpr.TypeString(ctyType))
}

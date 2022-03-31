package hclparser

import (
	"fmt"
	"strings"

	"github.com/hashicorp/hcl/v2"
	"github.com/hashicorp/hcl/v2/hclsyntax"
	"github.com/mineiros-io/terradoc/internal/entities"
	"github.com/mineiros-io/terradoc/internal/types"
)

// TODO: refactor - the next two functions are identical apart from the types slice that they use
func GetVarTypeFromExpression(expr hcl.Expression) (entities.Type, error) {
	t, err := processType(expr, types.TerraformEmptyType)
	if err != nil {
		return entities.Type{}, err
	}

	found := false
	for _, tt := range types.VariableTypes {
		if t.TFType == tt {
			found = true
		}
	}

	if !found {
		return entities.Type{}, fmt.Errorf("%q is not a valid variable type", t.TFType)
	}

	return t, nil
}

func GetOutputTypeFromExpression(expr hcl.Expression) (entities.Type, error) {
	t, err := processType(expr, types.TerraformEmptyType)
	if err != nil {
		return entities.Type{}, err
	}

	found := false
	for _, tt := range types.OutputTypes {
		if t.TFType == tt {
			found = true
		}
	}

	if !found {
		return entities.Type{}, fmt.Errorf("%q is not a valid output type", t.TFType)
	}

	return t, nil
}

// TODO: remove once we're sure we don't need `readme_example` attributes anymore
func getVarTypeFromString(str string, startRange hcl.Pos) (entities.Type, error) {
	expr, parseDiags := hclsyntax.ParseExpression([]byte(str), "", startRange)
	if parseDiags.HasErrors() {
		return entities.Type{}, fmt.Errorf("parsing type string expression: %v", parseDiags.Errs())
	}

	return GetVarTypeFromExpression(expr)
}

func processType(expr hcl.Expression, parentType types.TerraformType) (entities.Type, error) {
	kw := hcl.ExprAsKeyword(expr)

	tfType, found := types.TerraformTypes(kw)
	if found {
		if tfType.IsComplex() {
			return entities.Type{}, fmt.Errorf("type %q needs an argument", tfType.String())
		}

		return entities.Type{TFType: tfType}, nil
	}

	if kw != "" && !isComplexTypeExpression(kw) {
			if parentType.IsComplex() {
				// this is needed so we interpret stuff like list(my_object_label) as list(object(my_object_label))
				return entities.Type{TFType: types.TerraformObject, Label: kw}, nil
			}
			return entities.Type{}, fmt.Errorf("type %q is invalid", kw)
		}
	}

	call, diags := hcl.ExprCall(expr)
	if diags.HasErrors() {
		return entities.Type{}, fmt.Errorf("evaluating expression: %v", diags.Errs())
	}

	switch call.Name {
	case "bool", "string", "number", "any":
		return entities.Type{}, fmt.Errorf("type %q does not accept any argument", call.Name)
	}

	if len(call.Arguments) != 1 {
		return entities.Type{}, fmt.Errorf("type %q accepts only 1 argument", call.Name)
	}

	switch call.Name {
	case "list":
		nestedType, err := processType(call.Arguments[0], types.TerraformList)

		return entities.Type{
			TFType: types.TerraformList,
			Nested: &nestedType,
		}, err

	case "set":
		nestedType, err := processType(call.Arguments[0], types.TerraformSet)

		return entities.Type{
			TFType: types.TerraformSet,
			Nested: &nestedType,
		}, err

	case "map":
		nestedType, err := processType(call.Arguments[0], types.TerraformMap)

		return entities.Type{
			TFType: types.TerraformSet,
			Nested: &nestedType,
		}, err
	case "object":
		objectLabel := hcl.ExprAsKeyword(call.Arguments[0])

		return entities.Type{
			TFType: types.TerraformObject,
			Label:  objectLabel,
		}, nil
	case "resource":
		resourceLabel := hcl.ExprAsKeyword(call.Arguments[0])

		return entities.Type{
			TFType: types.TerraformResource,
			Label:  resourceLabel,
		}, nil
	case "module":
		moduleLabel := hcl.ExprAsKeyword(call.Arguments[0])

		return entities.Type{
			TFType: types.TerraformModule,
			Label:  moduleLabel,
		}, nil
	}

	return entities.Type{}, nil
}

func isComplexTypeExpression(expression string) bool {
	for _, tfType := range types.SupportedTerraformTypes {
		if strings.Contains(expression, tfType.String()) {
			return true
		}
	}

	return false
}

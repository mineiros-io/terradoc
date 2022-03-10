package hclparser

import (
	"fmt"
	"strings"

	"github.com/hashicorp/hcl/v2"
	"github.com/hashicorp/hcl/v2/hclsyntax"
	"github.com/mineiros-io/terradoc/internal/entities"
	"github.com/mineiros-io/terradoc/internal/types"
)

func GetVarTypeFromExpression(expr hcl.Expression) (entities.Type, error) {
	return processType(expr, types.TerraformEmptyType)
}

func GetOutputTypeFromExpression(expr hcl.Expression) (entities.Type, error) {
	return processType(expr, types.TerraformEmptyType)
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

	if kw != "" {
		if !isComplexTypeExpression(kw) {
			// TODO: no nested ifs, man!
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
		return entities.Type{}, fmt.Errorf("%q does not accept any argument", call.Name)
	}

	if len(call.Arguments) != 1 {
		return entities.Type{}, fmt.Errorf("%q accepts only 1 argument", call.Name)
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

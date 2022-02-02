package hclparser

import (
	"fmt"
	"strings"

	"github.com/hashicorp/hcl/v2"
	"github.com/hashicorp/hcl/v2/hclsyntax"
	"github.com/mineiros-io/terradoc/internal/entities"
	"github.com/mineiros-io/terradoc/internal/types"
	"github.com/zclconf/go-cty/cty"
	"github.com/zclconf/go-cty/cty/function"
	"github.com/zclconf/go-cty/cty/gocty"
)

func GetVarTypeFromExpression(expr hcl.Expression) (entities.Type, error) {
	return getTypeFromExpression(expr, varFunctions())
}

func GetOutputTypeFromExpression(expr hcl.Expression) (entities.Type, error) {
	return getTypeFromExpression(expr, outputFunctions())
}

func getTypeFromExpression(expr hcl.Expression, ctxFunctions map[string]function.Function) (entities.Type, error) {
	kw := hcl.ExprAsKeyword(expr)

	switch kw {
	case "string", "number", "bool":
		tfType, ok := types.TerraformTypes(kw)
		if !ok {
			return entities.Type{}, fmt.Errorf("could not get terraform type for %q", kw)
		}

		return entities.Type{TFType: tfType}, nil
	case "list", "object", "map", "tuple":
		// invalid as these types should be function calls
		return entities.Type{}, fmt.Errorf("type %q needs an argument", kw)
	}

	// TODO: how to make this decent?
	if kw != "" && !(strings.HasPrefix(kw, "list") ||
		strings.HasPrefix(kw, "object") ||
		strings.HasPrefix(kw, "map") ||
		strings.HasPrefix(kw, "tuple")) {
		return entities.Type{}, fmt.Errorf("type %q is invalid", kw)
	}

	return getComplexType(expr, ctxFunctions)
}

// this function exists to make it possible to parse `type` attribute expressions and `readme_type`
// attribute strings in the same way, so they are compatible even though they have different types
func getVarTypeFromString(str string, startRange hcl.Pos) (entities.Type, error) {
	expr, parseDiags := hclsyntax.ParseExpression([]byte(str), "", startRange)
	if parseDiags.HasErrors() {
		return entities.Type{}, fmt.Errorf("parsing type string expression: %v", parseDiags.Errs())
	}

	return GetVarTypeFromExpression(expr)
}

func getComplexType(expr hcl.Expression, ctxFunctions map[string]function.Function) (entities.Type, error) {
	got, exprDiags := expr.Value(getEvalContextForExpr(expr, ctxFunctions))
	if exprDiags.HasErrors() {
		return entities.Type{}, fmt.Errorf("getting expression value: %v", exprDiags.Errs())
	}
	var tfType, nestedTFType types.TerraformType

	err := gocty.FromCtyValue(got.GetAttr("type"), &tfType)
	if err != nil {
		return entities.Type{}, fmt.Errorf("getting type definition: %v", err)
	}

	err = gocty.FromCtyValue(got.GetAttr("nestedType"), &nestedTFType)
	if err != nil {
		return entities.Type{}, fmt.Errorf("getting nested type definition: %v", err)
	}

	var typeLabel, nestedTypeLabel string
	err = gocty.FromCtyValue(got.GetAttr("typeLabel"), &typeLabel)
	if err != nil {
		return entities.Type{}, fmt.Errorf("getting type label: %v", err)
	}

	err = gocty.FromCtyValue(got.GetAttr("nestedTypeLabel"), &nestedTypeLabel)
	if err != nil {
		return entities.Type{}, fmt.Errorf("getting nested type label: %v", err)
	}

	typeDef := entities.Type{
		TFType: tfType,
		Label:  typeLabel,
	}

	if nestedTFType != types.TerraformEmptyType {
		typeDef.Nested = &entities.Type{
			TFType: nestedTFType,
			Label:  nestedTypeLabel,
		}
	}

	return typeDef, nil
}

func nestedTypeFunc(tfType types.TerraformType) function.Function {
	return function.New(&function.Spec{
		Params: []function.Parameter{
			{
				Name:             "nestedTypeLabel",
				Type:             cty.String,
				AllowDynamicType: true,
			},
		},
		Type: function.StaticReturnType(cty.Object(typeObj())),
		Impl: func(args []cty.Value, retType cty.Type) (cty.Value, error) {
			var nestedType types.TerraformType
			var nestedLabel string

			nestedTypeName := args[0].AsString()

			nestedType, ok := types.TerraformTypes(nestedTypeName)
			if !ok {
				nestedType = types.TerraformObject
				nestedLabel = nestedTypeName
			}

			return cty.ObjectVal(map[string]cty.Value{
				"type":            cty.NumberIntVal(int64(tfType)),
				"typeLabel":       cty.StringVal(""), // need to pass empty value here so cty doesn't panic
				"nestedType":      cty.NumberIntVal(int64(nestedType)),
				"nestedTypeLabel": cty.StringVal(nestedLabel),
			}), nil
		},
	})
}

func complexTypeFunc(tfType types.TerraformType) function.Function {
	return function.New(&function.Spec{
		Params: []function.Parameter{
			{
				Name:             "typeLabel",
				Type:             cty.String,
				AllowDynamicType: true,
			},
		},
		Type: function.StaticReturnType(cty.Object(typeObj())),
		Impl: func(args []cty.Value, retType cty.Type) (cty.Value, error) {
			typeLabel := args[0].AsString()

			return cty.ObjectVal(map[string]cty.Value{
				"type":      cty.NumberIntVal(int64(tfType)),
				"typeLabel": cty.StringVal(typeLabel),
				// the following empty values need to be set to the attributes
				// otherwise cty panics
				"nestedType":      cty.NumberIntVal(int64(types.TerraformEmptyType)),
				"nestedTypeLabel": cty.StringVal(""), //
			}), nil
		},
	})
}

func getVariablesMap(expr hcl.Expression) map[string]cty.Value {
	varMap := make(map[string]cty.Value)
	for _, variable := range expr.Variables() {
		name := variable.RootName()

		varMap[name] = cty.StringVal(name)
	}

	return varMap
}

func getEvalContextForExpr(expr hcl.Expression, ctxFunctions map[string]function.Function) *hcl.EvalContext {
	return &hcl.EvalContext{
		Functions: ctxFunctions,
		Variables: getVariablesMap(expr),
	}
}

func varFunctions() map[string]function.Function {
	return map[string]function.Function{
		"object": complexTypeFunc(types.TerraformObject),
		"map":    nestedTypeFunc(types.TerraformMap),
		"list":   nestedTypeFunc(types.TerraformList),
		"set":    nestedTypeFunc(types.TerraformSet),
	}
}

func outputFunctions() map[string]function.Function {
	return map[string]function.Function{
		"resource": complexTypeFunc(types.TerraformResource),
		"object":   complexTypeFunc(types.TerraformObject),
		"map":      nestedTypeFunc(types.TerraformMap),
		"list":     nestedTypeFunc(types.TerraformList),
		"set":      nestedTypeFunc(types.TerraformSet),
	}
}

func typeObj() map[string]cty.Type {
	return map[string]cty.Type{
		"type":            cty.Number,
		"typeLabel":       cty.String,
		"nestedType":      cty.Number,
		"nestedTypeLabel": cty.String,
	}
}

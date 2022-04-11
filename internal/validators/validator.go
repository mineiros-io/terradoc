package validators

import (
	"github.com/mineiros-io/terradoc/internal/entities"
	"github.com/mineiros-io/terradoc/internal/types"
)

type TypeMismatchResult struct {
	Name           string
	DefinedType    string
	DocumentedType string
}

type Summary struct {
	Type                 string
	MissingDefinition    []string
	MissingDocumentation []string
	TypeMismatch         []TypeMismatchResult
}

func (vs Summary) Success() bool {
	return len(vs.MissingDocumentation) == 0 &&
		len(vs.MissingDefinition) == 0 &&
		len(vs.TypeMismatch) == 0
}

func TypesMatch(typeA, typeB *entities.Type) bool {
	if typeA == nil && typeB == nil {
		return true
	}

	if typeA.TFType == types.TerraformAny || typeB.TFType == types.TerraformAny {
		if typeA.TFType == types.TerraformObject || typeB.TFType == types.TerraformObject {
			return true
		}

		if typeA.HasNestedType() {
			return TypesMatch(typeA.Nested, typeB)
		}

		if typeB.HasNestedType() {
			return TypesMatch(typeA, typeB.Nested)
		}
	}

	if typeA.TFType == types.TerraformObject && typeB.TFType == types.TerraformObject {
		return true
	}

	return (typeA.TFType == typeB.TFType) &&
		TypesMatch(typeA.Nested, typeB.Nested)
}

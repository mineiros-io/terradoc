package validators

import (
	"github.com/mineiros-io/terradoc/internal/entities"
	"github.com/mineiros-io/terradoc/internal/types"
)

func TypesMatch(typeA, typeB *entities.Type) bool {
	if typeA == nil && typeB == nil {
		return true
	}

	// TODO: terraform accepts `any` for object
	if typeA.TFType == types.TerraformObject && typeB.TFType == types.TerraformObject {
		return true
	}

	return (typeA == nil && typeB == nil) ||
		(typeA.Label == typeB.Label) &&
			(typeA.TFType == typeB.TFType) &&
			TypesMatch(typeA.Nested, typeB.Nested)
}

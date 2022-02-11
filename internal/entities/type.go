package entities

import (
	"fmt"

	"github.com/mineiros-io/terradoc/internal/types"
)

// Type represents a variable or attribute type with its readme and Terraform type data
type Type struct {
	// TFType is the specific Terraform type definition for this type
	TFType types.TerraformType `json:"type"`
	// Label is an optional label for the TerraformType
	Label string `json:"label"`
	// Nested is an optional nested type definition
	Nested *Type `json:"nested,omitempty"`
}

func (t Type) AsString() string {
	switch {
	case t.HasNestedType():
		return fmt.Sprintf("%s(%s)", t.TFType.String(), t.Nested.AsString())
	case t.Label != "":
		return fmt.Sprintf("%s(%s)", t.TFType.String(), t.Label)
	}

	return t.TFType.String()
}

func (t Type) HasNestedType() bool {
	return t.Nested != nil
}

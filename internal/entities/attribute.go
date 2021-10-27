package entities

import (
	"encoding/json"

	"github.com/mineiros-io/terradoc/internal/types"
)

// Attribute represents an `attribute` block from the input file
type Attribute struct {
	Name             string              `json:"name"`                     // Name is the attribute name as defined in the `attribute` block label
	TerraformType    types.TerraformType `json:"type"`                     // TerraformType is the Terraform type assigned to the variable
	ReadmeType       string              `json:"readme_type,omitempty"`    // ReadmeType is an optional type to be rendered on the README document. Can be any string value as opposed to the `Type` field
	Default          json.RawMessage     `json:"default,omitempty"`        // Default is an optional default value for this variable in case none is given. Must be a valid JSON value.
	Description      string              `json:"description,omitempty"`    // Description is an optional attribute description
	ForcesRecreation bool                `json:"forces_recreation"`        // ForcesRecreation specifies if a change in the attribute will force the resource recreation
	ReadmeExample    string              `json:"readme_example,omitempty"` // ReadmeExample is an optional readme example to be used in the documentation
	Required         bool                `json:"required"`                 // Required specifies if the attribute is required
	Attributes       []Attribute         `json:"attributes,omitempty"`     // Attributes is a collection of nested attributes contained in the attribute block definition
	Level            int                 `json:"-"`                        // Level is the nesting level of this attribute
}

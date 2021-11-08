package entities

import (
	"encoding/json"

	"github.com/mineiros-io/terradoc/internal/types"
)

// Variable represents a `variable` block from the input file.
type Variable struct {
	Name             string              `json:"name"`                        // Name as defined in the `variable` block label.
	TerraformType    types.TerraformType `json:"type"`                        // Type is the Terraform type assigned to the variable
	Description      string              `json:"description,omitempty"`       // Description is an optional variable description
	ReadmeType       string              `json:"readme_type,omitempty"`       // ReadmeType is an optional type to be rendered on the README document. Can be any string value as opposed to the `Type` field
	Default          json.RawMessage     `json:"default,omitempty"`           // Default is an optional default value for this variable in case none is given. Must be a valid JSON value.
	Required         bool                `json:"required,omitempty"`          // Required specifies if the variable is required
	ForcesRecreation bool                `json:"forces_recreation,omitempty"` // ForcesRecreation specifies if a change in the variable triggers the recreation of the resource.
	ReadmeExample    string              `json:"readme_example,omitempty"`    // ReadmeExample is an optional readme example to be used in the documentation
	Attributes       []Attribute         `json:"attributes,omitempty"`        // Attributes is a collection attributes that make up the value of this variable.
}

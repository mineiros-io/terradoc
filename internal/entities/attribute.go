package entities

import (
	"encoding/json"
)

// Attribute represents an `attribute` block from the input file
type Attribute struct {
	// Name is the attribute name as defined in the `attribute` block label
	Name string `json:"name"`
	// Type is the type definition for the attribute
	Type Type `json:"type_definition"`
	// Default is an optional default value for this variable in case none is given. Must be a valid JSON value.
	Default json.RawMessage `json:"default,omitempty"`
	// Description is an optional attribute description
	Description string `json:"description,omitempty"`
	// ForcesRecreation specifies if a change in the attribute will force the resource recreation
	ForcesRecreation bool `json:"forces_recreation"`
	// ReadmeExample is an optional readme example to be used in the documentation
	ReadmeExample string `json:"readme_example,omitempty"`
	// Required specifies if the attribute is required
	Required bool `json:"required"`
	// Attributes is a collection of nested attributes contained in the attribute block definition
	Attributes []Attribute `json:"attributes,omitempty"`
	// Level is the nesting level of this attribute
	Level int `json:"-"`
}

package entities

import (
	"encoding/json"
)

type VariableCollection []Variable

func (vc VariableCollection) VarByName(name string) (Variable, bool) {
	for _, v := range vc {
		if v.Name == name {
			return v, true
		}
	}

	return Variable{}, false
}

// Variable represents a `variable` block from the input file.
type Variable struct {
	// Name as defined in the `variable` block label.
	Name string `json:"name"`
	// Type is a type definition for the variable
	Type Type `json:"type_definition"`
	// Description is an optional variable description
	Description string `json:"description,omitempty"`
	// Default is an optional default value for this variable in case none is given. Must be a valid JSON value.
	Default json.RawMessage `json:"default,omitempty"`
	// Required specifies if the variable is required
	Required bool `json:"required,omitempty"`
	// ForcesRecreation specifies if a change in the variable triggers the recreation of the resource.
	ForcesRecreation bool `json:"forces_recreation,omitempty"`
	// ReadmeExample is an optional readme example to be used in the documentation
	ReadmeExample string `json:"readme_example,omitempty"`
	// Attributes is a collection attributes that make up the value of this variable.
	Attributes []Attribute `json:"attributes,omitempty"`
}

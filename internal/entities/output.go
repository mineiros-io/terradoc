package entities

// Output represents an `output` block from the input file.
type Output struct {
	// Name as defined in the `output` block label.
	Name string `json:"name"`
	// Type is a type definition for the output
	Type Type `json:"type_definition"`
	// Description is an optional output description
	Description string `json:"description,omitempty"`
}

type OutputCollection []Output

func (oc OutputCollection) OutputByName(name string) (Output, bool) {
	for _, o := range oc {
		if o.Name == name {
			return o, true
		}
	}

	return Output{}, false
}

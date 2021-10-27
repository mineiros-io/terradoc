package entities

// Section represents a `section` block from the input file.
type Section struct {
	// Title is the title of the section.
	Title string `json:"title"`
	// Description is an optional section description.
	Description string `json:"description,omitempty"`
	// Variables is a collection of variable definitions contained in the section block.
	Variables []Variable `json:"variables,omitempty"`
	// SubSections is a collection of nested sections contained in the section block.
	SubSections []Section `json:"subsections,omitempty"`
	// Level is the nesting of this section
	Level int `json:"-"`
}

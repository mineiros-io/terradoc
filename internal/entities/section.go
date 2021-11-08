package entities

// Section represents a `section` block from the input file.
type Section struct {
	Title       string     `json:"title"`                 // Title is the title of the section.
	Description string     `json:"description,omitempty"` // Description is an optional section description.
	Variables   []Variable `json:"variables,omitempty"`   // Variables is a collection of variable definitions contained in the section block.
	SubSections []Section  `json:"subsections,omitempty"` // SubSections is a collection of nested sections contained in the section block.
	Level       int        `json:"-"`                     // Level is the nesting of this section
}

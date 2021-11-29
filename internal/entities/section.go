package entities

// Section represents a `section` block from the input file.
type Section struct {
	// Title is an optional title for the section.
	Title string `json:"title"`
	// Content is an optional text content for the section.
	Content string `json:"content,omitempty"`
	// Variables is a collection of variable definitions contained in the section block.
	Variables []Variable `json:"variables,omitempty"`
	// SubSections is a collection of nested sections contained in the section block.
	SubSections []Section `json:"subsections,omitempty"`
	// Level is the nesting of this section
	Level int `json:"-"`
	// TOC is a flag for generating table of contents for nested sections
	TOC bool `json:"-"`
}

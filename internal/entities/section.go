package entities

type Section struct {
	Title       string      `json:"title,omitempty"`
	Description string      `json:"description,omitempty"`
	Variables   []*Variable `json:"variables,omitempty"`
	Sections    []*Section  `json:"sections,omitempty"`
	Level       int         `json:"-"`
}

package entities

type Attribute struct {
	Name             string       `json:"name"`
	Type             string       `json:"type"`
	Description      string       `json:"description,omitempty"`
	ForcesRecreation bool         `json:"forces_recreation"`
	Required         bool         `json:"required"`
	Attributes       []*Attribute `json:"attributes,omitempty"`
	Level            int          `json:"-"`
}

package entities

// Definition represents a parsed source file.
type Definition struct {
	Sections []Section `json:"sections"` // Sections is a collection of sections defined in the source file.
}

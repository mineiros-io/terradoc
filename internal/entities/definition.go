package entities

// Definition represents a parsed source file.
type Definition struct {
	// Sections is a collection of sections defined in the source file.
	Sections []Section `json:"sections"`
	// References is a collection of references defined in the source file.
	References []Reference `json:"references"`
}

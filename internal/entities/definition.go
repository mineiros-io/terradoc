package entities

// Definition represents a parsed source file.
type Definition struct {
	// Header is the header section block from the source file
	Header Header
	// Sections is a collection of root level section blocks from the source file
	Sections []Section
	// References is a collection of references as found in the source file
	References []Reference `json:"reference"`
}

// Reference represents a reference block on the parsed source fil
type Reference struct {
	Name  string `json:"name"`  // Name is the identifier for the reference
	Value string `json:"value"` // Value is the value for the reference
}

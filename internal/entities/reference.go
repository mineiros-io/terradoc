package entities

// Reference represents a `ref` block on the source document
type Reference struct {
	Name  string `json:"name"`  // Name is the identifier for the reference
	Value string `json:"value"` // Value is the value that the reference holds
}

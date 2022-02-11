package entities

// TFDoc represents a parsed source file.
type TFDoc struct {
	// Header is the header section block from the source file
	Header Header
	// Sections is a collection of sections defined in the source file.
	Sections []Section `json:"sections"`
	// References is a collection of references defined in the source file.
	References []Reference `json:"references"`
}

// Header represents the `header` block on the parsed source file
type Header struct {
	Image  string  `json:"image"`  // Image is the image url to be displayed on the header
	URL    string  `json:"url"`    // URL is the target url for the image
	Badges []Badge `json:"badges"` // Badges is a collection of Badge entities
}

// Badge represents a `badge` block on the parsed source file
type Badge struct {
	Image string `json:"image"` // Image is the badge's image url
	URL   string `json:"url"`   // URL is the target url for the badge
	Text  string `json:"text"`  // Text is the text label for the badge
	Name  string `json:"name"`  // Name is an identifier for the badge
}

func (d TFDoc) AllVariables() (result []Variable) {
	for _, s := range d.Sections {
		result = append(result, s.AllVariables()...)
	}

	return result
}

func (d TFDoc) AllOutputs() (result []Output) {
	for _, s := range d.Sections {
		result = append(result, s.AllOutputs()...)
	}

	return result
}

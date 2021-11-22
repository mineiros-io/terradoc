package entities

// Badge represents the `header` block on the parsed source file
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

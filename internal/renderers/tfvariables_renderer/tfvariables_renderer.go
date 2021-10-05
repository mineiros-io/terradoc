package tfvariables_renderer

import (
	"fmt"
	"html/template"
	"io"

	"github.com/mineiros-io/terradoc/internal/entities"
)

const (
	templateFile = "templates/variables.tf"
)

func Render(writer io.Writer, definition *entities.Definition) error {
	t, err := template.ParseFiles(templateFile)
	if err != nil {
		return fmt.Errorf("Error parsing template %q: %s", templateFile, err)
	}

	err = t.Execute(writer, definition)

	if err != nil {
		return err
	}

	return nil
}

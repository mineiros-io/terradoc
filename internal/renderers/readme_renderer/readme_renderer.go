package readme_renderer

import (
	"io"
	"text/template"

	"github.com/mineiros-io/terradoc/internal/entities"
	"github.com/mineiros-io/terradoc/internal/renderers"
)

const (
	defaultFilename = "README.md"
	templateFile    = "templates/readme.md"
	templateName    = "readme.md"
)

func Render(writer io.Writer, definition *entities.Definition) error {
	t := template.New(templateName)

	funcMap := renderers.TemplateFuncMap(t)

	t, err := t.Funcs(funcMap).ParseFiles(templateFile)
	if err != nil {
		return err
	}

	return t.Execute(writer, definition)
}

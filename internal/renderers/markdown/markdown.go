package markdown

import (
	"io"

	"github.com/mineiros-io/terradoc/internal/entities"
)

func Render(writer io.Writer, definition entities.Doc) error {
	mdWriter, err := newMarkdownWriter(writer)
	if err != nil {
		return err
	}

	return mdWriter.writeDefinition(definition)
}

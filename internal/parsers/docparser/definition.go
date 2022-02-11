package docparser

import (
	"fmt"

	"github.com/hashicorp/hcl/v2"
	"github.com/mineiros-io/terradoc/internal/entities"
	"github.com/mineiros-io/terradoc/internal/schemas/docschema"
)

func parseDoc(f *hcl.File) (entities.Doc, error) {
	docContent, diags := f.Body.Content(docschema.RootSchema())
	if diags.HasErrors() {
		return entities.Doc{}, fmt.Errorf("parsing Terradoc doc: %v", diags.Errs())
	}

	var err error

	def := entities.Doc{}

	def.Header, err = parseHeader(docContent.Blocks.OfType(headerBlockName))
	if err != nil {
		return entities.Doc{}, fmt.Errorf("parsing header: %v", err)
	}

	def.Sections, err = parseSections(docContent.Blocks.OfType(sectionBlockName))
	if err != nil {
		return entities.Doc{}, err
	}

	def.References, err = parseReferences(docContent.Blocks.OfType(referencesBlockName))
	if err != nil {
		return entities.Doc{}, err
	}

	return def, nil
}

package tfdocparser

import (
	"fmt"

	"github.com/hashicorp/hcl/v2"
	"github.com/mineiros-io/terradoc/internal/entities"
	"github.com/mineiros-io/terradoc/internal/schemas/tfdocschema"
)

func parseDefinition(f *hcl.File) (entities.Definition, error) {
	definitionContent, diags := f.Body.Content(tfdocschema.RootSchema())
	if diags.HasErrors() {
		return entities.Definition{}, fmt.Errorf("parsing Terradoc definition: %v", diags.Errs())
	}

	var err error

	def := entities.Definition{}

	def.Header, err = parseHeader(definitionContent.Blocks.OfType(headerBlockName))
	if err != nil {
		return entities.Definition{}, fmt.Errorf("parsing header: %v", err)
	}

	def.Sections, err = parseSections(definitionContent.Blocks.OfType(sectionBlockName))
	if err != nil {
		return entities.Definition{}, err
	}

	def.References, err = parseReferences(definitionContent.Blocks.OfType(referencesBlockName))
	if err != nil {
		return entities.Definition{}, err
	}

	return def, nil
}

package tfdocparser

import (
	"fmt"

	"github.com/hashicorp/hcl/v2"
	"github.com/mineiros-io/terradoc/internal/entities"
	"github.com/mineiros-io/terradoc/internal/parsers/hclparser"
	"github.com/mineiros-io/terradoc/internal/schemas/docschema"
)

func parseReferences(referencesBlocks hcl.Blocks) ([]entities.Reference, error) {
	switch {
	case len(referencesBlocks) == 0:
		return nil, nil
	case len(referencesBlocks) != 1:
		return nil,
			fmt.Errorf("parsing references: expected at most 1 `references` block but got %d instead", len(referencesBlocks))
	}

	referencesContent, diags := referencesBlocks[0].Body.Content(docschema.ReferencesSchema())
	if diags.HasErrors() {
		return nil, fmt.Errorf("parsing references: %v", diags.Errs())
	}

	return parseRefs(referencesContent.Blocks.OfType(refBlockName))
}

func parseRefs(refBlocks hcl.Blocks) (refs []entities.Reference, err error) {
	for _, refBlock := range refBlocks {
		ref, err := parseRef(refBlock)
		if err != nil {
			return nil, err
		}

		refs = append(refs, ref)
	}

	return refs, nil
}

func parseRef(refBlock *hcl.Block) (entities.Reference, error) {
	// reference blocks are required to have a label as defined in the schema
	name := refBlock.Labels[0]

	refContent, diags := refBlock.Body.Content(docschema.RefSchema())
	if diags.HasErrors() {
		return entities.Reference{}, fmt.Errorf("parsing Terradoc `references`: %v", diags.Errs())
	}

	value, err := hclparser.GetAttribute(refContent.Attributes, valueAttributeName).String()
	if err != nil {
		return entities.Reference{}, err
	}

	return entities.Reference{Name: name, Value: value}, nil
}

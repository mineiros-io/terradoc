package hclparser

import (
	"fmt"

	"github.com/hashicorp/hcl/v2"
	"github.com/mineiros-io/terradoc/internal/entities"
	"github.com/mineiros-io/terradoc/internal/parser/hclparser/hclschema"
)

func parseReferences(referencesBlocks []*hcl.Block) ([]entities.Reference, error) {
	switch {
	case len(referencesBlocks) == 0:
		return nil, nil
	case len(referencesBlocks) != 1:
		return nil,
			fmt.Errorf("parsing references: expected at most 1 `references` block but got %d instead", len(referencesBlocks))
	}

	referencesContent, diags := referencesBlocks[0].Body.Content(hclschema.ReferencesSchema())
	if diags.HasErrors() {
		return nil, fmt.Errorf("parsing references: %v", diags.Errs())
	}

	return parseRefs(referencesContent.Blocks.OfType(refBlockName))
}

func parseRefs(refBlocks []*hcl.Block) (refs []entities.Reference, err error) {
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
	// name
	name := refBlock.Labels[0]

	refContent, diags := refBlock.Body.Content(hclschema.RefSchema())
	if diags.HasErrors() {
		return entities.Reference{}, fmt.Errorf("parsing Terradoc `references`: %v", diags.Errs())
	}

	// value
	value, err := getAttribute(refContent.Attributes, valueAttributeName).String()
	if err != nil {
		return entities.Reference{}, err
	}

	return entities.Reference{Name: name, Value: value}, nil
}

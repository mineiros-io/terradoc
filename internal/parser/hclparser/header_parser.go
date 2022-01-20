package hclparser

import (
	"fmt"

	"github.com/hashicorp/hcl/v2"
	"github.com/mineiros-io/terradoc/internal/entities"
	"github.com/mineiros-io/terradoc/internal/parser/hclparser/hclschema"
)

func parseHeader(headerBlocks hcl.Blocks) (entities.Header, error) {
	switch {
	case len(headerBlocks) == 0:
		return entities.Header{}, nil
	case len(headerBlocks) > 1:
		return entities.Header{}, fmt.Errorf("parsing Terradoc header: expected 1 but documenta has %d header blocks", len(headerBlocks))
	}

	headerBlock := headerBlocks[0]

	headerContent, diags := headerBlock.Body.Content(hclschema.HeaderSchema())
	if diags.HasErrors() {
		return entities.Header{}, fmt.Errorf("parsing Terradoc header: %v", diags.Errs())
	}

	header, err := createHeaderFromHCLAttributes(headerContent.Attributes)
	if err != nil {
		return entities.Header{}, fmt.Errorf("parsing header: %s", err)
	}

	// parse `badge` blocks
	for _, badgeBlk := range headerContent.Blocks.OfType(badgeBlockName) {
		badge, err := parseBadge(badgeBlk)
		if err != nil {
			return entities.Header{}, fmt.Errorf("parsing header badge: %s", err)
		}

		header.Badges = append(header.Badges, badge)
	}

	return header, nil
}

func parseBadge(badgeBlock *hcl.Block) (entities.Badge, error) {
	// badge blocks are required to have a label as defined in the schema
	name := badgeBlock.Labels[0]

	badgeContent, diags := badgeBlock.Body.Content(hclschema.BadgeSchema())
	if diags.HasErrors() {
		return entities.Badge{}, fmt.Errorf("parsing badge: %v", diags.Errs())
	}

	return createBadgeFromHCLAttributes(badgeContent.Attributes, name)
}

func createHeaderFromHCLAttributes(attrs hcl.Attributes) (entities.Header, error) {
	header := entities.Header{}

	// image
	image, err := getAttribute(attrs, imageAttributeName).String()
	if err != nil {
		return entities.Header{}, err
	}
	header.Image = image

	// url
	url, err := getAttribute(attrs, urlAttributeName).String()
	if err != nil {
		return entities.Header{}, err
	}
	header.URL = url

	return header, nil
}

func createBadgeFromHCLAttributes(attrs hcl.Attributes, name string) (entities.Badge, error) {
	badge := entities.Badge{Name: name}

	// image
	image, err := getAttribute(attrs, imageAttributeName).String()
	if err != nil {
		return entities.Badge{}, err
	}
	badge.Image = image

	// url
	url, err := getAttribute(attrs, urlAttributeName).String()
	if err != nil {
		return entities.Badge{}, err
	}
	badge.URL = url

	// url
	text, err := getAttribute(attrs, textAttributeName).String()
	if err != nil {
		return entities.Badge{}, err
	}
	badge.Text = text

	return badge, nil
}

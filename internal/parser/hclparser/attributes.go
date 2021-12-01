package hclparser

import (
	"fmt"

	"github.com/hashicorp/hcl/v2"
	"github.com/mineiros-io/terradoc/internal/entities"
	"github.com/mineiros-io/terradoc/internal/parser/hclparser/hclschema"
)

func parseVariableAttributes(attributeBlocks []*hcl.Block) (attributes []entities.Attribute, err error) {
	const variableAttributeLevel = 1

	for _, attrBlk := range attributeBlocks {
		attribute, err := parseAttribute(attrBlk, variableAttributeLevel)
		if err != nil {
			return nil, fmt.Errorf("parsing attributes: %s", err)
		}

		attributes = append(attributes, attribute)
	}

	return attributes, nil
}

func parseAttribute(attrBlock *hcl.Block, level int) (entities.Attribute, error) {
	attrContent, diags := attrBlock.Body.Content(hclschema.AttributeSchema())
	if diags.HasErrors() {
		return entities.Attribute{}, fmt.Errorf("parsing attribute block: %v", diags.Errs())
	}

	if len(attrBlock.Labels) != 1 {
		return entities.Attribute{}, fmt.Errorf("expected single 'name' label, got %v", attrBlock.Labels)
	}

	name := attrBlock.Labels[0]
	attr, err := createAttributeFromHCLAttributes(attrContent.Attributes, name, level)
	if err != nil {
		return entities.Attribute{}, fmt.Errorf("parsing attribute: %s", err)
	}

	nestedAttributeLevel := level + 1
	// attribute blocks have only `attribute` blocks
	for _, blk := range attrContent.Blocks.OfType(attributeBlockName) {
		nestedAttr, err := parseAttribute(blk, nestedAttributeLevel)
		if err != nil {
			return entities.Attribute{}, fmt.Errorf("parsing nested attribute: %s", err)
		}

		attr.Attributes = append(attr.Attributes, nestedAttr)
	}

	return attr, nil
}

func createAttributeFromHCLAttributes(attrs hcl.Attributes, name string, level int) (entities.Attribute, error) {
	var err error

	attr := entities.Attribute{Name: name, Level: level}

	attr.Description, err = getAttribute(attrs, descriptionAttributeName).String()
	if err != nil {
		return entities.Attribute{}, err
	}

	attr.Required, err = getAttribute(attrs, requiredAttributeName).Bool()
	if err != nil {
		return entities.Attribute{}, err
	}

	attr.ForcesRecreation, err = getAttribute(attrs, forcesRecreationAttributeName).Bool()
	if err != nil {
		return entities.Attribute{}, err
	}

	attr.ReadmeExample, err = getAttribute(attrs, readmeExampleAttributeName).String()
	if err != nil {
		return entities.Attribute{}, err
	}

	attr.Type, err = getType(attrs, name)
	if err != nil {
		return entities.Attribute{}, err
	}

	attr.Default, err = getAttribute(attrs, defaultAttributeName).RawJSON()
	if err != nil {
		return entities.Attribute{}, err
	}

	return attr, nil
}
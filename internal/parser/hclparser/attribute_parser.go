package hclparser

import (
	"errors"
	"fmt"

	"github.com/hashicorp/hcl/v2"
	"github.com/mineiros-io/terradoc/internal/entities"
	"github.com/mineiros-io/terradoc/internal/parser/hclparser/hclschema"
)

func parseAttribute(attrBlock *hcl.Block, level int) (entities.Attribute, error) {
	attrContent, diags := attrBlock.Body.Content(hclschema.AttributeSchema())
	if diags.HasErrors() {
		return entities.Attribute{}, fmt.Errorf("parsing attribute block: %v", diags.Errs())
	}

	if len(attrBlock.Labels) != 1 {
		return entities.Attribute{}, errors.New("attribute block does not have a name")
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
	attribute := entities.Attribute{Name: name, Level: level}

	// description
	description, err := getAttribute(attrs, descriptionAttributeName).String()
	if err != nil {
		return entities.Attribute{}, err
	}
	attribute.Description = description

	// required
	required, err := getAttribute(attrs, requiredAttributeName).Bool()
	if err != nil {
		return entities.Attribute{}, err
	}
	attribute.Required = required

	// forcesRecreation
	forcesRecreation, err := getAttribute(attrs, forcesRecreationAttributeName).Bool()
	if err != nil {
		return entities.Attribute{}, err
	}
	attribute.ForcesRecreation = forcesRecreation

	// readme example
	readmeExample, err := getAttribute(attrs, readmeExampleAttributeName).HCLString()
	if err != nil {
		return entities.Attribute{}, err
	}
	attribute.ReadmeExample = readmeExample

	// type definition
	typeDefinition, err := getType(attrs, name)
	if err != nil {
		return entities.Attribute{}, err
	}
	attribute.Type = typeDefinition

	// default
	attrDefault, err := getAttribute(attrs, defaultAttributeName).RawJSON()
	if err != nil {
		return entities.Attribute{}, err
	}
	attribute.Default = attrDefault

	return attribute, nil
}

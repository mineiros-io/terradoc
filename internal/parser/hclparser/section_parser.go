package hclparser

import (
	"fmt"

	"github.com/hashicorp/hcl/v2"
	"github.com/mineiros-io/terradoc/internal/entities"
	"github.com/mineiros-io/terradoc/internal/parser/hclparser/hclschema"
)

func parseRootSection(sectionBlock *hcl.Block) (entities.Section, error) {
	sectionContent, diags := sectionBlock.Body.Content(hclschema.RootSectionSchema())
	if diags.HasErrors() {
		return entities.Section{}, fmt.Errorf("parsing root section: %v", diags.Errs())
	}

	section := entities.Section{}

	// title
	title, err := getAttribute(sectionContent.Attributes, nameAttributeName).String()
	if err != nil {
		return entities.Section{}, err
	}
	section.Title = title

	// description
	description, err := getAttribute(sectionContent.Attributes, descriptionAttributeName).String()
	if err != nil {
		return entities.Section{}, err
	}
	section.Description = description

	for _, sectionBlock := range sectionContent.Blocks.OfType(sectionBlockName) {
		section, err := parseSection(sectionBlock, rootSectionLevel)
		if err != nil {
			return entities.Section{}, err
		}

		section.SubSections = append(section.SubSections, section)
	}

	return section, nil
}

func parseSection(sectionBlock *hcl.Block, level int) (entities.Section, error) {
	sectionContent, diags := sectionBlock.Body.Content(hclschema.SectionSchema())
	if diags.HasErrors() {
		return entities.Section{}, fmt.Errorf("parsing Terradoc section: %v", diags.Errs())
	}

	section, err := createSectionFromAttributes(sectionContent.Attributes, level)
	if err != nil {
		return entities.Section{}, fmt.Errorf("parsing section: %s", err)
	}

	// parse `variable` blocks
	for _, varBlk := range sectionContent.Blocks.OfType(variableBlockName) {
		variable, err := parseVariable(varBlk)
		if err != nil {
			return entities.Section{}, fmt.Errorf("parsing section variable: %s", err)
		}

		section.Variables = append(section.Variables, variable)
	}

	subSectionLevel := level + 1
	// parse `section` blocks
	for _, subSectionBlk := range sectionContent.Blocks.OfType(sectionBlockName) {
		subSection, err := parseSection(subSectionBlk, subSectionLevel)
		if err != nil {
			return entities.Section{}, fmt.Errorf("parsing subsection: %s", err)
		}

		section.SubSections = append(section.SubSections, subSection)
	}

	return section, nil
}

func createSectionFromAttributes(attrs hcl.Attributes, level int) (entities.Section, error) {
	section := entities.Section{Level: level}

	// title
	title, err := getAttribute(attrs, titleAttributeName).String()
	if err != nil {
		return entities.Section{}, err
	}
	section.Title = title

	// fetch section description
	description, err := getAttribute(attrs, descriptionAttributeName).String()
	if err != nil {
		return entities.Section{}, err
	}

	section.Description = description

	return section, nil
}

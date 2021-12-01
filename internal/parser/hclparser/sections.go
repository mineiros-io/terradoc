package hclparser

import (
	"fmt"

	"github.com/hashicorp/hcl/v2"
	"github.com/mineiros-io/terradoc/internal/entities"
	"github.com/mineiros-io/terradoc/internal/parser/hclparser/hclschema"
)

const (
	rootSectionLevel = 1
)

func parseSections(sectionBlocks []*hcl.Block) (sections []entities.Section, err error) {
	for _, sectionBlock := range sectionBlocks {
		section, err := parseSection(sectionBlock, rootSectionLevel) // initial level
		if err != nil {
			return nil, fmt.Errorf("parsing sections: %s", err)
		}

		sections = append(sections, section)
	}

	return sections, nil
}

func parseSection(sectionBlock *hcl.Block, level int) (entities.Section, error) {
	sectionContent, diags := sectionBlock.Body.Content(hclschema.SectionSchema())
	if diags.HasErrors() {
		return entities.Section{}, fmt.Errorf("parsing Terradoc section: %v", diags.Errs())
	}

	section, err := createSectionFromHCLAttributes(sectionContent.Attributes, level)
	if err != nil {
		return entities.Section{}, fmt.Errorf("parsing section: %s", err)
	}

	// parse `variable` blocks
	variables, err := parseVariables(sectionContent.Blocks.OfType(variableBlockName))
	if err != nil {
		return entities.Section{}, fmt.Errorf("parsing section variable: %v", err)
	}
	section.Variables = variables

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

func createSectionFromHCLAttributes(attrs hcl.Attributes, level int) (entities.Section, error) {
	var err error

	section := entities.Section{Level: level}

	section.Title, err = getAttribute(attrs, titleAttributeName).String()
	if err != nil {
		return entities.Section{}, err
	}

	section.Content, err = getAttribute(attrs, contentAttributeName).String()
	if err != nil {
		return entities.Section{}, err
	}

	section.TOC, err = getAttribute(attrs, tocAttributeName).Bool()
	if err != nil {
		return entities.Section{}, err
	}

	return section, nil
}

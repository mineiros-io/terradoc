package hclparser

import (
	"errors"
	"fmt"
	"io"

	"github.com/hashicorp/hcl/v2"
	"github.com/hashicorp/hcl/v2/hclparse"
	"github.com/mineiros-io/terradoc/internal/entities"
	"github.com/mineiros-io/terradoc/internal/parser/hclparser/hclschema"
)

const (
	titleAttributeName            = "title"
	descriptionAttributeName      = "description"
	typeAttributeName             = "type"
	readmeTypeAttributeName       = "readme_type"
	defaultAttributeName          = "default"
	requiredAttributeName         = "required"
	forcesRecreationAttributeName = "forces_recreation"
	readmeExampleAttributeName    = "readme_example"

	sectionBlockName   = "section"
	variableBlockName  = "variable"
	attributeBlockName = "attribute"

	variableAttributeLevel = 1
	rootSubSectionLevel    = 1
)

// Parse reads the content of a io.Reader and returns a Definition entity from its parsed values
func Parse(r io.Reader, filename string) (*entities.Definition, error) {
	src, err := io.ReadAll(r)
	if err != nil {
		return nil, err
	}

	return parseHCL(src, filename)
}

func parseHCL(src []byte, filename string) (*entities.Definition, error) {
	p := hclparse.NewParser()

	f, diags := p.ParseHCL(src, filename)
	if diags.HasErrors() {
		return nil, fmt.Errorf("Error parsing HCL: %v", diags.Errs())
	}

	return parseDefinition(f)
}

func parseDefinition(f *hcl.File) (*entities.Definition, error) {
	definitionContent, diags := f.Body.Content(hclschema.RootSchema())
	if diags.HasErrors() {
		return nil, fmt.Errorf("Error parsing Terradoc definition: %v", diags.Errs())
	}

	def := &entities.Definition{}

	// Root does not have attributes and only has `section` blocks
	for _, sectionBlock := range definitionContent.Blocks {
		section, err := parseSection(sectionBlock, rootSubSectionLevel) // initial level
		if err != nil {
			return nil, fmt.Errorf("Error parseing sections: %s", err)
		}

		def.Sections = append(def.Sections, *section)
	}

	return def, nil
}

func parseSection(sectionBlock *hcl.Block, level int) (*entities.Section, error) {
	sectionContent, diags := sectionBlock.Body.Content(hclschema.SectionSchema())
	if diags.HasErrors() {
		return nil, fmt.Errorf("Error parsing Terradoc section: %v", diags.Errs())
	}

	section, err := createSectionFromAttributes(sectionContent.Attributes, level)
	if err != nil {
		return nil, fmt.Errorf("Error parsing section: %s", err)
	}

	// parse `variable` blocks
	for _, varBlk := range sectionContent.Blocks.OfType(variableBlockName) {
		variable, err := parseVariable(varBlk)
		if err != nil {
			return nil, fmt.Errorf("Error parsing section variable: %s", err)
		}

		section.Variables = append(section.Variables, *variable)
	}

	subSectionLevel := level + 1
	// parse `section` blocks
	for _, subSectionBlk := range sectionContent.Blocks.OfType(sectionBlockName) {
		subSection, err := parseSection(subSectionBlk, subSectionLevel)
		if err != nil {
			return nil, fmt.Errorf("Error parsing subsection: %s", err)
		}

		section.SubSections = append(section.SubSections, *subSection)
	}

	return section, nil
}

func parseVariable(variableBlock *hcl.Block) (*entities.Variable, error) {
	variableContent, diags := variableBlock.Body.Content(hclschema.VariableSchema())
	if diags.HasErrors() {
		return nil, fmt.Errorf("Error parsing variable: %v", diags.Errs())
	}

	if len(variableBlock.Labels) != 1 {
		return nil, errors.New("Attribute block does not have a name")
	}

	name := variableBlock.Labels[0]
	variable, err := createVariableFromAttributes(variableContent.Attributes, name)
	if err != nil {
		return nil, fmt.Errorf("Error parsing variable: %s", err)
	}

	// variables have only `attribute` blocks
	for _, blk := range variableContent.Blocks.OfType(attributeBlockName) {
		attribute, err := parseAttribute(blk, variableAttributeLevel)
		if err != nil {
			return nil, fmt.Errorf("Error parsing variable attributes: %s", err)
		}

		variable.Attributes = append(variable.Attributes, *attribute)
	}

	return variable, nil
}

func parseAttribute(attrBlock *hcl.Block, level int) (*entities.Attribute, error) {
	attrContent, diags := attrBlock.Body.Content(hclschema.AttributeSchema())
	if diags.HasErrors() {
		return nil, fmt.Errorf("Error parsing attribute block: %v", diags.Errs())
	}

	if len(attrBlock.Labels) != 1 {
		return nil, errors.New("Attribute block does not have a name")
	}

	name := attrBlock.Labels[0]
	attr, err := createAttributeFromAttributes(attrContent.Attributes, name, level)
	if err != nil {
		return nil, fmt.Errorf("Error parsing attribute: %s", err)
	}

	nestedAttributeLevel := level + 1
	// attribute blocks have only `attribute` blocks
	for _, blk := range attrContent.Blocks.OfType(attributeBlockName) {
		nestedAttr, err := parseAttribute(blk, nestedAttributeLevel)
		if err != nil {
			return nil, fmt.Errorf("Error parsing nested attribute: %s", err)
		}

		attr.Attributes = append(attr.Attributes, *nestedAttr)
	}

	return attr, nil
}

func createSectionFromAttributes(attrs hcl.Attributes, level int) (*entities.Section, error) {
	section := &entities.Section{Level: level}

	// title
	title, err := getAttribute(attrs, titleAttributeName).String()
	if err != nil {
		return nil, err
	}
	section.Title = title

	// fetch section description
	description, err := getAttribute(attrs, descriptionAttributeName).String()
	if err != nil {
		return nil, err
	}

	section.Description = description

	return section, nil
}

func createVariableFromAttributes(attrs hcl.Attributes, name string) (*entities.Variable, error) {
	variable := &entities.Variable{Name: name}

	// type
	varType, err := getAttribute(attrs, typeAttributeName).TerraformType()
	if err != nil {
		return nil, err
	}
	variable.TerraformType = varType

	// description
	description, err := getAttribute(attrs, descriptionAttributeName).String()
	if err != nil {
		return nil, err
	}
	variable.Description = description

	// readme type
	readmeType, err := getAttribute(attrs, readmeTypeAttributeName).String()
	if err != nil {
		return nil, err
	}
	variable.ReadmeType = readmeType

	// default
	varDefault, err := getAttribute(attrs, defaultAttributeName).RawJSON()
	if err != nil {
		return nil, err
	}
	variable.Default = varDefault

	// required
	required, err := getAttribute(attrs, requiredAttributeName).Bool()
	if err != nil {
		return nil, err
	}
	variable.Required = required

	// forcesRecreation
	forcesRecreation, err := getAttribute(attrs, forcesRecreationAttributeName).Bool()
	if err != nil {
		return nil, err
	}
	variable.ForcesRecreation = forcesRecreation

	// readme example
	readmeExample, err := getAttribute(attrs, readmeExampleAttributeName).HCLString()
	if err != nil {
		return nil, err
	}
	variable.ReadmeExample = readmeExample

	return variable, nil
}

func createAttributeFromAttributes(attrs hcl.Attributes, name string, level int) (*entities.Attribute, error) {
	attribute := &entities.Attribute{Name: name, Level: level}

	// type
	attrType, err := getAttribute(attrs, typeAttributeName).TerraformType()
	if err != nil {
		return nil, err
	}
	attribute.TerraformType = attrType

	// description
	description, err := getAttribute(attrs, descriptionAttributeName).String()
	if err != nil {
		return nil, err
	}
	attribute.Description = description

	// required
	required, err := getAttribute(attrs, requiredAttributeName).Bool()
	if err != nil {
		return nil, err
	}
	attribute.Required = required

	// forcesRecreation
	forcesRecreation, err := getAttribute(attrs, forcesRecreationAttributeName).Bool()
	if err != nil {
		return nil, err
	}
	attribute.ForcesRecreation = forcesRecreation

	// readme example
	readmeExample, err := getAttribute(attrs, readmeExampleAttributeName).HCLString()
	if err != nil {
		return nil, err
	}
	attribute.ReadmeExample = readmeExample

	// readme type
	readmeType, err := getAttribute(attrs, readmeTypeAttributeName).String()
	if err != nil {
		return nil, err
	}
	attribute.ReadmeType = readmeType

	// default
	varDefault, err := getAttribute(attrs, defaultAttributeName).RawJSON()
	if err != nil {
		return nil, err
	}
	attribute.Default = varDefault

	return attribute, nil
}

func getAttribute(attrs hcl.Attributes, name string) *hclAttribute {
	attr, exists := attrs[name]
	if exists {
		return &hclAttribute{attr}
	}

	return &hclAttribute{}
}

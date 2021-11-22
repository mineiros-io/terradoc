package hclparser

import (
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
	imageAttributeName            = "image"
	urlAttributeName              = "url"
	textAttributeName             = "text"
	valueAttributeName            = "value"
	nameAttributeName             = "name"

	sectionBlockName    = "section"
	referencesBlockName = "references"
	headerBlockName     = "header"
	variableBlockName   = "variable"
	attributeBlockName  = "attribute"
	badgeBlockName      = "badge"
	referenceBlockName  = "reference"

	variableAttributeLevel = 1
	rootSectionLevel       = 2
)

// Parse reads the content of a io.Reader and returns a Definition entity from its parsed values
func Parse(r io.Reader, filename string) (entities.Definition, error) {
	src, err := io.ReadAll(r)
	if err != nil {
		return entities.Definition{}, err
	}

	return parseHCL(src, filename)
}

func parseHCL(src []byte, filename string) (entities.Definition, error) {
	p := hclparse.NewParser()

	f, diags := p.ParseHCL(src, filename)
	if diags.HasErrors() {
		return entities.Definition{}, fmt.Errorf("parsing HCL: %v", diags.Errs())
	}

	return parseDefinition(f)
}

func parseDefinition(f *hcl.File) (entities.Definition, error) {
	def := entities.Definition{}

	definitionContent, diags := f.Body.Content(hclschema.RootSchema())
	if diags.HasErrors() {
		return entities.Definition{}, fmt.Errorf("parsing Terradoc definition: %v", diags.Errs())
	}

	// parse header

	for _, headerBlock := range definitionContent.Blocks.OfType(headerBlockName) {
		header, err := parseHeader(headerBlock) // initial level
		if err != nil {
			return entities.Definition{}, fmt.Errorf("parsing sections: %s", err)
		}

		def.Header = header
	}

	// parse sections
	for _, sectionBlock := range definitionContent.Blocks.OfType(sectionBlockName) {
		section, err := parseSection(sectionBlock, rootSectionLevel)
		if err != nil {
			return entities.Definition{}, fmt.Errorf("parsing sections: %s", err)
		}

		def.Sections = append(def.Sections, section)
	}

	for _, referenceBlock := range definitionContent.Blocks.OfType(referencesBlockName) {
		references, err := parseReferences(referenceBlock)
		if err != nil {
			return entities.Definition{}, fmt.Errorf("parsing sections: %s", err)
		}

		def.References = references
	}

	return def, nil
}

func parseReferences(referenceBlock *hcl.Block) (references []entities.Reference, err error) {
	referencesContent, diags := referenceBlock.Body.Content(hclschema.ReferencesSchema())
	if diags.HasErrors() {
		return nil, fmt.Errorf("parsing Terradoc section: %v", diags.Errs())
	}

	for _, referenceBlock := range referencesContent.Blocks.OfType("reference") {
		name := referenceBlock.Labels[0]

		reference := entities.Reference{Name: name}

		referenceContent, diags := referenceBlock.Body.Content(hclschema.ReferenceSchema())
		if diags.HasErrors() {
			return nil, fmt.Errorf("parsing Terradoc section: %v", diags.Errs())
		}

		value, err := getAttribute(referenceContent.Attributes, valueAttributeName).String()
		if err != nil {
			return nil, err
		}
		reference.Value = value

		references = append(references, reference)
	}

	return references, nil
}

func getAttribute(attrs hcl.Attributes, name string) *hclAttribute {
	attr, exists := attrs[name]
	if exists {
		return &hclAttribute{attr}
	}

	return &hclAttribute{}
}

func getType(attrs hcl.Attributes, attrName string) (entities.Type, error) {
	readmeType, err := getAttribute(attrs, readmeTypeAttributeName).String()
	if err != nil {
		return entities.Type{}, err
	}

	terraformType, err := getAttribute(attrs, typeAttributeName).TerraformType()
	if err != nil {
		return entities.Type{}, err
	}

	return entities.Type{
		TerraformType: terraformType,
		ReadmeType:    readmeType,
		Name:          attrName,
	}, nil
}

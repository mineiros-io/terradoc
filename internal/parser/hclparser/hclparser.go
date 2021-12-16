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
	contentAttributeName          = "content"
	descriptionAttributeName      = "description"
	typeAttributeName             = "type"
	readmeTypeAttributeName       = "readmeType"
	defaultAttributeName          = "default"
	requiredAttributeName         = "required"
	forcesRecreationAttributeName = "forces_recreation"
	readmeExampleAttributeName    = "readme_example"
	valueAttributeName            = "value"
	imageAttributeName            = "image"
	urlAttributeName              = "url"
	textAttributeName             = "text"
	tocAttributeName              = "toc"

	sectionBlockName    = "section"
	variableBlockName   = "variable"
	attributeBlockName  = "attribute"
	referencesBlockName = "references"
	refBlockName        = "ref"
	headerBlockName     = "header"
	badgeBlockName      = "badge"
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

	// header
	header, err := parseHeader(definitionContent.Blocks.OfType(headerBlockName))
	if err != nil {
		return entities.Definition{}, fmt.Errorf("parsing header: %v", err)
	}
	def.Header = header

	// sections
	sections, err := parseSections(definitionContent.Blocks.OfType(sectionBlockName))
	if err != nil {
		return entities.Definition{}, err
	}
	def.Sections = sections

	def.References, err = parseReferences(definitionContent.Blocks.OfType(referencesBlockName))
	if err != nil {
		return entities.Definition{}, err
	}

	return def, nil
}

func getAttribute(attrs hcl.Attributes, name string) *hclAttribute {
	attr, exists := attrs[name]
	if exists {
		return &hclAttribute{attr}
	}

	return &hclAttribute{}
}

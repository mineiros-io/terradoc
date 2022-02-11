package docparser

import (
	"fmt"
	"io"

	"github.com/hashicorp/hcl/v2/hclparse"
	"github.com/mineiros-io/terradoc/internal/entities"
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
	outputBlockName     = "output"
)

// Parse reads the content of a io.Reader and returns a Definition entity from its parsed values
func Parse(r io.Reader, filename string) (entities.Doc, error) {
	src, err := io.ReadAll(r)
	if err != nil {
		return entities.Doc{}, err
	}

	return parseHCL(src, filename)
}

func parseHCL(src []byte, filename string) (entities.Doc, error) {
	p := hclparse.NewParser()

	f, diags := p.ParseHCL(src, filename)
	if diags.HasErrors() {
		return entities.Doc{}, fmt.Errorf("parsing HCL: %v", diags.Errs())
	}

	return parseDoc(f)
}

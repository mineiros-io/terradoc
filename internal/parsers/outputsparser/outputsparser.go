package outputsparser

import (
	"errors"
	"fmt"
	"io"

	"github.com/hashicorp/hcl/v2"
	"github.com/hashicorp/hcl/v2/hclparse"
	"github.com/mineiros-io/terradoc/internal/entities"
	"github.com/mineiros-io/terradoc/internal/parsers/hclparser"
	"github.com/mineiros-io/terradoc/internal/schemas/outputsschema"
)

func Parse(r io.Reader, filename string) (entities.OutputsFile, error) {
	src, err := io.ReadAll(r)
	if err != nil {
		return entities.OutputsFile{}, err
	}

	return parseOutputsHCL(src, filename)
}

func parseOutputsHCL(src []byte, filename string) (entities.OutputsFile, error) {
	p := hclparse.NewParser()

	f, diags := p.ParseHCL(src, filename)
	if diags.HasErrors() {
		return entities.OutputsFile{}, fmt.Errorf("parsing HCL: %v", diags.Errs())
	}

	content, diags := f.Body.Content(outputsschema.RootSchema())
	if diags.HasErrors() {
		return entities.OutputsFile{}, fmt.Errorf("getting body content: %v", diags.Errs())
	}

	outputs, err := parseOutputs(content.Blocks.OfType("output"))
	if err != nil {
		return entities.OutputsFile{}, fmt.Errorf("parsing outputs: %v", err)
	}

	return entities.OutputsFile{Outputs: outputs}, nil
}

func parseOutputs(outputBlocks hcl.Blocks) (outputs []entities.Output, err error) {
	for _, outBlk := range outputBlocks {
		output, err := parseOutput(outBlk)
		if err != nil {
			return nil, fmt.Errorf("parsing output: %s", err)
		}

		outputs = append(outputs, output)
	}

	return outputs, nil
}

func parseOutput(outputBlock *hcl.Block) (entities.Output, error) {
	if len(outputBlock.Labels) != 1 {
		return entities.Output{}, errors.New("output block must have a single label")
	}

	outputContent, diags := outputBlock.Body.Content(outputsschema.OutputSchema())
	if diags.HasErrors() {
		return entities.Output{}, fmt.Errorf("parsing output: %v", diags.Errs())
	}

	// output blocks are required to have a label as defined in the schema
	name := outputBlock.Labels[0]
	output, err := createOutputFromHCLAttributes(outputContent.Attributes, name)
	if err != nil {
		return entities.Output{}, fmt.Errorf("parsing output: %s", err)
	}

	return output, nil
}

func createOutputFromHCLAttributes(attrs hcl.Attributes, name string) (entities.Output, error) {
	var err error
	output := entities.Output{Name: name}

	// description
	output.Description, err = hclparser.GetAttribute(attrs, "description").String()
	if err != nil {
		return entities.Output{}, err
	}

	return output, nil
}

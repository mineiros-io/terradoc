package docparser

import (
	"errors"
	"fmt"

	"github.com/hashicorp/hcl/v2"
	"github.com/mineiros-io/terradoc/internal/entities"
	"github.com/mineiros-io/terradoc/internal/parsers/hclparser"
	"github.com/mineiros-io/terradoc/internal/schemas/docschema"
)

func parseOutputs(outputBlocks hcl.Blocks) (outputs []entities.Output, err error) {
	for _, outputBlk := range outputBlocks {
		output, err := parseOutput(outputBlk)
		if err != nil {
			return nil, fmt.Errorf("parsing output: %s", err)
		}

		outputs = append(outputs, output)
	}

	return outputs, nil
}

func parseOutput(outputBlock *hcl.Block) (entities.Output, error) {
	if len(outputBlock.Labels) != 1 {
		return entities.Output{}, errors.New("output block does not have a name")
	}

	outputContent, diags := outputBlock.Body.Content(docschema.OutputSchema())
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

	output.Description, err = hclparser.GetAttribute(attrs, descriptionAttributeName).String()
	if err != nil {
		return entities.Output{}, err
	}

	// type definition
	output.Type, err = hclparser.GetAttribute(attrs, typeAttributeName).OutputType()
	if err != nil {
		return entities.Output{}, err
	}

	return output, nil
}

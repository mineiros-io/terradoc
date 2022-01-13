package hclparser

import (
	"errors"
	"fmt"

	"github.com/hashicorp/hcl/v2"
	"github.com/mineiros-io/terradoc/internal/entities"
	"github.com/mineiros-io/terradoc/internal/parser/hclparser/hclschema"
)

func parseOutputs(outputBlocks []*hcl.Block) (outputs []entities.Output, err error) {
	for _, outputBlk := range outputBlocks {
		variable, err := parseOutput(outputBlk)
		if err != nil {
			return nil, fmt.Errorf("parsing variable: %s", err)
		}

		outputs = append(outputs, variable)
	}

	return outputs, nil
}

func parseOutput(outputBlock *hcl.Block) (entities.Output, error) {
	if len(outputBlock.Labels) != 1 {
		return entities.Output{}, errors.New("variable block does not have a name")
	}

	outputContent, diags := outputBlock.Body.Content(hclschema.OutputSchema())
	if diags.HasErrors() {
		return entities.Output{}, fmt.Errorf("parsing variable: %v", diags.Errs())
	}

	// variables have only the `name` label
	name := outputBlock.Labels[0]
	output, err := createOutputFromHCLAttributes(outputContent.Attributes, name)
	if err != nil {
		return entities.Output{}, fmt.Errorf("parsing variable: %s", err)
	}

	return output, nil
}

func createOutputFromHCLAttributes(attrs hcl.Attributes, name string) (entities.Output, error) {
	var err error

	output := entities.Output{Name: name}

	output.Description, err = getAttribute(attrs, descriptionAttributeName).String()
	if err != nil {
		return entities.Output{}, err
	}

	// type definition
	output.Type, err = getAttribute(attrs, typeAttributeName).Type()
	if err != nil {
		return entities.Output{}, err
	}

	return output, nil
}

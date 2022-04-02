package validationparser

import (
	"errors"
	"fmt"
	"io"
	"strings"

	"github.com/hashicorp/hcl/v2"
	"github.com/hashicorp/hcl/v2/hclparse"
	"github.com/mineiros-io/terradoc/internal/entities"
	"github.com/mineiros-io/terradoc/internal/parsers/hclparser"
	"github.com/mineiros-io/terradoc/internal/schemas/validationschema"
	"github.com/mineiros-io/terradoc/internal/schemas/varsschema"
)

func Parse(r io.Reader, filename string, variablesEnabled bool, outputsEnabled bool) (entities.ValidationContents, error) {
	src, err := io.ReadAll(r)
	if err != nil {
		return entities.ValidationContents{}, err
	}

	return parseContentHCL(src, filename, variablesEnabled, outputsEnabled)
}

func parseContentHCL(src []byte, filename string, variablesEnabled bool, outputsEnabled bool) (entities.ValidationContents, error) {
	p := hclparse.NewParser()
	validationContents := entities.ValidationContents{}

	f, diags := p.ParseHCL(src, filename)
	if diags.HasErrors() {
		return entities.ValidationContents{}, fmt.Errorf("parsing HCL: %v", diags.Errs())
	}

	// Ignore errors, only focus on content to ignore non outputs/variables
	content, _ := f.Body.Content(validationschema.RootSchema())

	if variablesEnabled {
		variables, err := parseVariables(content.Blocks.OfType("variable"))
		if err != nil {
			return entities.ValidationContents{}, fmt.Errorf("parsing variables: %v", err)
		}
		validationContents.Variables = variables
	}

	if outputsEnabled {
		outputs, err := parseOutputs(content.Blocks.OfType("output"))
		if err != nil {
			return entities.ValidationContents{}, fmt.Errorf("parsing outputs: %v", err)
		}
		validationContents.Outputs = outputs
	}

	return validationContents, nil
}

func parseVariables(variableBlocks hcl.Blocks) (variables []entities.Variable, err error) {
	for _, varBlk := range variableBlocks {
		variable, err := parseVariable(varBlk)
		if err != nil {
			return nil, fmt.Errorf("parsing variable: %s", err)
		}

		variables = append(variables, variable)
	}

	return variables, nil
}

func parseVariable(variableBlock *hcl.Block) (entities.Variable, error) {
	if len(variableBlock.Labels) != 1 {
		return entities.Variable{}, errors.New("variable block must have a single label")
	}

	variableContent, diags := variableBlock.Body.Content(varsschema.VariableSchema())
	if diags.HasErrors() && len(variableContent.Attributes) != 1 {
		for _, e := range diags.Errs() {
			// Only return if parsing error is relevant to `type`
			if strings.Contains(e.Error(), "type") {
				return entities.Variable{}, fmt.Errorf("parsing variable: %v", e)
			}
		}
	}

	// variable blocks are required to have a label as defined in the schema
	name := variableBlock.Labels[0]
	variable, err := createVariableFromHCLAttributes(variableContent.Attributes, name)
	if err != nil {
		return entities.Variable{}, fmt.Errorf("parsing variable: %s", err)
	}

	return variable, nil
}

func createVariableFromHCLAttributes(attrs hcl.Attributes, name string) (entities.Variable, error) {
	var err error

	variable := entities.Variable{Name: name}

	variable.Default, err = hclparser.GetAttribute(attrs, "default").RawJSON()
	if err != nil {
		return entities.Variable{}, err
	}

	// type definition
	variable.Type, err = hclparser.GetAttribute(attrs, "type").VarType()
	if err != nil {
		return entities.Variable{}, err
	}

	return variable, nil
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

	// output blocks are required to have a label as defined in the schema
	name := outputBlock.Labels[0]
	output := entities.Output{Name: name}

	return output, nil
}

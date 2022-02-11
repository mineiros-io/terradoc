package varsparser

import (
	"errors"
	"fmt"
	"io"

	"github.com/hashicorp/hcl/v2"
	"github.com/hashicorp/hcl/v2/hclparse"
	"github.com/mineiros-io/terradoc/internal/entities"
	"github.com/mineiros-io/terradoc/internal/parsers/hclparser"
	"github.com/mineiros-io/terradoc/internal/schemas/varsschema"
)

func Parse(r io.Reader, filename string) (entities.VariablesFile, error) {
	src, err := io.ReadAll(r)
	if err != nil {
		return entities.VariablesFile{}, err
	}

	return parseVariablesHCL(src, filename)
}

func parseVariablesHCL(src []byte, filename string) (entities.VariablesFile, error) {
	p := hclparse.NewParser()

	f, diags := p.ParseHCL(src, filename)
	if diags.HasErrors() {
		return entities.VariablesFile{}, fmt.Errorf("parsing HCL: %v", diags.Errs())
	}

	content, diags := f.Body.Content(varsschema.RootSchema())
	if diags.HasErrors() {
		return entities.VariablesFile{}, fmt.Errorf("getting body content: %v", diags.Errs())
	}

	variables, err := parseVariables(content.Blocks.OfType("variable"))
	if err != nil {
		return entities.VariablesFile{}, fmt.Errorf("parsing variables: %v", err)
	}

	return entities.VariablesFile{Variables: variables}, nil
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
	if diags.HasErrors() {
		return entities.Variable{}, fmt.Errorf("parsing variable: %v", diags.Errs())
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

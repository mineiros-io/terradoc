package hclparser

import (
	"errors"
	"fmt"

	"github.com/hashicorp/hcl/v2"
	"github.com/mineiros-io/terradoc/internal/entities"
	"github.com/mineiros-io/terradoc/internal/parser/hclparser/hclschema"
)

func parseVariables(variableBlocks []*hcl.Block) (variables []entities.Variable, err error) {
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
		return entities.Variable{}, errors.New("variable block does not have a name")
	}

	variableContent, diags := variableBlock.Body.Content(hclschema.VariableSchema())
	if diags.HasErrors() {
		return entities.Variable{}, fmt.Errorf("parsing variable: %v", diags.Errs())
	}

	// variables have only the `name` label
	name := variableBlock.Labels[0]
	variable, err := createVariableFromHCLAttributes(variableContent.Attributes, name)
	if err != nil {
		return entities.Variable{}, fmt.Errorf("parsing variable: %s", err)
	}

	// variables have only `attribute` blocks
	attributes, err := parseVariableAttributes(variableContent.Blocks.OfType(attributeBlockName))
	if err != nil {
		return entities.Variable{}, fmt.Errorf("parsing variable attributes: %s", err)
	}
	variable.Attributes = attributes

	return variable, nil
}

func createVariableFromHCLAttributes(attrs hcl.Attributes, name string) (entities.Variable, error) {
	var err error

	variable := entities.Variable{Name: name}

	variable.Description, err = getAttribute(attrs, descriptionAttributeName).String()
	if err != nil {
		return entities.Variable{}, err
	}

	variable.Default, err = getAttribute(attrs, defaultAttributeName).RawJSON()
	if err != nil {
		return entities.Variable{}, err
	}

	variable.Required, err = getAttribute(attrs, requiredAttributeName).Bool()
	if err != nil {
		return entities.Variable{}, err
	}

	variable.ForcesRecreation, err = getAttribute(attrs, forcesRecreationAttributeName).Bool()
	if err != nil {
		return entities.Variable{}, err
	}

	variable.ReadmeExample, err = getAttribute(attrs, readmeExampleAttributeName).String()
	if err != nil {
		return entities.Variable{}, err
	}

	// type definition
	readmeType := getAttribute(attrs, readmeTypeAttributeName)
	if !readmeType.isNil() {
		variable.Type, err = readmeType.VarTypeFromString()
	} else {
		variable.Type, err = getAttribute(attrs, typeAttributeName).VarType()
	}

	if err != nil {
		return entities.Variable{}, err
	}

	return variable, nil
}

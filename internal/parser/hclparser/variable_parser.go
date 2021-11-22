package hclparser

import (
	"errors"
	"fmt"

	"github.com/hashicorp/hcl/v2"
	"github.com/mineiros-io/terradoc/internal/entities"
	"github.com/mineiros-io/terradoc/internal/parser/hclparser/hclschema"
)

func parseVariable(variableBlock *hcl.Block) (entities.Variable, error) {
	variableContent, diags := variableBlock.Body.Content(hclschema.VariableSchema())
	if diags.HasErrors() {
		return entities.Variable{}, fmt.Errorf("parsing variable: %v", diags.Errs())
	}

	if len(variableBlock.Labels) != 1 {
		return entities.Variable{}, errors.New("variable block does not have a name")
	}

	name := variableBlock.Labels[0]
	variable, err := createVariableFromHCLAttributes(variableContent.Attributes, name)
	if err != nil {
		return entities.Variable{}, fmt.Errorf("parsing variable: %s", err)
	}

	// variables have only `attribute` blocks
	for _, blk := range variableContent.Blocks.OfType(attributeBlockName) {
		attribute, err := parseAttribute(blk, variableAttributeLevel)
		if err != nil {
			return entities.Variable{}, fmt.Errorf("parsing variable attributes: %s", err)
		}

		variable.Attributes = append(variable.Attributes, attribute)
	}

	return variable, nil
}

func createVariableFromHCLAttributes(attrs hcl.Attributes, name string) (entities.Variable, error) {
	variable := entities.Variable{Name: name}

	// description
	description, err := getAttribute(attrs, descriptionAttributeName).String()
	if err != nil {
		return entities.Variable{}, err
	}
	variable.Description = description

	// default
	varDefault, err := getAttribute(attrs, defaultAttributeName).RawJSON()
	if err != nil {
		return entities.Variable{}, err
	}
	variable.Default = varDefault

	// required
	required, err := getAttribute(attrs, requiredAttributeName).Bool()
	if err != nil {
		return entities.Variable{}, err
	}
	variable.Required = required

	// forcesRecreation
	forcesRecreation, err := getAttribute(attrs, forcesRecreationAttributeName).Bool()
	if err != nil {
		return entities.Variable{}, err
	}
	variable.ForcesRecreation = forcesRecreation

	// readme example
	readmeExample, err := getAttribute(attrs, readmeExampleAttributeName).HCLString()
	if err != nil {
		return entities.Variable{}, err
	}
	variable.ReadmeExample = readmeExample

	// type definition
	typeDefinition, err := getType(attrs, name)
	if err != nil {
		return entities.Variable{}, err
	}
	variable.Type = typeDefinition

	return variable, nil
}

package varsvalidator_test

import (
	"testing"

	"github.com/mineiros-io/terradoc/internal/entities"
	"github.com/mineiros-io/terradoc/internal/types"
	"github.com/mineiros-io/terradoc/internal/validators"
	"github.com/mineiros-io/terradoc/internal/validators/varsvalidator"
	"github.com/mineiros-io/terradoc/test"
)

func TestValidate(t *testing.T) {
	tests := []struct {
		desc                   string
		docVariables           entities.VariableCollection
		variablesFileVariables entities.VariableCollection
		wantMissingDoc         []string
		wantMissingDef         []string
		wantTypeMismatch       []validators.TypeMismatchResult
	}{
		{

			desc:                   "when a variable is missing from variables file",
			variablesFileVariables: entities.VariableCollection{},
			docVariables: entities.VariableCollection{
				{
					Name: "name",
					Type: entities.Type{TFType: types.TerraformString},
				},
			},
			wantMissingDef: []string{"name"},
		},
		{
			desc: "when a variable is missing from doc file",
			variablesFileVariables: entities.VariableCollection{
				{
					Name: "age",
					Type: entities.Type{TFType: types.TerraformNumber},
				},
			},
			docVariables:   entities.VariableCollection{},
			wantMissingDoc: []string{"age"},
		},
		{
			desc: "when doc and variables file have the same variables",
			variablesFileVariables: entities.VariableCollection{
				{
					Name: "name",
					Type: entities.Type{TFType: types.TerraformString},
				},
				{
					Name: "age",
					Type: entities.Type{TFType: types.TerraformNumber},
				},
			},
			docVariables: entities.VariableCollection{
				{
					Name: "name",
					Type: entities.Type{TFType: types.TerraformString},
				},
				{
					Name: "age",
					Type: entities.Type{TFType: types.TerraformNumber},
				},
			},
		},
		{
			desc: "when a variable has different types on doc and variables file",
			variablesFileVariables: entities.VariableCollection{
				{
					Name: "age",
					Type: entities.Type{TFType: types.TerraformNumber},
				},
			},
			docVariables: entities.VariableCollection{
				{
					Name: "age",
					Type: entities.Type{TFType: types.TerraformString},
				},
			},
			wantTypeMismatch: []validators.TypeMismatchResult{
				{
					Name:           "age",
					DefinedType:    "number",
					DocumentedType: "string",
				},
			},
		},
		{
			desc: "when a variable is missing from variables file, another missing from doc and another with type mismatch",
			docVariables: entities.VariableCollection{
				{
					Name: "name",
					Type: entities.Type{TFType: types.TerraformString},
				},
				{
					Name: "birth",
					Type: entities.Type{TFType: types.TerraformBool},
				},
			},
			variablesFileVariables: entities.VariableCollection{
				{
					Name: "age",
					Type: entities.Type{TFType: types.TerraformNumber},
				},
				{
					Name: "birth",
					Type: entities.Type{TFType: types.TerraformString},
				},
			},
			wantMissingDef: []string{"name"},
			wantMissingDoc: []string{"age"},
			wantTypeMismatch: []validators.TypeMismatchResult{
				{
					Name:           "birth",
					DefinedType:    "string",
					DocumentedType: "bool",
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.desc, func(t *testing.T) {
			def := definitionFromVariables(tt.docVariables)
			of := variableFileFromVariables(tt.variablesFileVariables)

			got := varsvalidator.Validate(def, of)

			test.AssertHasStrings(t, tt.wantMissingDef, got.MissingDefinition)
			test.AssertHasStrings(t, tt.wantMissingDoc, got.MissingDocumentation)
			test.AssertHasTypeMismatches(t, tt.wantTypeMismatch, got.TypeMismatch)
		})
	}
}

func definitionFromVariables(variables entities.VariableCollection) entities.Doc {
	section := entities.Section{Variables: variables}

	return entities.Doc{Sections: []entities.Section{section}}
}

func variableFileFromVariables(variables entities.VariableCollection) entities.VariablesFile {
	return entities.VariablesFile{Variables: variables}
}

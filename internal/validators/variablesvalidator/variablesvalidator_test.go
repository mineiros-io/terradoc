package variablesvalidator_test

import (
	"testing"

	"github.com/mineiros-io/terradoc/internal/entities"
	"github.com/mineiros-io/terradoc/internal/types"
	"github.com/mineiros-io/terradoc/internal/validators/variablesvalidator"
)

func TestValidate(t *testing.T) {
	tests := []struct {
		desc                   string
		tfdocVariables         entities.VariableCollection
		variablesFileVariables entities.VariableCollection
		wantMissingDoc         []string
		wantMissingDef         []string
		wantTypeMismatch       []string
	}{
		{

			desc:                   "when a variable is missing from variables file",
			variablesFileVariables: entities.VariableCollection{},
			tfdocVariables: entities.VariableCollection{
				{
					Name: "name",
					Type: entities.Type{TFType: types.TerraformString},
				},
			},
			wantMissingDef: []string{"name"},
		},
		{
			desc: "when a variable is missing from tfdoc file",
			variablesFileVariables: entities.VariableCollection{
				{
					Name: "age",
					Type: entities.Type{TFType: types.TerraformNumber},
				},
			},
			tfdocVariables: entities.VariableCollection{},
			wantMissingDoc: []string{"age"},
		},
		{
			desc: "when tfdoc and variables file have the same variables",
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
			tfdocVariables: entities.VariableCollection{
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
			desc: "when a variable has different types on tfdoc and variables file",
			variablesFileVariables: entities.VariableCollection{
				{
					Name: "age",
					Type: entities.Type{TFType: types.TerraformNumber},
				},
			},
			tfdocVariables: entities.VariableCollection{
				{
					Name: "age",
					Type: entities.Type{TFType: types.TerraformString},
				},
			},
			wantTypeMismatch: []string{"age"},
		},
		{
			desc: "when a variable is missing from variables file, another missing from tfdoc and another with type mismatch",
			tfdocVariables: entities.VariableCollection{
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
			wantMissingDef:   []string{"name"},
			wantMissingDoc:   []string{"age"},
			wantTypeMismatch: []string{"birth"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.desc, func(t *testing.T) {
			def := definitionFromVariables(tt.tfdocVariables)
			of := variableFileFromVariables(tt.variablesFileVariables)

			got := variablesvalidator.Validate(def, of)

			assertHasVariables(t, tt.wantMissingDef, got.MissingDefinition)
			assertHasVariables(t, tt.wantMissingDoc, got.MissingDocumentation)
			assertHasVariables(t, tt.wantTypeMismatch, got.TypeMismatch)
		})
	}
}

func assertHasVariables(t *testing.T, want []string, got entities.VariableCollection) {
	t.Helper()

	for _, name := range want {
		found := false

		for _, o := range got {
			if name == o.Name {
				found = true
			}
		}

		if !found {
			t.Errorf("wanted %q to be found in %v", name, got)
		}
	}
}

func definitionFromVariables(variables entities.VariableCollection) entities.Definition {
	section := entities.Section{Variables: variables}

	return entities.Definition{Sections: []entities.Section{section}}
}

func variableFileFromVariables(variables entities.VariableCollection) entities.VariablesFile {
	return entities.VariablesFile{Variables: variables}
}

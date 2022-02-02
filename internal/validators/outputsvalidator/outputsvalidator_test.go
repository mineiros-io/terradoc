package outputsvalidator_test

import (
	"testing"

	"github.com/mineiros-io/terradoc/internal/entities"
	"github.com/mineiros-io/terradoc/internal/types"
	"github.com/mineiros-io/terradoc/internal/validators/outputsvalidator"
)

func TestValidate(t *testing.T) {
	tests := []struct {
		desc               string
		tfdocOutputs       entities.OutputCollection
		outputsFileOutputs entities.OutputCollection
		wantMissingDoc     []string
		wantMissingDef     []string
		wantTypeMismatch   []string
	}{
		{

			desc:               "when an output is missing from outputs file",
			outputsFileOutputs: entities.OutputCollection{},
			tfdocOutputs: entities.OutputCollection{
				{
					Name: "name",
					Type: entities.Type{TFType: types.TerraformString},
				},
			},
			wantMissingDef: []string{"name"},
		},
		{
			desc: "when an output is missing from tfdoc file",
			outputsFileOutputs: entities.OutputCollection{
				{
					Name: "age",
					Type: entities.Type{TFType: types.TerraformNumber},
				},
			},
			tfdocOutputs:   entities.OutputCollection{},
			wantMissingDoc: []string{"age"},
		},
		{
			desc: "when tfdoc and outputs file have the same outputs",
			outputsFileOutputs: entities.OutputCollection{
				{
					Name: "name",
					Type: entities.Type{TFType: types.TerraformString},
				},
				{
					Name: "age",
					Type: entities.Type{TFType: types.TerraformNumber},
				},
			},
			tfdocOutputs: entities.OutputCollection{
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
			desc: "when an output has different types on tfdoc and outputs file",
			outputsFileOutputs: entities.OutputCollection{
				{
					Name: "age",
					Type: entities.Type{TFType: types.TerraformNumber},
				},
			},
			tfdocOutputs: entities.OutputCollection{
				{
					Name: "age",
					Type: entities.Type{TFType: types.TerraformString},
				},
			},
			wantTypeMismatch: []string{"age"},
		},
		{
			desc: "when an output is missing from outputs file, another missing from tfdoc and another with type mismatch",
			tfdocOutputs: entities.OutputCollection{
				{
					Name: "name",
					Type: entities.Type{TFType: types.TerraformString},
				},
				{
					Name: "birth",
					Type: entities.Type{TFType: types.TerraformBool},
				},
			},
			outputsFileOutputs: entities.OutputCollection{
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
			def := definitionFromOutputs(tt.tfdocOutputs)
			of := outputFileFromOutputs(tt.outputsFileOutputs)

			got := outputsvalidator.Validate(def, of)

			assertHasOutputs(t, tt.wantMissingDef, got.MissingDefinition)
			assertHasOutputs(t, tt.wantMissingDoc, got.MissingDocumentation)
			assertHasOutputs(t, tt.wantTypeMismatch, got.TypeMismatch)
		})
	}
}

func assertHasOutputs(t *testing.T, want []string, got entities.OutputCollection) {
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

func definitionFromOutputs(outputs entities.OutputCollection) entities.Definition {
	section := entities.Section{Outputs: outputs}

	return entities.Definition{Sections: []entities.Section{section}}
}

func outputFileFromOutputs(outputs entities.OutputCollection) entities.OutputsFile {
	return entities.OutputsFile{Outputs: outputs}
}

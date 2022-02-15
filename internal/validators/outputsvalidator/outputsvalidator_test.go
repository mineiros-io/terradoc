package outputsvalidator_test

import (
	"testing"

	"github.com/mineiros-io/terradoc/internal/entities"
	"github.com/mineiros-io/terradoc/internal/types"
	"github.com/mineiros-io/terradoc/internal/validators/outputsvalidator"
	"github.com/mineiros-io/terradoc/test"
)

func TestValidate(t *testing.T) {
	tests := []struct {
		desc               string
		docOutputs         entities.OutputCollection
		outputsFileOutputs entities.OutputCollection
		wantMissingDoc     []string
		wantMissingDef     []string
	}{
		{

			desc:               "when an output is missing from outputs file",
			outputsFileOutputs: entities.OutputCollection{},
			docOutputs: entities.OutputCollection{
				{
					Name: "name",
					Type: entities.Type{TFType: types.TerraformString},
				},
			},
			wantMissingDef: []string{"name"},
		},
		{
			desc: "when an output is missing from doc file",
			outputsFileOutputs: entities.OutputCollection{
				{
					Name: "age",
					Type: entities.Type{TFType: types.TerraformNumber},
				},
			},
			docOutputs:     entities.OutputCollection{},
			wantMissingDoc: []string{"age"},
		},
		{
			desc: "when doc and outputs file have the same outputs",
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
			docOutputs: entities.OutputCollection{
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
	}

	for _, tt := range tests {
		t.Run(tt.desc, func(t *testing.T) {
			def := definitionFromOutputs(tt.docOutputs)
			of := outputFileFromOutputs(tt.outputsFileOutputs)

			got := outputsvalidator.Validate(def, of)

			test.AssertHasStrings(t, tt.wantMissingDef, got.MissingDefinition)
			test.AssertHasStrings(t, tt.wantMissingDoc, got.MissingDocumentation)
		})
	}
}

func definitionFromOutputs(outputs entities.OutputCollection) entities.Doc {
	section := entities.Section{Outputs: outputs}

	return entities.Doc{Sections: []entities.Section{section}}
}

func outputFileFromOutputs(outputs entities.OutputCollection) entities.OutputsFile {
	return entities.OutputsFile{Outputs: outputs}
}

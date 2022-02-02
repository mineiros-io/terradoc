package outputsvalidator

import (
	"github.com/mineiros-io/terradoc/internal/entities"
	"github.com/mineiros-io/terradoc/internal/validators"
)

type OutputsValidationSummary struct {
	MissingDefinition    entities.OutputCollection
	MissingDocumentation entities.OutputCollection
	TypeMismatch         entities.OutputCollection
}

func (vs OutputsValidationSummary) Valid() bool {
	return len(vs.MissingDocumentation) == 0 &&
		len(vs.MissingDefinition) == 0 &&
		len(vs.TypeMismatch) == 0
}

// TODO: receive resource or collection of outputs directly?
func Validate(tfdoc entities.Definition, outputsFile entities.OutputsFile) OutputsValidationSummary {
	summary := OutputsValidationSummary{}

	outputsFileOutputs := outputsFile.Outputs
	// TODO
	tfdocOutputs := entities.OutputCollection{}
	for _, s := range tfdoc.Sections {
		tfdocOutputs = append(tfdocOutputs, s.AllOutputs()...)
	}

	for _, defOutput := range tfdocOutputs {
		fOutput, exists := outputsFileOutputs.OutputByName(defOutput.Name)
		if !exists {
			summary.MissingDefinition = append(summary.MissingDefinition, defOutput)

			continue
		}

		if !validators.TypesMatch(&defOutput.Type, &fOutput.Type) {
			summary.TypeMismatch = append(summary.TypeMismatch, defOutput)

			continue
		}
	}

	for _, fOutput := range outputsFileOutputs {
		_, exists := tfdocOutputs.OutputByName(fOutput.Name)
		if !exists {
			summary.MissingDocumentation = append(summary.MissingDocumentation, fOutput)
		}
	}

	return summary
}

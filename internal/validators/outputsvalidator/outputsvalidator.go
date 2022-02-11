package outputsvalidator

import (
	"github.com/mineiros-io/terradoc/internal/entities"
	"github.com/mineiros-io/terradoc/internal/validators"
)

const CheckType = "output"

type outputValidationChecks map[string]outputValidation

type outputValidation struct {
	defined    entities.Output
	documented entities.Output
}

func Validate(doc entities.Doc, outputsFile entities.OutputsFile) validators.Summary {
	summary := validators.Summary{Type: CheckType}

	validationResult := validateOutputs(doc.AllOutputs(), outputsFile.Outputs)

	for outputName, check := range validationResult {
		switch {
		// TODO: using missing name as indicator of missing output
		case check.defined.Name == "":
			summary.MissingDefinition = append(summary.MissingDefinition, outputName)
		case check.documented.Name == "":
			summary.MissingDocumentation = append(summary.MissingDocumentation, outputName)
		case !validators.TypesMatch(&check.defined.Type, &check.documented.Type):
			summary.TypeMismatch = append(
				summary.TypeMismatch,
				validators.TypeMismatchResult{
					Name:           outputName,
					DefinedType:    check.defined.Type.AsString(),
					DocumentedType: check.documented.Type.AsString(),
				},
			)
		}
	}

	return summary
}

func validateOutputs(docOutputs, outputsFileOutputs []entities.Output) outputValidationChecks {
	result := outputValidationChecks{}

	for _, tfOutput := range docOutputs {
		result[tfOutput.Name] = outputValidation{documented: tfOutput}
	}

	for _, fOutput := range outputsFileOutputs {
		val, ok := result[fOutput.Name]
		if !ok {
			val = outputValidation{}
		}
		val.defined = fOutput

		result[fOutput.Name] = val
	}

	return result
}

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

func Validate(doc entities.Doc, outputsFile entities.ValidationContents) validators.Summary {
	summary := validators.Summary{Type: CheckType}

	validationResult := validateOutputs(doc.AllOutputs(), outputsFile.Outputs)

	for outputName, check := range validationResult {
		switch {
		case check.defined.Name == "":
			summary.MissingDefinition = append(summary.MissingDefinition, outputName)
		case check.documented.Name == "":
			summary.MissingDocumentation = append(summary.MissingDocumentation, outputName)
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

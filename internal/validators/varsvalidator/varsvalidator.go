package varsvalidator

import (
	"github.com/mineiros-io/terradoc/internal/entities"
	"github.com/mineiros-io/terradoc/internal/validators"
)

const CheckType = "variable"

type variableValidationChecks map[string]variableValidation

type variableValidation struct {
	defined    entities.Variable
	documented entities.Variable
}

func Validate(doc entities.Doc, varsFile entities.ValidationContents) validators.Summary {
	summary := validators.Summary{Type: CheckType}

	validationResult := validateVariables(doc.AllVariables(), varsFile.Variables)

	for varName, check := range validationResult {
		switch {
		case check.defined.Name == "":
			summary.MissingDefinition = append(summary.MissingDefinition, varName)
		case check.documented.Name == "":
			summary.MissingDocumentation = append(summary.MissingDocumentation, varName)
		case !validators.TypesMatch(&check.defined.Type, &check.documented.Type):
			summary.TypeMismatch = append(
				summary.TypeMismatch,
				validators.TypeMismatchResult{
					Name:           varName,
					DefinedType:    check.defined.Type.AsString(),
					DocumentedType: check.documented.Type.AsString(),
				},
			)
		}
	}

	return summary
}

func validateVariables(docVars, varFileVars []entities.Variable) variableValidationChecks {
	result := variableValidationChecks{}

	for _, tfVar := range docVars {
		result[tfVar.Name] = variableValidation{documented: tfVar}
	}

	for _, fVar := range varFileVars {
		val, ok := result[fVar.Name]
		if !ok {
			val = variableValidation{}
		}
		val.defined = fVar

		result[fVar.Name] = val
	}

	return result
}

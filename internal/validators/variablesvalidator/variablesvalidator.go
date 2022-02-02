package variablesvalidator

import (
	"github.com/mineiros-io/terradoc/internal/entities"
	"github.com/mineiros-io/terradoc/internal/validators"
)

type VariablesValidationSummary struct {
	MissingDefinition    entities.VariableCollection
	MissingDocumentation entities.VariableCollection
	TypeMismatch         entities.VariableCollection
}

func (vs VariablesValidationSummary) Valid() bool {
	return len(vs.MissingDocumentation) == 0 &&
		len(vs.MissingDefinition) == 0 &&
		len(vs.TypeMismatch) == 0
}

func Validate(tfdoc entities.Definition, varFile entities.VariablesFile) VariablesValidationSummary {
	summary := VariablesValidationSummary{}

	variablesFileVars := varFile.Variables
	//TODO
	tfdocVars := entities.VariableCollection{}
	for _, s := range tfdoc.Sections {
		tfdocVars = append(tfdocVars, s.AllVariables()...)
	}

	for _, defVar := range tfdocVars {
		fVar, exists := variablesFileVars.VarByName(defVar.Name)
		if !exists {
			summary.MissingDefinition = append(summary.MissingDefinition, defVar)
			continue
		}

		if !validators.TypesMatch(&defVar.Type, &fVar.Type) {
			summary.TypeMismatch = append(summary.TypeMismatch, defVar)
			continue
		}
	}

	for _, fVar := range variablesFileVars {
		_, exists := tfdocVars.VarByName(fVar.Name)
		if !exists {
			summary.MissingDocumentation = append(summary.MissingDocumentation, fVar)
		}
	}

	return summary
}

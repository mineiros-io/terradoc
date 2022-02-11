package cli

import (
	"errors"
	"fmt"
	"os"

	"github.com/mineiros-io/terradoc/internal/parsers/outputsparser"
	"github.com/mineiros-io/terradoc/internal/parsers/tfdocparser"
	"github.com/mineiros-io/terradoc/internal/parsers/varsparser"
	"github.com/mineiros-io/terradoc/internal/validators"
	"github.com/mineiros-io/terradoc/internal/validators/outputsvalidator"
	"github.com/mineiros-io/terradoc/internal/validators/varsvalidator"
)

type ValidateCmd struct {
	DocFile       string `arg:"" help:"Input file." type:"existingfile"`
	VariablesFile string `name:"variables" optional:"" short:"v" help:"Variables file" type:"existingfile"`
	OutputsFile   string `name:"outputs" short:"o" optional:"" help:"Outputs file" type:"existingfile"`
}

func (vcm ValidateCmd) Run() error {
	var hasVarsErrors, hasOutputsErrors bool
	// DOC
	t, tCloser, err := openInput(vcm.DocFile)
	if err != nil {
		return err
	}
	defer tCloser()

	doc, err := tfdocparser.Parse(t, t.Name())
	if err != nil {
		return err
	}

	// VARIABLES
	if vcm.VariablesFile != "" {
		v, vCloser, err := openInput(vcm.VariablesFile)
		if err != nil {
			return err
		}
		defer vCloser()

		tfvars, err := varsparser.Parse(v, v.Name())
		if err != nil {
			return err
		}

		varsSummary := varsvalidator.Validate(doc, tfvars)

		printValidationSummary(varsSummary, t.Name(), v.Name())

		hasVarsErrors = !varsSummary.Success()
	}

	// OUTPUTS
	if vcm.OutputsFile != "" {
		o, oCloser, err := openInput(vcm.OutputsFile)
		if err != nil {
			return err
		}
		defer oCloser()

		tfoutputs, err := outputsparser.Parse(o, o.Name())
		if err != nil {
			return err
		}

		outputsSummary := outputsvalidator.Validate(doc, tfoutputs)

		printValidationSummary(outputsSummary, t.Name(), o.Name())

		hasOutputsErrors = !outputsSummary.Success()
	}

	if hasVarsErrors || hasOutputsErrors {
		return errors.New("Found validation errors")
	}

	return nil
}

func printValidationSummary(summary validators.Summary, docFilename, defFilename string) {
	for _, missingDef := range summary.MissingDefinition {
		fmt.Fprintf(os.Stderr, "Missing %s definition: %q is not defined in %q\n", summary.Type, missingDef, defFilename)
	}

	for _, missingDoc := range summary.MissingDocumentation {
		fmt.Fprintf(os.Stderr, "Missing %s documentation: %q is not documented in %q\n", summary.Type, missingDoc, docFilename)
	}

	for _, tMismatch := range summary.TypeMismatch {
		fmt.Fprintf(os.Stderr, "Type mismatch for %s: %q is documented as %q in %q but defined as %q in %q\n", summary.Type, tMismatch.Name, tMismatch.DocumentedType, docFilename, tMismatch.DefinedType, defFilename)
	}

}

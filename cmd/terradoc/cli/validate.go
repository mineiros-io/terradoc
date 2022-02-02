package cli

import (
	"errors"
	"fmt"
	"os"

	"github.com/mineiros-io/terradoc/internal/parsers/outputsparser"
	"github.com/mineiros-io/terradoc/internal/parsers/tfdocparser"
	"github.com/mineiros-io/terradoc/internal/parsers/variablesparser"
	"github.com/mineiros-io/terradoc/internal/validators/outputsvalidator"
	"github.com/mineiros-io/terradoc/internal/validators/variablesvalidator"
)

type ValidateCmd struct {
	TFDocFile     string `arg:"" help:"Input file." type:"existingfile"`
	VariablesFile string `name:"variables" optional:"" short:"v" help:"Variables file" type:"existingfile"`
	OutputsFile   string `name:"outputs" short:"o" optional:"" help:"Outputs file" type:"existingfile"`
}

func (vcm ValidateCmd) Run() error {
	var hasVarsErrors, hasOutputsErrors bool
	// TFDOC
	t, tCloser, err := openInput(vcm.TFDocFile)
	if err != nil {
		return err
	}
	defer tCloser()

	tfdoc, err := tfdocparser.Parse(t, t.Name())
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

		tfvars, err := variablesparser.Parse(v, v.Name())
		if err != nil {
			return err
		}

		varsSummary := variablesvalidator.Validate(tfdoc, tfvars)

		printVariablesValidationSummary(varsSummary, v.Name(), t.Name())

		// TODO
		hasVarsErrors = !varsSummary.Valid()
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

		outputsSummary := outputsvalidator.Validate(tfdoc, tfoutputs)

		printOutputsValidationSummary(outputsSummary, o.Name(), t.Name())

		// TODO
		hasOutputsErrors = !outputsSummary.Valid()
	}

	if hasVarsErrors || hasOutputsErrors {
		return errors.New("Found validation errors")
	}

	return nil
}

func printVariablesValidationSummary(summary variablesvalidator.VariablesValidationSummary, varsFilename, tfdocFilename string) {
	for _, mVar := range summary.MissingDefinition {
		fmt.Fprintf(os.Stderr, "Missing variable definition: %q is not defined in %q\n", mVar.Name, varsFilename)
	}

	for _, mVar := range summary.MissingDocumentation {
		fmt.Fprintf(os.Stderr, "Missing variable documentation: %q is not documented in %q\n", mVar.Name, tfdocFilename)
	}

	for _, mVar := range summary.TypeMismatch {
		fmt.Fprintf(os.Stderr, "Variable type mismatch: %q\n", mVar.Name)
	}
}

func printOutputsValidationSummary(summary outputsvalidator.OutputsValidationSummary, outputsFilename, tfdocFilename string) {
	for _, mOutput := range summary.MissingDefinition {
		fmt.Fprintf(os.Stderr, "Missing output definition: %q is not defined in %q\n", mOutput.Name, outputsFilename)
	}

	for _, mOutput := range summary.MissingDocumentation {
		fmt.Fprintf(os.Stderr, "Missing output documentation: %q is not documented in %q\n", mOutput.Name, tfdocFilename)
	}

	for _, mOutput := range summary.TypeMismatch {
		fmt.Fprintf(os.Stderr, "Output type mismatch: %q\n", mOutput.Name)
	}
}

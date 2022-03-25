package cli

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"

	"github.com/mineiros-io/terradoc/internal/entities"
	"github.com/mineiros-io/terradoc/internal/parsers/docparser"
	"github.com/mineiros-io/terradoc/internal/parsers/validationparser"
	"github.com/mineiros-io/terradoc/internal/validators"
	"github.com/mineiros-io/terradoc/internal/validators/outputsvalidator"
	"github.com/mineiros-io/terradoc/internal/validators/varsvalidator"
)

type ValidateCmd struct {
	DocFile          string `arg:"" help:"Input file." type:"existingfile"`
	VariablesEnabled bool   `name:"variables" optional:"" short:"v" help:"Whether to validate variables."`
	OutputsEnabled   bool   `name:"outputs" short:"o" optional:"" help:"Whether to validate outputs."`
}

func (vcm ValidateCmd) Run() error {
	var hasVarsErrors, hasOutputsErrors bool
	// DOC
	t, tCloser, err := openInput(vcm.DocFile)
	if err != nil {
		return err
	}
	defer tCloser()

	doc, err := docparser.Parse(t, t.Name())
	if err != nil {
		return err
	}

	path, err := os.Getwd()
	if err != nil {
		return err
	}

	files, err := WalkMatch(path, "*.tf")
	if err != nil {
		return err
	}

	varsEnabled := false
	if vcm.VariablesEnabled {
		varsEnabled = true
	} else {
		varsEnabled = !vcm.VariablesEnabled && !vcm.OutputsEnabled
	}

	outputsEnabled := false
	if vcm.OutputsEnabled {
		outputsEnabled = true
	} else {
		outputsEnabled = !vcm.VariablesEnabled && !vcm.OutputsEnabled
	}

	hasVarsErrors = false
	hasOutputsErrors = false

	tfContent := entities.ValidationContents{}

	for _, file := range files {
		f, vCloser, err := openInput(file)
		if err != nil {
			return err
		}
		defer vCloser()

		content, err := validationparser.Parse(f, f.Name(), varsEnabled, outputsEnabled)
		if err != nil {
			return err
		}

		tfContent.Variables = append(tfContent.Variables, content.Variables...)
		tfContent.Outputs = append(tfContent.Outputs, content.Outputs...)
	}

	// VARIABLES
	if varsEnabled {
		varsSummary := varsvalidator.Validate(doc, tfContent)

		printValidationSummary(varsSummary, t.Name())

		if !varsSummary.Success() {
			hasVarsErrors = !varsSummary.Success()
		}
	}

	// OUTPUTS
	if outputsEnabled {
		outputsSummary := outputsvalidator.Validate(doc, tfContent)

		printValidationSummary(outputsSummary, t.Name())

		if !outputsSummary.Success() {
			hasOutputsErrors = !outputsSummary.Success()
		}
	}

	if hasVarsErrors || hasOutputsErrors {
		return errors.New("Found validation errors")
	}

	return nil
}

func printValidationSummary(summary validators.Summary, docFilename string) {
	for _, missingDef := range summary.MissingDefinition {
		fmt.Fprintf(os.Stderr, "Missing %s definition: %q is not defined in .tf files\n", summary.Type, missingDef)
	}

	for _, missingDoc := range summary.MissingDocumentation {
		fmt.Fprintf(os.Stderr, "Missing %s documentation: %q is not documented in %q\n", summary.Type, missingDoc, docFilename)
	}

	for _, tMismatch := range summary.TypeMismatch {
		fmt.Fprintf(os.Stderr, "Type mismatch for %s: %q is documented as %q in %q but defined as %q in .tf files\n", summary.Type, tMismatch.Name, tMismatch.DocumentedType, docFilename, tMismatch.DefinedType)
	}

}

func WalkMatch(root, pattern string) ([]string, error) {
	var matches []string
	err := filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if info.IsDir() && path != root {
			return filepath.SkipDir
		}
		if matched, err := filepath.Match(pattern, filepath.Base(path)); err != nil {
			return err
		} else if matched {
			matches = append(matches, path)
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	return matches, nil
}

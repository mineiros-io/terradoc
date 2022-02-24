package cli

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"

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

	allFiles, err := WalkMatch(path, "*.tf")
	if err != nil {
		return err
	}

	// Ignore any folder or file matching ignore array - cold replace with file in future
	var ignore = []string{"/.terraform/", "/example/", "/.vscode/"}

	var files []string
	for _, file := range allFiles {
		contains := false

		for _, ignoreString := range ignore {
			if !strings.Contains(file, ignoreString) {
				contains = true
				break
			}
		}

		if !contains {
			files = append(files, file)
			fmt.Fprintf(os.Stderr, file, "\n")
		}
	}

	hasVarsErrors = false
	hasOutputsErrors = false

	for _, file := range files {
		f, vCloser, err := openInput(file)
		if err != nil {
			return err
		}
		defer vCloser()

		tfContent, err := validationparser.Parse(f, f.Name(), vcm.VariablesEnabled, vcm.OutputsEnabled)
		if err != nil {
			return err
		}

		// VARIABLES
		varsSummary := varsvalidator.Validate(doc, tfContent)

		printValidationSummary(varsSummary, t.Name(), f.Name())

		if !varsSummary.Success() {
			hasVarsErrors = !varsSummary.Success()
		}

		// OUTPUTS
		outputsSummary := outputsvalidator.Validate(doc, tfContent)

		printValidationSummary(outputsSummary, t.Name(), f.Name())

		if !outputsSummary.Success() {
			hasOutputsErrors = !outputsSummary.Success()
		}
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

func WalkMatch(root, pattern string) ([]string, error) {
	var matches []string
	err := filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() {
			return nil
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

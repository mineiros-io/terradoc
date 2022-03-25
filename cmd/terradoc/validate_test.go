package main_test

import (
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"strings"
	"testing"

	"github.com/madlambda/spells/assert"
	"github.com/mineiros-io/terradoc/internal/validators"
	"github.com/mineiros-io/terradoc/test"
)

func TestValidateVariables(t *testing.T) {
	tests := []struct {
		desc                     string
		doc                      string
		variables                string
		wantMissingDocumentation []string
		wantMissingDefinition    []string
		wantTypeMismatch         []validators.TypeMismatchResult
		wantError                bool
	}{
		{
			desc:      "when `variables.tf` and terradoc file have the same variables",
			doc:       "validate/variables/complete.tfdoc.hcl",
			variables: "validate/variables/complete-variables.tf",
		},
		{
			desc:                     "when `variables.tf` has missing variables",
			doc:                      "validate/variables/complete.tfdoc.hcl",
			variables:                "validate/variables/missing-variables.tf",
			wantMissingDocumentation: []string{},
			wantMissingDefinition:    []string{"beer"},
			wantError:                true,
		},
		{
			desc:                     "when terradoc file has missing variables",
			doc:                      "validate/variables/missing-variables.tfdoc.hcl",
			variables:                "validate/variables/complete-variables.tf",
			wantMissingDocumentation: []string{"beer"},
			wantMissingDefinition:    []string{},
			wantError:                true,
		},
		{
			desc:                     "when `variables.tf` has type mismatch",
			doc:                      "validate/variables/complete.tfdoc.hcl",
			variables:                "validate/variables/type-mismatch.tf",
			wantMissingDocumentation: []string{},
			wantMissingDefinition:    []string{},
			wantTypeMismatch: []validators.TypeMismatchResult{
				{
					Name:           "number",
					DefinedType:    "list(string)",
					DocumentedType: "number",
				},
			},
			wantError: true,
		},
		{
			desc:                     "when `variables.tf` has type mismatch and missing variables",
			doc:                      "validate/variables/complete.tfdoc.hcl",
			variables:                "validate/variables/type-mismatch-with-missing.tf",
			wantMissingDocumentation: []string{},
			wantMissingDefinition:    []string{"beer"},
			wantTypeMismatch: []validators.TypeMismatchResult{
				{
					Name:           "person",
					DefinedType:    "string",
					DocumentedType: "object(person)",
				},
			},
			wantError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.desc, func(t *testing.T) {
			doc := test.ReadFixture(t, tt.doc)
			docFile, err := ioutil.TempFile(t.TempDir(), "terradoc-validate-doc-")
			assert.NoError(t, err)
			_, err = docFile.Write(doc)
			assert.NoError(t, err)

			variables := test.ReadFixture(t, tt.variables)

			// Break variables up into multiple files
			variablesDir := t.TempDir()

			for _, v := range strings.Split(string(variables[:]), "\n\n") {
				variablesFile, err := ioutil.TempFile(variablesDir, "terradoc-validate-variables-*.tf")
				assert.NoError(t, err)
				_, err = variablesFile.Write([]byte(v))
				assert.NoError(t, err)
			}

			err = os.Chdir(variablesDir)
			assert.NoError(t, err)

			cmd := exec.Command(terradocBinPath, "validate", docFile.Name(), "-v")

			output, err := cmd.CombinedOutput()

			gotResult := splitOutputMessages(t, output, "variable")

			if tt.wantError {
				assert.Error(t, err)

				assertHasMissingDocumentation(t, docFile.Name(), gotResult.missingDocumentation, tt.wantMissingDocumentation, "variable")
				assertHasMissingDefinition(t, gotResult.missingDefinition, tt.wantMissingDefinition, "variable")
				assertHasTypeMismatch(t, docFile.Name(), tt.wantTypeMismatch, gotResult.typeMismatch, "variable")
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestValidateOutputs(t *testing.T) {
	tests := []struct {
		desc                     string
		doc                      string
		outputs                  string
		wantMissingDocumentation []string
		wantMissingDefinition    []string
		wantTypeMismatch         []validators.TypeMismatchResult
		wantError                bool
	}{
		{
			desc:    "when `outputs.tf` and terradoc file have the same outputs",
			doc:     "validate/outputs/complete.tfdoc.hcl",
			outputs: "validate/outputs/complete-outputs.tf",
		},
		{
			desc:                     "when `outputs.tf` has missing outputs",
			doc:                      "validate/outputs/complete.tfdoc.hcl",
			outputs:                  "validate/outputs/missing-outputs.tf",
			wantMissingDocumentation: []string{},
			wantMissingDefinition:    []string{"beer"},
			wantError:                true,
		},
		{
			desc:                     "when terradoc file has missing outputs",
			doc:                      "validate/outputs/missing-outputs.tfdoc.hcl",
			outputs:                  "validate/outputs/complete-outputs.tf",
			wantMissingDocumentation: []string{"beer"},
			wantMissingDefinition:    []string{},
			wantError:                true,
		},
		{
			desc:                     "when `outputs.tf` has type mismatch",
			doc:                      "validate/outputs/complete.tfdoc.hcl",
			outputs:                  "validate/outputs/type-mismatch.tf",
			wantMissingDocumentation: []string{},
			wantMissingDefinition:    []string{},
			wantTypeMismatch: []validators.TypeMismatchResult{
				{
					Name:           "number",
					DefinedType:    "list(string)",
					DocumentedType: "number",
				},
			},
			wantError: true,
		},
		{
			desc:                     "when `outputs.tf` has type mismatch and missing outputs",
			doc:                      "validate/outputs/complete.tfdoc.hcl",
			outputs:                  "validate/outputs/type-mismatch-with-missing.tf",
			wantMissingDocumentation: []string{},
			wantMissingDefinition:    []string{"beer"},
			wantTypeMismatch: []validators.TypeMismatchResult{
				{
					Name:           "person",
					DefinedType:    "string",
					DocumentedType: "object(person)",
				},
			},
			wantError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.desc, func(t *testing.T) {
			doc := test.ReadFixture(t, tt.doc)
			docFile, err := ioutil.TempFile(t.TempDir(), "terradoc-validate-doc-")
			assert.NoError(t, err)
			_, err = docFile.Write(doc)
			assert.NoError(t, err)

			outputs := test.ReadFixture(t, tt.outputs)

			// Break outputs up into multiple files
			outputsDir := t.TempDir()

			for _, v := range strings.Split(string(outputs[:]), "\n\n") {
				outputsFile, err := ioutil.TempFile(outputsDir, "terradoc-validate-outputs-*.tf")
				assert.NoError(t, err)
				_, err = outputsFile.Write([]byte(v))
				assert.NoError(t, err)
			}

			err = os.Chdir(outputsDir)
			assert.NoError(t, err)

			cmd := exec.Command(terradocBinPath, "validate", docFile.Name(), "-o")

			output, err := cmd.CombinedOutput()

			gotResult := splitOutputMessages(t, output, "output")

			if tt.wantError {
				assert.Error(t, err)

				assertHasMissingDocumentation(t, docFile.Name(), gotResult.missingDocumentation, tt.wantMissingDocumentation, "output")
				assertHasMissingDefinition(t, gotResult.missingDefinition, tt.wantMissingDefinition, "output")
				assertHasTypeMismatch(t, docFile.Name(), tt.wantTypeMismatch, gotResult.typeMismatch, "output")
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

type validationResult struct {
	missingDocumentation []string
	missingDefinition    []string
	typeMismatch         []string
}

func splitOutputMessages(t *testing.T, output []byte, validationType string) validationResult {
	result := validationResult{}
	outputStrings := strings.Split(string(output), "\n")

	for _, oo := range outputStrings {
		switch {
		case strings.HasPrefix(oo, fmt.Sprintf("Missing %s definition:", validationType)):
			result.missingDefinition = append(result.missingDefinition, oo)
		case strings.HasPrefix(oo, fmt.Sprintf("Missing %s documentation:", validationType)):
			result.missingDocumentation = append(result.missingDocumentation, oo)
		case strings.HasPrefix(oo, fmt.Sprintf("Type mismatch for %s:", validationType)):
			result.typeMismatch = append(result.typeMismatch, oo)
		}
	}

	return result
}

func assertHasMissingDocumentation(t *testing.T, filename string, got, want []string, validationType string) {
	t.Helper()

	if len(got) != len(want) {
		t.Errorf("wanted %v but got %v", want, got)
	}

	for _, wantStr := range want {
		found := false

		completeWantString := fmt.Sprintf("Missing %s documentation: %q is not documented in %q", validationType, wantStr, filename)

		for _, msg := range got {
			if msg == completeWantString {
				found = true

				break
			}
		}

		if !found {
			t.Errorf("wanted output to have missing documentation for %q but didn't find it", wantStr)
		}
	}
}

func assertHasMissingDefinition(t *testing.T, got, want []string, validationType string) {
	t.Helper()

	if len(got) != len(want) {
		t.Errorf("wanted %v but got %v", want, got)
	}

	for _, wantStr := range want {
		found := false

		completeWantString := fmt.Sprintf("Missing %s definition: %q is not defined in .tf files", validationType, wantStr)
		for _, msg := range got {
			if msg == completeWantString {
				found = true

				break
			}
		}

		if !found {
			t.Errorf("wanted output to have missing definition %q but didn't find it", wantStr)
		}
	}
}

func assertHasTypeMismatch(t *testing.T, docFilename string, want []validators.TypeMismatchResult, got []string, validationType string) {
	t.Helper()

	if len(got) != len(want) {
		t.Errorf("wanted %v but got %v", want, got)
	}

	for _, tm := range want {
		found := false

		completeWantString := fmt.Sprintf("Type mismatch for %s: %q is documented as %q in %q but defined as %q in .tf files", validationType, tm.Name, tm.DocumentedType, docFilename, tm.DefinedType)

		for _, msg := range got {
			if msg == completeWantString {
				found = true

				break
			}
		}

		if !found {
			t.Errorf("wanted output to have type mismatch for %q but didn't find it", tm.Name)
		}
	}
}

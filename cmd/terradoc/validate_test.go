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
			doc:       "validate/complete.tfdoc.hcl",
			variables: "validate/complete-variables.tf",
		},
		{
			desc:                     "when `variables.tf` has missing variables",
			doc:                      "validate/complete.tfdoc.hcl",
			variables:                "validate/missing-variables.tf",
			wantMissingDocumentation: []string{},
			wantMissingDefinition:    []string{"beer"},
			wantError:                true,
		},
		{
			desc:                     "when terradoc file has missing variables",
			doc:                      "validate/missing-variables.tfdoc.hcl",
			variables:                "validate/complete-variables.tf",
			wantMissingDocumentation: []string{"beer"},
			wantMissingDefinition:    []string{},
			wantError:                true,
		},
		{
			desc:                     "when `variables.tf` has type mismatch",
			doc:                      "validate/complete.tfdoc.hcl",
			variables:                "validate/type-mismatch.tf",
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
			doc:                      "validate/complete.tfdoc.hcl",
			variables:                "validate/type-mismatch-with-missing.tf",
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
			variablesFile, err := ioutil.TempFile(t.TempDir(), "terradoc-validate-variables-")
			assert.NoError(t, err)
			_, err = variablesFile.Write(variables)
			assert.NoError(t, err)

			os.Chdir(os.TempDir())

			cmd := exec.Command(terradocBinPath, "validate", docFile.Name(), "-v")

			output, err := cmd.CombinedOutput()

			gotResult := splitOutputMessages(t, output)

			if tt.wantError {
				assert.Error(t, err)

				assertHasMissingDocumentation(t, docFile.Name(), gotResult.missingDocumentation, tt.wantMissingDocumentation)
				assertHasMissingDefinition(t, variablesFile.Name(), gotResult.missingDefinition, tt.wantMissingDefinition)
				assertHasTypeMismatch(t, docFile.Name(), variablesFile.Name(), tt.wantTypeMismatch, gotResult.typeMismatch)
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

func splitOutputMessages(t *testing.T, output []byte) validationResult {
	result := validationResult{}
	outputStrings := strings.Split(string(output), "\n")

	for _, oo := range outputStrings {
		switch {
		case strings.HasPrefix(oo, "Missing variable definition:"):
			result.missingDefinition = append(result.missingDefinition, oo)
		case strings.HasPrefix(oo, "Missing variable documentation:"):
			result.missingDocumentation = append(result.missingDocumentation, oo)
		case strings.HasPrefix(oo, "Type mismatch for variable:"):
			result.typeMismatch = append(result.typeMismatch, oo)
		}
	}

	return result
}

func assertHasMissingDocumentation(t *testing.T, filename string, got, want []string) {
	t.Helper()

	if len(got) != len(want) {
		t.Errorf("wanted %v but got %v", want, got)
	}

	for _, wantStr := range want {
		found := false

		completeWantString := fmt.Sprintf("Missing variable documentation: %q is not documented in %q", wantStr, filename)

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

func assertHasMissingDefinition(t *testing.T, filename string, got, want []string) {
	t.Helper()

	if len(got) != len(want) {
		t.Errorf("wanted %v but got %v", want, got)
	}

	for _, wantStr := range want {
		found := false

		completeWantString := fmt.Sprintf("Missing variable definition: %q is not defined in %q", wantStr, filename)
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

func assertHasTypeMismatch(t *testing.T, docFilename, defFilename string, want []validators.TypeMismatchResult, got []string) {
	t.Helper()

	if len(got) != len(want) {
		t.Errorf("wanted %v but got %v", want, got)
	}

	for _, tm := range want {
		found := false

		completeWantString := fmt.Sprintf("Type mismatch for variable: %q is documented as %q in %q but defined as %q in %q", tm.Name, tm.DocumentedType, docFilename, tm.DefinedType, defFilename)

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

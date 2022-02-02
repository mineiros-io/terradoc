package main_test

import (
	"fmt"
	"io/ioutil"
	"log"
	"os/exec"
	"strings"
	"testing"

	"github.com/madlambda/spells/assert"
	"github.com/mineiros-io/terradoc/test"
)

func TestValidateVariables(t *testing.T) {
	tests := []struct {
		desc                     string
		tfdocFixture             string
		variablesFixture         string
		wantMissingDocumentation []string
		wantMissingDefinition    []string
		wantTypeMismatch         []string
		wantError                bool
	}{
		{
			desc:             "when `variables.tf` and terradoc file have the same variables",
			tfdocFixture:     "validate/complete.tfdoc.hcl",
			variablesFixture: "validate/complete-variables.tf",
		},
		{
			desc:                     "when `variables.tf` has missing variables",
			tfdocFixture:             "validate/complete.tfdoc.hcl",
			variablesFixture:         "validate/missing-variables.tf",
			wantMissingDocumentation: []string{},
			wantMissingDefinition:    []string{"beer"},
			wantTypeMismatch:         []string{},
			wantError:                true,
		},
		{
			desc:                     "when terradoc file has missing variables",
			tfdocFixture:             "validate/missing-variables.tfdoc.hcl",
			variablesFixture:         "validate/complete-variables.tf",
			wantMissingDocumentation: []string{"beer"},
			wantMissingDefinition:    []string{},
			wantTypeMismatch:         []string{},
			wantError:                true,
		},
		{
			desc:                     "when `variables.tf` has type mismatch",
			tfdocFixture:             "validate/complete.tfdoc.hcl",
			variablesFixture:         "validate/type-mismatch.tf",
			wantMissingDocumentation: []string{},
			wantMissingDefinition:    []string{},
			wantTypeMismatch:         []string{"number"},
			wantError:                true,
		},
		{
			desc:                     "when `variables.tf` has type mismatch and missing variables",
			tfdocFixture:             "validate/complete.tfdoc.hcl",
			variablesFixture:         "validate/type-mismatch-with-missing.tf",
			wantMissingDocumentation: []string{},
			wantMissingDefinition:    []string{"beer"},
			wantTypeMismatch:         []string{"person"},
			wantError:                true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.desc, func(t *testing.T) {
			tfdoc := test.ReadFixture(t, tt.tfdocFixture)
			tfdocFile, err := ioutil.TempFile("", "terradoc-validate-tfdoc-")
			assert.NoError(t, err)
			_, err = tfdocFile.Write(tfdoc)
			assert.NoError(t, err)

			variables := test.ReadFixture(t, tt.variablesFixture)
			variablesFile, err := ioutil.TempFile("", "terradoc-validate-variables-")
			assert.NoError(t, err)
			_, err = variablesFile.Write(variables)
			assert.NoError(t, err)

			cmd := exec.Command(terradocBinPath, "validate", tfdocFile.Name(), "-v", variablesFile.Name())

			output, err := cmd.CombinedOutput()

			gotResult := splitOutputMessages(t, output)

			if tt.wantError {
				assertHasMissingDocumentation(t, tfdocFile.Name(), gotResult.missingDocumentation, tt.wantMissingDocumentation)
				assertHasMissingDefinition(t, variablesFile.Name(), gotResult.missingDefinition, tt.wantMissingDefinition)
				assertHasTypeMismatch(t, gotResult.typeMismatch, tt.wantTypeMismatch)
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
		case strings.HasPrefix(oo, "Variable type mismatch:"):
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

func assertHasTypeMismatch(t *testing.T, got, want []string) {
	t.Helper()

	if len(got) != len(want) {
		t.Errorf("wanted %v but got %v", want, got)
	}

	for _, wantStr := range want {
		found := false

		completeWantString := fmt.Sprintf("Variable type mismatch: %q", wantStr)

		for _, msg := range got {
			if msg == completeWantString {
				found = true

				break
			}
		}

		if !found {
			t.Errorf("wanted output to have type mismatch for %q but didn't find it", wantStr)
		}
	}
}

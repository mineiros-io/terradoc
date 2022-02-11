package test

import (
	"embed"
	"io/fs"
	"testing"

	"github.com/madlambda/spells/assert"
	"github.com/mineiros-io/terradoc/internal/entities"
	"github.com/mineiros-io/terradoc/internal/validators"
)

//go:embed testdata/*
var testDataFS embed.FS

func ReadFixture(t *testing.T, filename string) []byte {
	data, err := testDataFS.ReadFile("testdata/" + filename)
	assert.NoError(t, err)

	return data
}

func OpenFixture(t *testing.T, filename string) fs.File {
	f, err := testDataFS.Open("testdata/" + filename)
	assert.NoError(t, err)

	return f
}

func AssertEqualTypes(t *testing.T, want, got entities.Type) {
	t.Helper()

	assert.EqualStrings(t, want.TFType.String(), got.TFType.String())
	assert.EqualStrings(t, want.Label, got.Label)

	if want.Nested != nil {
		assert.EqualStrings(t, want.Nested.TFType.String(), got.Nested.TFType.String())
		assert.EqualStrings(t, want.Nested.Label, got.Nested.Label)
	}
}

func AssertHasTypeMismatches(t *testing.T, want, got []validators.TypeMismatchResult) {
	for _, tm := range want {
		found := false

		for _, tms := range got {
			if tms.Name == tm.Name {
				found = true

				assert.EqualStrings(t, tm.DefinedType, tms.DefinedType)
				assert.EqualStrings(t, tm.DocumentedType, tms.DocumentedType)
			}
		}

		if !found {
			t.Errorf("wanted %q to be found in %v", tm.Name, got)
		}
	}
}

func AssertHasStrings(t *testing.T, want, got []string) {
	t.Helper()

	for _, name := range want {
		found := false

		for _, o := range got {
			if name == o {
				found = true
			}
		}

		if !found {
			t.Errorf("wanted %q to be found in %v", name, got)
		}
	}
}

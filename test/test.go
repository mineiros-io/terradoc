package test

import (
	"embed"
	"io/fs"
	"testing"

	"github.com/madlambda/spells/assert"
	"github.com/mineiros-io/terradoc/internal/entities"
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

	if want.Nested == nil && got.Nested != nil {
		t.Fatalf("wanted nested to be nil but got %+v", got.Nested)
	}

	if want.Nested != nil {
		if got.Nested == nil {
			t.Fatal("wanted a nested type but found none")
		}

		assert.EqualStrings(t, want.Nested.TFType.String(), got.Nested.TFType.String())
		assert.EqualStrings(t, want.Nested.Label, got.Nested.Label)
	}

}

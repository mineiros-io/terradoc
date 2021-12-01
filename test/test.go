package test

import (
	"embed"
	"io/fs"
	"testing"

	"github.com/madlambda/spells/assert"
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

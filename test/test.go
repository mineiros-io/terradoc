package test

import (
	"embed"
	"io/fs"
	"testing"
)

//go:embed testdata/*
var testDataFS embed.FS

func ReadFixture(t *testing.T, filename string) []byte {
	data, err := testDataFS.ReadFile("testdata/" + filename)
	if err != nil {
		t.Fatal(err)
	}

	return data
}

func OpenFixture(t *testing.T, filename string) fs.File {
	f, err := testDataFS.Open("testdata/" + filename)
	if err != nil {
		t.Fatal(err)
	}

	return f
}

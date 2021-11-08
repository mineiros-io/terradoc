package test

import (
	"embed"
	"io/fs"
)

//go:embed testdata/*
var testDataFS embed.FS

func OpenFixture(filename string) (fs.File, error) {
	return testDataFS.Open("testdata/" + filename)
}

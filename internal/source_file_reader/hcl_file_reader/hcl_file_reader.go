package hcl_file_reader

import (
	"fmt"
	"path/filepath"

	"github.com/hashicorp/hcl/v2/hclparse"
	"github.com/mineiros-io/terradoc/internal/entities"
)

func Read(path string) (*entities.SourceFile, error) {
	absPath, err := getAbsPath(path)
	if err != nil {
		return nil, err
	}

	sf, err := readHCLFile(absPath)
	if err != nil {
		return nil, err
	}

	return sf, nil
}

func readHCLFile(filepath string) (*entities.SourceFile, error) {
	p := hclparse.NewParser()

	f, diags := p.ParseHCLFile(filepath)
	if diags.HasErrors() {
		return nil, diags
	}

	return &entities.SourceFile{HCLFile: f}, nil
}

func getAbsPath(path string) (string, error) {
	absPath, err := filepath.Abs(path)
	if err != nil {
		return "", fmt.Errorf("Could not generate absolute path for `%s`", path)
	}

	return absPath, nil
}

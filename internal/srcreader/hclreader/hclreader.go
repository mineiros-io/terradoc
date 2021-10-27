package hclreader

import (
	"errors"
	"fmt"
	"path/filepath"
	"strings"

	"github.com/hashicorp/hcl/v2"
	"github.com/hashicorp/hcl/v2/hclparse"
	"github.com/mineiros-io/terradoc/internal/entities"
)

func Read(src []byte) (*entities.SourceFile, error) {
	p := hclparse.NewParser()

	f, diags := p.ParseHCL(src, "")
	if diags.HasErrors() {
		return nil, diagsToError(diags)
	}

	return &entities.SourceFile{HCLFile: f}, nil
}

func ReadFromFile(path string) (*entities.SourceFile, error) {
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
		return nil, diagsToError(diags)
	}

	return &entities.SourceFile{HCLFile: f}, nil
}

func getAbsPath(path string) (string, error) {
	absPath, err := filepath.Abs(path)
	if err != nil {
		return "", fmt.Errorf("Could not get absolute path for `%s`", path)
	}

	return absPath, nil
}

func diagsToError(diags hcl.Diagnostics) error {
	var errsStr []string

	for _, err := range diags.Errs() {
		errsStr = append(errsStr, err.Error())
	}

	return errors.New(strings.Join(errsStr, "; "))
}

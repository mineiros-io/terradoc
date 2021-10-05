package main

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"log"

	"github.com/hashicorp/hcl/v2/hclparse"
	"github.com/mineiros-io/terradoc/internal/entities"
	"github.com/mineiros-io/terradoc/internal/source_file_reader/hcl_file_reader"
	"github.com/mineiros-io/terradoc/internal/writer_factory"
)

func getWriterFactory(outputDir string) writer_factory.WriterFactory {
	if outputDir == "" {
		return &writer_factory.STDOUTWriter{}
	}

	return &writer_factory.FileWriter{DirPath: outputDir}
}

func readInput(r io.Reader, args ...string) (*entities.SourceFile, error) {
	if len(args) > 0 {
		return readInputFromFile(args[0])
	}

	return readInputFromReader(r)
}

func readInputFromFile(filepath string) (*entities.SourceFile, error) {
	sf, err := hcl_file_reader.Read(filepath)
	if err != nil {
		log.Fatalf("Error reading file: %v", err)
	}

	return sf, nil
}

func readInputFromReader(r io.Reader) (*entities.SourceFile, error) {
	var result []byte
	s := bufio.NewScanner(r)

	s.Split(bufio.ScanBytes)

	for s.Scan() {
		if err := s.Err(); err != nil {
			return nil, err
		}

		result = append(result, s.Bytes()...)
	}

	if len(bytes.TrimSpace(result)) == 0 {
		return nil, fmt.Errorf("Input cannot be blank")
	}

	f, diags := hclparse.NewParser().ParseHCL(result, "")
	if diags.HasErrors() {
		return nil, diags
	}

	return &entities.SourceFile{HCLFile: f}, nil
}

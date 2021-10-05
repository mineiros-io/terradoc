package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/mineiros-io/terradoc/pkg/terradoc"
)

func main() {
	doReadme := flag.Bool("r", false, "Generate `README.md` file")
	doVariables := flag.Bool("v", false, "Generate `variables.tf` file")
	outputDir := flag.String("d", "", "Directory to write result files in")

	flag.Parse()

	if !(*doReadme || *doVariables) {
		flag.Usage()

		os.Exit(1)
	}

	wf := getWriterFactory(*outputDir)
	td := terradoc.NewTerradoc(wf, *doReadme, *doVariables)

	sf, err := readInput(os.Stdin, flag.Args()...)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error reading input: %s\n", err.Error())

		os.Exit(1)
	}

	err = td.CreateDocumentation(sf)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error creating documentation: %s\n", err.Error())

		os.Exit(1)
	}
}

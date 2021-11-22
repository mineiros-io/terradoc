package main

import (
	"flag"
	"fmt"
	"io"
	"os"

	"github.com/mineiros-io/terradoc/internal/parser/hclparser"
	"github.com/mineiros-io/terradoc/internal/renderers/markdown"
)

func main() {
	outputFile := flag.String("o", "", "Output file name")

	flag.Parse()

	r, rCloser, err := openInput(flag.Args()...)
	if err != nil {
		fmt.Fprintf(os.Stderr, "reading input: %s\n", err)

		os.Exit(1)
	}
	defer rCloser()

	def, err := hclparser.Parse(r, r.Name())
	if err != nil {
		fmt.Fprintf(os.Stderr, "parsing document: %v", err)

		os.Exit(1)
	}

	w, wCloser, err := getOutputWriter(*outputFile)
	if err != nil {
		fmt.Fprintf(os.Stderr, "creating writer: %s\n", err)

		os.Exit(1)
	}
	defer wCloser()

	err = markdown.Render(w, def)
	if err != nil {
		fmt.Fprintf(os.Stderr, "rendering document: %v", err)

		os.Exit(1)
	}
}

func getOutputWriter(filename string) (io.Writer, func(), error) {
	if filename == "" {
		return os.Stdout, noopClose, nil
	}

	f, err := os.Create(filename)
	if err != nil {
		return nil, nil, err
	}

	closer := func() {
		if err := f.Close(); err != nil {
			fmt.Fprintf(os.Stderr, "closing output stream: %v", err)
		}
	}

	return f, closer, nil
}

func openInput(args ...string) (*os.File, func(), error) {
	switch len(args) {
	case 0:
		return os.Stdin, noopClose, nil
	case 1:
		f, err := os.Open(args[0])
		if err != nil {
			return nil, nil, err
		}

		closer := func() {
			if err := f.Close(); err != nil {
				fmt.Fprintf(os.Stderr, "closing input stream: %v", err)
			}
		}

		return f, closer, nil
	default:
		return nil, nil, fmt.Errorf("expects none or one input file but given %d", len(args))
	}
}

func noopClose() {}

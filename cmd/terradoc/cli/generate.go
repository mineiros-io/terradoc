package cli

import (
	"fmt"
	"io"
	"os"

	"github.com/mineiros-io/terradoc/internal/parser/hclparser"
	"github.com/mineiros-io/terradoc/internal/renderers/markdown"
)

type GenerateCmd struct {
	InputFile  string `arg:"" required:"" help:"Input file." type:"existingfile"`
	OutputFile string `name:"output" short:"o" optional:"" default:"-" help:"Output file to write resulting markdown to" type:"path"`
}

func (g GenerateCmd) Run() error {
	r, rCloser, err := openInput(g.InputFile)
	if err != nil {
		return fmt.Errorf("opening input: %s\n", err)
	}
	defer rCloser()

	w, wCloser, err := getOutputWriter(g.OutputFile)
	if err != nil {
		return fmt.Errorf("creating writer: %s\n", err)
	}
	defer wCloser()

	def, err := hclparser.Parse(r, r.Name())
	if err != nil {
		return fmt.Errorf("parsing input: %v", err)
	}

	err = markdown.Render(w, def)
	if err != nil {
		return fmt.Errorf("rendering document: %v", err)
	}

	return nil
}

func openInput(path string) (*os.File, func(), error) {
	if path == "-" {
		return os.Stdin, noopClose, nil
	}

	f, err := os.Open(path)
	if err != nil {
		return nil, nil, err
	}

	closer := func() {
		if err := f.Close(); err != nil {
			fmt.Fprintf(os.Stderr, "closing input stream: %v", err)
		}
	}

	return f, closer, nil
}

func getOutputWriter(filename string) (io.Writer, func(), error) {
	if filename == "-" {
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

func noopClose() {}

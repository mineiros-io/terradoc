package cli

import (
	"fmt"
	"io/ioutil"
	"os"

	"github.com/hashicorp/hcl/v2/hclwrite"
)

type FormatCmd struct {
	InputFile string `arg:"" help:"Input file." type:"existingfile"`
	Write     bool   `name:"write" short:"w" help:"Overwrite file with formatted version."`
}

func (f FormatCmd) Run() error {
	inSrc, err := ioutil.ReadFile(f.InputFile)
	if err != nil {
		return fmt.Errorf("reading input: %s", err)
	}

	outSrc := hclwrite.Format(inSrc)

	if f.Write {
		err = ioutil.WriteFile(f.InputFile, outSrc, 0644)
	} else {
		_, err = os.Stdout.Write(outSrc)
	}

	if err != nil {
		return fmt.Errorf("writing result: %s", err)
	}

	return nil
}

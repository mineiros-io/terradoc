package main

import (
	"fmt"

	"github.com/mineiros-io/terradoc/config"
	"gopkg.in/alecthomas/kingpin.v2"
)

var (
	fileName = kingpin.Arg("file", "The .tf or .hcl file that contains the variable blocks that should be parsed").Required().String()
)

func main() {
	kingpin.Version("0.0.1")
	kingpin.Parse()

	module := config.ParseVariables(*fileName)
	outputBuffer := config.Format(module)

	// Write buffer as a string to stdout
	fmt.Print(outputBuffer.String())
}

// - **`region`**: *(Optional `string`)*
// If specified, the AWS region this bucket should reside in.
// Default is the region used by the callee.

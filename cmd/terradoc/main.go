package main

import (
	"github.com/alecthomas/kong"
	"github.com/mineiros-io/terradoc/cmd/terradoc/cli"
)

func main() {
	ctx := kong.Parse(&cli.Cli)
	err := ctx.Run()
	ctx.FatalIfErrorf(err)
}

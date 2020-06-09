package config

import (
	"log"

	"github.com/hashicorp/hcl/v2/hclparse"
)

func Load() {
	filename := "test.tf"

	parser := hclparse.NewParser()

	file, fileDiags := parser.ParseHCLFile(filename)

	content, _, contentDiags := file.Body.PartialContent(rootSchema)
	log.Print(content)
	log.Print(contentDiags)
	log.Print(fileDiags)

	for _, block := range content.Blocks {
		log.Print(string(block.Type))
	}
}

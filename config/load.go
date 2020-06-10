package config

import (
	"log"

	"github.com/hashicorp/hcl/v2/gohcl"
	"github.com/hashicorp/hcl/v2/hclparse"
)

func ParseVariables(filename string) *Module {

	module := newModule()
	parser := hclparse.NewParser()

	file, fileDiags := parser.ParseHCLFile(filename)

	if fileDiags.HasErrors() {
		log.Fatal(fileDiags.Error())
	}

	content, _, contentDiags := file.Body.PartialContent(rootSchema)

	if contentDiags.HasErrors() {
		log.Fatal(contentDiags.Error())
	}

	// loop over all found blocks
	for _, block := range content.Blocks {
		switch block.Type {
		// handle variable blocks
		case variableBlock:
			content, _, _ := block.Body.PartialContent(variableSchema)

			variableName := block.Labels[0]
			variable := &Variable{
				Name: variableName,
			}

			module.Variables[variableName] = variable

			// parse type
			if attr, defined := content.Attributes["type"]; defined {
				var typeExpr string

				typeDiags := gohcl.DecodeExpression(attr.Expr, nil, &typeExpr)

				if typeDiags.HasErrors() {
					rng := attr.Expr.Range()
					sourceFilename := rng.Filename
					source, exists := parser.Sources()[sourceFilename]

					if exists {
						typeExpr = string(rng.SliceBytes(source))
					} else {
						// This should never happen, so we'll just warn about
						// it and leave the type unspecified.
						typeExpr = ""
					}
				}

				variable.Type = typeExpr
			}

			// // parse description
			if attr, defined := content.Attributes["description"]; defined {
				var descriptionExpr string

				descDiags := gohcl.DecodeExpression(attr.Expr, nil, &descriptionExpr)

				if descDiags.HasErrors() {
					log.Fatal(descDiags.Error())
				}

				variable.Description = descriptionExpr
			}

		}
	}

	return module
}

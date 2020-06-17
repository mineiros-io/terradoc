package config

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/hashicorp/hcl/v2/gohcl"
	"github.com/hashicorp/hcl/v2/hclparse"
	ctyjson "github.com/zclconf/go-cty/cty/json"
)

// ParseVariables laods a given file by it's filename and
// loads its structure into a Module struct
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
		case "variable":
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

			// parse description
			if attr, defined := content.Attributes["description"]; defined {
				var descriptionExpr string

				descDiags := gohcl.DecodeExpression(attr.Expr, nil, &descriptionExpr)

				if descDiags.HasErrors() {
					log.Fatal(descDiags.Error())
				}

				variable.Description = descriptionExpr
			}

			// parse default value
			if attr, defined := content.Attributes["default"]; defined {
				// To avoid the caller needing to deal with cty here, we'll
				// use its JSON encoding to convert into an
				// approximately-equivalent plain Go interface{} value
				// to return.

				val, defaultDiags := attr.Expr.Value(nil)

				if defaultDiags != nil {
					log.Fatal(defaultDiags.Error())
				}

				if val.IsWhollyKnown() { // should only be false if there are errors in the input
					valJSON, err := ctyjson.Marshal(val, val.Type())
					if err != nil {
						// Should never happen, since all possible known
						// values have a JSON mapping.
						log.Fatal(fmt.Errorf("failed to serialize default value as JSON: %s", err))
					}
					var def interface{}
					err = json.Unmarshal(valJSON, &def)
					if err != nil {
						// Again should never happen, because valJSON is
						// guaranteed valid by ctyjson.Marshal.
						log.Fatal(fmt.Errorf("failed to re-parse default value from JSON: %s", err))
					}
					variable.Default = def
				}
			} else {
				variable.Required = true
			}

		}
	}

	return module
}

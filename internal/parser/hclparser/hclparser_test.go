package hclparser_test

import (
	"errors"
	"strings"
	"testing"

	"github.com/hashicorp/hcl/v2/hclparse"
	"github.com/mineiros-io/terradoc/internal/entities"
	"github.com/mineiros-io/terradoc/internal/parser/hclparser"
)

func TestParse(t *testing.T) {
	content := `
	section {
	  title       = "Module Argument Reference"
	  description = "See [variables.tf] and [examples/] for details and use-cases."

	  section {
	    title = "Main Resource Configuration"

	    variable "local_secondary_indexes" {
	      type        = any
	      readme_type = "list(local_secondary_index)"

	      description = "Describe an LSI on the table; these can only be allocated creation so you cannot change this definition after you have created the resource."
	      default     = []

	      required = true

	      forces_recreation = true

	      readme_example = {
	        local_secondary_indexes = [
	          {
	            range_key = "someKey"
	          }
	        ]
	      }

	      attribute "range_key" {
	        type = string

	        description = "The attribute to use as the range (sort) key. Must also be defined as an attribute, see below."

	        forces_recreation = true
	      }
	    }
	  }

	}
	`

	sf, err := generateSourceFile(content)
	if err != nil {
		t.Fatal(err)
	}

	// parsed definition
	definition, err := hclparser.Parse(sf)
	if err != nil {
		t.Fatal(err)
	}

	if len(definition.Sections) != 1 {
		t.Errorf("Expected definition to contain a section but found %d", len(definition.Sections))
	}

	rootSection := definition.Sections[0]
	t.Run("ParseRootSection", func(t *testing.T) {
		wantTitle := "Module Argument Reference"
		wantDesc := "See [variables.tf] and [examples/] for details and use-cases."
		wantNestedSectionsNum := 1

		if wantTitle != rootSection.Title {
			t.Errorf("Wanted root section title to be %q. Got %q instead", wantTitle, rootSection.Title)
		}

		if wantDesc != rootSection.Description {
			t.Errorf("Wanted root section description to be %q. Got %q instead", wantDesc, rootSection.Description)
		}

		if len(rootSection.Sections) != wantNestedSectionsNum {
			t.Errorf(
				"Expected root section to contain %d nested sections. Found %d instead",
				wantNestedSectionsNum,
				len(rootSection.Sections),
			)
		}
	})

	subSection := rootSection.Sections[0]
	t.Run("ParseSubSection", func(t *testing.T) {
		wantTitle := "Main Resource Configuration"
		wantDesc := ""
		wantVariablesNum := 1

		if wantTitle != subSection.Title {
			t.Errorf("Wanted root section title to be %q. Got %q instead", wantTitle, subSection.Title)
		}

		if wantDesc != subSection.Description {
			t.Errorf("Wanted root section description to be %q. Got %q instead", wantDesc, subSection.Description)
		}

		if len(subSection.Variables) != wantVariablesNum {
			t.Errorf(
				"Expected sub section to contain %d variable blocks. Found %d instead",
				wantVariablesNum,
				len(subSection.Variables),
			)
		}
	})

	variable := subSection.Variables[0]
	t.Run("ParseVariable", func(t *testing.T) {
		wantName := "local_secondary_indexes"
		wantType := "dynamic"
		wantDesc := "Describe an LSI on the table; these can only be allocated creation so you cannot change this definition after you have created the resource."
		wantDefault := "[]"
		wantRequired := true
		wantForcesRecreation := true

		// TODO: remove ugly hardcoded stuff
		wantReadmeType := "list(local_secondary_index)"
		wantReadmeExample := "local_secondary_indexes = [{\n    range_key = \"someKey\"\n  }]"
		wantAttrNum := 1

		if wantName != variable.Name {
			t.Errorf("Wanted variable name to be %q. Got %q instead", wantName, variable.Name)
		}

		if wantDesc != variable.Description {
			t.Errorf("Wanted variable description to be %q. Got %q instead", wantDesc, variable.Description)
		}

		if wantType != variable.Type {
			t.Errorf("Wanted variable type to be %q. Got %q instead", wantType, variable.Type)
		}

		if wantDefault != variable.Default {
			t.Errorf("Wanted variable default to be %q. Got %q instead", wantDefault, variable.Default)
		}

		if wantRequired != variable.Required {
			t.Errorf("Wanted variable required to be %t. Got %t instead", wantRequired, variable.Required)
		}

		if wantForcesRecreation != variable.ForcesRecreation {
			t.Errorf("Wanted variable forces recreation to be %t. Got %t instead", wantForcesRecreation, variable.ForcesRecreation)
		}

		if wantReadmeType != variable.ReadmeType {
			t.Errorf("Wanted variable readme type to be %q. Got %q instead", wantDefault, variable.Default)
		}

		if wantReadmeExample != variable.ReadmeExample {
			t.Errorf("Wanted variable readme example to be %q. Got %q instead", wantDefault, variable.Default)
		}

		if len(variable.Attributes) != wantAttrNum {
			t.Errorf(
				"Expected variable to contain %d attribute blocks. Found %d instead",
				wantAttrNum,
				len(variable.Attributes),
			)
		}
	})

	// attribute
	attribute := variable.Attributes[0]
	t.Run("ParseAttribute", func(t *testing.T) {
		wantName := "range_key"
		wantDesc := "The attribute to use as the range (sort) key. Must also be defined as an attribute, see below."
		wantType := "string"
		wantForcesRecreation := true
		wantRequired := false
		wantNestedAttrNum := 0

		if wantName != attribute.Name {
			t.Errorf("Wanted attribute name to be %q. Got %q instead", wantName, attribute.Name)
		}

		if wantDesc != attribute.Description {
			t.Errorf("Wanted attribute description to be %q. Got %q instead", wantDesc, attribute.Description)
		}

		if wantType != attribute.Type {
			t.Errorf("Wanted attrybute type to be %q. Got %q instead", wantType, attribute.Type)
		}

		if wantRequired != attribute.Required {
			t.Errorf("Wanted attribute required to be %t. Got %t instead", wantRequired, attribute.Required)
		}

		if wantForcesRecreation != attribute.ForcesRecreation {
			t.Errorf("Wanted attribute forces recreation to be %t. Got %t instead", wantForcesRecreation, attribute.ForcesRecreation)
		}

		if len(attribute.Attributes) != wantNestedAttrNum {
			t.Errorf(
				"Expected attribute to contain %d nested attribute blocks. Found %d instead",
				wantNestedAttrNum,
				len(attribute.Attributes),
			)
		}
	})
}

func generateSourceFile(content string) (*entities.SourceFile, error) {
	p := hclparse.NewParser()

	f, diag := p.ParseHCL([]byte(content), "")
	// TODO: refactor
	if diag.HasErrors() {
		var errStr []string

		for _, err := range diag.Errs() {
			errStr = append(errStr, err.Error())
		}

		return nil, errors.New(strings.Join(errStr, "; "))
	}

	return &entities.SourceFile{HCLFile: f}, nil
}

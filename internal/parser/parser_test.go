package parser

import (
	"errors"
	"strings"
	"fmt"
	"encoding/json"
	"testing"

	"github.com/hashicorp/hcl/v2/hclparse"
	"github.com/mineiros-io/terradoc/internal/entities"
	"github.com/stretchr/testify/require"
)

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
	require.NoError(t, err)

	result, err := Parse(sf)
	require.NoError(t, err)

	require.Len(t, result.Sections, 1)

	section := result.Sections[0]

	require.Equal(t, "Module Argument Reference", section.Title)
	require.Equal(t, "See [variables.tf] and [examples/] for details and use-cases.", section.Description)

	require.Len(t, section.Sections, 1)

	subSection := section.Sections[0]

	require.Equal(t, "Main Resource Configuration", subSection.Title)
	require.Equal(t, "", subSection.Description)

	require.Len(t, subSection.Variables, 1)

	variable := subSection.Variables[0]

	require.Equal(t, "local_secondary_indexes", variable.Name)
	require.Equal(t, "dynamic", variable.Type)
	require.Equal(t, "Describe an LSI on the table; these can only be allocated creation so you cannot change this definition after you have created the resource.", variable.Description)
	require.Equal(t, "[]", variable.Default)
	require.True(t, variable.Required)
	require.True(t, variable.ForcesRecreation)
	// TODO: remove ugly hardcoded stuff
	require.Equal(t, "list(local_secondary_index)", variable.ReadmeType)
	require.Equal(t, `local_secondary_indexes = [{\nrange_key = "someKey"\n}]`, variable.ReadmeExample)


	require.Len(t, variable.Attributes, 1)

	attribute := variable.Attributes[0]

	require.Equal(t, "The attribute to use as the range (sort) key. Must also be defined as an attribute, see below.", attribute.Description)
	require.Equal(t, "string", attribute.Type)
	require.Equal(t, "range_key", attribute.Name)
	require.True(t, attribute.ForcesRecreation)
	require.False(t, attribute.Required)

	// TODO: remove
	jsonStuff, err := json.MarshalIndent(result, "", "  ")
	require.NoError(t, err)

	fmt.Print(string(jsonStuff))
}

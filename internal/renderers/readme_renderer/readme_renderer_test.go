package readme_renderer

import (
	"bytes"
	"errors"
	"strings"
	"testing"

	"github.com/hashicorp/hcl/v2/hclparse"
	"github.com/mineiros-io/terradoc/internal/entities"
	"github.com/mineiros-io/terradoc/internal/parser"
	"github.com/stretchr/testify/require"
)

func renderSourceFile(content string) (*entities.SourceFile, error) {
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

func TestRender(t *testing.T) {
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

	sf, err := renderSourceFile(content)
	require.NoError(t, err)

	result, err := parser.Parse(sf)
	require.NoError(t, err)

	buf := new(bytes.Buffer)

	err = Render(buf, result)
	require.NoError(t, err)

	require.Equal(t, buf.String(), "foo")
}

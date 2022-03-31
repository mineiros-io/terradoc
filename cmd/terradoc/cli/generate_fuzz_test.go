//go:build go1.18

package cli

import (
	"bytes"
	"strings"
	"testing"

	"github.com/mineiros-io/terradoc/internal/parsers/docparser"
	"github.com/mineiros-io/terradoc/internal/renderers/markdown"
)

func FuzzGenerate(f *testing.F) {
	seedCorpus := []string{
		`section {
			title = "test"
			variable "module_enabled" {
				type        = bool
			}
		}`,
		`section {
			title = "test"
			section {
				a = 1
			}
		}`,
		`section {
			variable "local_secondary_indexes" {
				type        = any
				readme_type = "list(local_secondary_index)"
			}
		}`,
	}

	for _, seed := range seedCorpus {
		f.Add(seed)
	}

	f.Fuzz(func(t *testing.T, str string) {
		r := strings.NewReader(str)
		w := bytes.Buffer{}

		def, err := docparser.Parse(r, "in.hcl")
		if err != nil {
			return
		}

		_ = markdown.Render(&w, def)
	})
}

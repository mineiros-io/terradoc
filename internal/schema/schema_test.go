package schema_test

import (
	"reflect"
	"testing"

	"github.com/hashicorp/hcl/v2"
	"github.com/mineiros-io/terradoc/internal/schema"
)

func assertHasBlock(t *testing.T, blocks []hcl.BlockHeaderSchema, expectedType string, expectedLabels ...string) {
	t.Helper()

	var found bool

	for _, blk := range blocks {
		if blk.Type == expectedType {
			found = true

			if !reflect.DeepEqual(blk.LabelNames, expectedLabels) {
				t.Errorf("Expected block to have labels %q. Got %q instead", expectedLabels, blk.LabelNames)
			}
		}
	}

	if !found {
		t.Errorf("Expected to have block with type %q but none was found", expectedType)
	}
}

func TestRootSchema(t *testing.T) {
	s := schema.RootSchema

	if len(s.Attributes) > 0 {
		t.Errorf("Expected root schema to not have attributes. Found %+v instead", s.Attributes)
	}

	assertHasBlock(t, s.Blocks, "section")
}

func TestSectionSchema(t *testing.T) {
	s := schema.SectionSchema

	assertHasBlock(t, s.Blocks, "section")
	assertHasBlock(t, s.Blocks, "variable", "name")

}

func TestVariableSchema(t *testing.T) {
	s := schema.VariableSchema

	assertHasBlock(t, s.Blocks, "attribute", "name")

}

func AttributeSchema(t *testing.T) {
	s := schema.AttributeSchema

	assertHasBlock(t, s.Blocks, "attribute", "name")
}

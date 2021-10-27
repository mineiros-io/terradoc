package hclschema_test

import (
	"testing"

	"github.com/hashicorp/hcl/v2"
	"github.com/mineiros-io/terradoc/internal/parser/hclparser/hclschema"
)

func TestRootSchema(t *testing.T) {
	s := hclschema.RootSchema()

	if len(s.Attributes) > 0 {
		t.Errorf("Expected root schema to not have attributes. Found %+v instead", s.Attributes)
	}

	sectionBlocks := getBlocks(s.Blocks, "section")

	if len(sectionBlocks) != 1 {
		t.Errorf("Expected 1 section block. Found %d instead", len(sectionBlocks))
	}
}

func TestSectionSchema(t *testing.T) {
	s := hclschema.SectionSchema()

	// schema attributes
	assertHasAttribute(t, s, "title", true)
	assertHasAttribute(t, s, "description", false)

	// schema blocks
	nestedSectionBlocks := getBlocks(s.Blocks, "section")
	if len(nestedSectionBlocks) != 1 {
		t.Errorf("Expected 1 nested section block. Found %d instead", len(nestedSectionBlocks))
	}
	assertDoesNotHaveLabel(t, nestedSectionBlocks[0])

	variableBlocks := getBlocks(s.Blocks, "variable")
	if len(variableBlocks) != 1 {
		t.Errorf("Expected 1 variable block. Found %d instead", len(variableBlocks))
	}

	assertBlockHasLabel(t, variableBlocks[0], "name")
}

func TestVariableSchema(t *testing.T) {
	s := hclschema.VariableSchema()

	// schema attributes
	assertHasAttribute(t, s, "type", true)
	assertHasAttribute(t, s, "readme_type", false)
	assertHasAttribute(t, s, "description", false)
	assertHasAttribute(t, s, "default", false)
	assertHasAttribute(t, s, "required", false)
	assertHasAttribute(t, s, "forces_recreation", false)
	assertHasAttribute(t, s, "readme_example", false)

	// schema blocks
	attrBlocks := getBlocks(s.Blocks, "attribute")
	if len(attrBlocks) != 1 {
		t.Errorf("Expected 1 attribute block. Found %d instead", len(attrBlocks))
	}
	assertBlockHasLabel(t, attrBlocks[0], "name")
}

func TestAttributeSchema(t *testing.T) {
	s := hclschema.AttributeSchema()

	// schema attributes
	assertHasAttribute(t, s, "type", true)
	assertHasAttribute(t, s, "description", false)
	assertHasAttribute(t, s, "required", false)
	assertHasAttribute(t, s, "forces_recreation", false)

	// schema blocks
	attrBlocks := getBlocks(s.Blocks, "attribute")
	if len(attrBlocks) != 1 {
		t.Errorf("Expected 1 attribute block. Found %d instead", len(attrBlocks))
	}
	assertBlockHasLabel(t, attrBlocks[0], "name")
}

func getBlocks(blockList []hcl.BlockHeaderSchema, blockType string) (result []hcl.BlockHeaderSchema) {
	for _, blk := range blockList {
		if blk.Type == blockType {
			result = append(result, blk)
		}
	}

	return result
}

func assertDoesNotHaveLabel(t *testing.T, blk hcl.BlockHeaderSchema) {
	if len(blk.LabelNames) != 0 {
		t.Errorf("Expected block %q to not have labels. Found %v instead", blk.Type, blk.LabelNames)
	}
}

func assertBlockHasLabel(t *testing.T, blk hcl.BlockHeaderSchema, labelName string) {
	t.Helper()

	if len(blk.LabelNames) == 0 {
		t.Errorf("Expected block %q to have label %q but it was not found", blk.Type, labelName)
	}

	var found bool

	for _, label := range blk.LabelNames {
		if label == labelName {
			found = true
		}
	}

	if !found {
		t.Errorf("Expected block %q to have label %q but existing labels are %v", blk.Type, labelName, blk.LabelNames)
	}
}

func assertHasAttribute(t *testing.T, s *hcl.BodySchema, attrName string, isRequired bool) {
	t.Helper()

	if len(s.Attributes) == 0 {
		t.Errorf("Expected attribute %q to exist but no attribute is defined", attrName)
	}

	var found bool

	for _, attr := range s.Attributes {
		if attr.Name == attrName {
			found = true

			if attr.Required != isRequired {
				t.Errorf(
					"Expected attribute %q to have required as %t. Got %t instead",
					attrName,
					isRequired,
					attr.Required,
				)
			}
		}
	}

	if !found {
		t.Errorf("Attribute %q not found", attrName)
	}
}

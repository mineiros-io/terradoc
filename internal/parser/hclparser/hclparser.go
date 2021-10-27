package hclparser

import (
	"encoding/json"
	"log"
	"strings"

	"github.com/hashicorp/hcl/v2"
	"github.com/hashicorp/hcl/v2/ext/typeexpr"
	"github.com/hashicorp/hcl/v2/hclwrite"
	"github.com/mineiros-io/terradoc/internal/entities"
	"github.com/mineiros-io/terradoc/internal/parser/hclparser/hclschema"
	ctyJson "github.com/zclconf/go-cty/cty/json"
)

func Parse(sourceFile *entities.SourceFile) (*entities.Definition, error) {
	// TODO: check what to do here
	return parseRoot(sourceFile.HCLBody())
}

func parseRoot(body hcl.Body) (*entities.Definition, error) {
	rootContent, diags := body.Content(hclschema.RootSchema)
	if diags.HasErrors() {
		return nil, diags
	}

	def := &entities.Definition{}
	// Root does not have attributes and only has `section` blocks
	for _, sectionBlock := range rootContent.Blocks {
		section, err := parseSection(sectionBlock, 1) // initial level
		if err != nil {
			return nil, err
		}

		def.Sections = append(def.Sections, section)
	}

	return def, nil
}

func parseSection(sectionBlock *hcl.Block, level int) (*entities.Section, error) {
	sectionContent, diags := sectionBlock.Body.Content(hclschema.SectionSchema)
	if diags.HasErrors() {
		log.Print("SECTION")
		log.Fatal(diags.Errs())
	}

	section := &entities.Section{
		Level:       level,
		Title:       getString(sectionContent, "title"),
		Description: getString(sectionContent, "description"),
	}

	for _, varBlk := range sectionContent.Blocks.OfType("variable") {
		variable, err := parseVariable(varBlk)
		if err != nil {
			log.Fatal(err)
		}

		section.Variables = append(section.Variables, variable)
	}

	for _, secBlk := range sectionContent.Blocks.OfType("section") {
		nestedSection, err := parseSection(secBlk, level+1)
		if err != nil {
			log.Fatal(err)
		}

		section.Sections = append(section.Sections, nestedSection)
	}

	return section, nil
}

func parseVariable(variableBlock *hcl.Block) (*entities.Variable, error) {
	name := variableBlock.Labels[0]

	variableContent, diags := variableBlock.Body.Content(hclschema.VariableSchema)
	if diags.HasErrors() {
		return nil, diags
	}

	variable := &entities.Variable{
		Name:             name,
		Type:             getType(variableContent, "type"),
		Description:      getString(variableContent, "description"),
		ReadmeType:       getString(variableContent, "readme_type"),
		Default:          getJSON(variableContent, "default"),
		Required:         getBool(variableContent, "required"),
		ForcesRecreation: getBool(variableContent, "forces_recreation"),
		ReadmeExample:    getHCL(variableContent, "readme_example"),
	}

	// variables have only `attribute` blocks
	for _, blk := range variableContent.Blocks.OfType("attribute") {
		attribute, err := parseAttribute(blk, 1)
		if err != nil {
			return nil, err
		}

		variable.Attributes = append(variable.Attributes, attribute)
	}

	return variable, nil
}

func parseAttribute(attrBlock *hcl.Block, level int) (*entities.Attribute, error) {
	name := attrBlock.Labels[0]

	attrContent, diags := attrBlock.Body.Content(hclschema.AttributeSchema)
	if diags.HasErrors() {
		return nil, diags
	}

	attr := &entities.Attribute{
		Name:             name,
		Description:      getString(attrContent, "description"),
		Required:         getBool(attrContent, "required"),
		ForcesRecreation: getBool(attrContent, "forces_recreation"),
		Type:             getType(attrContent, "type"),
		Level:            level,
	}

	// attribute blocks have only `attribute` blocks
	for _, blk := range attrContent.Blocks.OfType("attribute") {
		nestedAttr, err := parseAttribute(blk, level+1)
		if err != nil {
			return nil, err
		}

		attr.Attributes = append(attr.Attributes, nestedAttr)
	}

	return attr, nil
}

func getString(content *hcl.BodyContent, attrName string) string {
	attr, exists := content.Attributes[attrName]
	if !exists {
		return ""
	}

	val, diags := attr.Expr.Value(nil)
	if diags.HasErrors() {
		log.Fatal(diags.Errs())
	}

	return val.AsString()
}

func getBool(content *hcl.BodyContent, attrName string) bool {
	attr, exists := content.Attributes[attrName]
	if !exists {
		return false
	}

	val, diags := attr.Expr.Value(nil)
	if diags.HasErrors() {
		log.Fatal(diags.Errs())
	}

	return val.True()
}

func getHCL(content *hcl.BodyContent, attrName string) string {
	attr, exists := content.Attributes[attrName]
	if !exists {
		return ""
	}

	val, diags := attr.Expr.Value(nil)
	if diags.HasErrors() {
		log.Printf("VAL ERROR: %+v", diags.Errs())
	}
	tk := hclwrite.TokensForValue(val)

	str := string(tk.Bytes())

	return strings.Trim(str, "{}\n ")
}

func getJSON(content *hcl.BodyContent, attrName string) string {
	attr, exists := content.Attributes[attrName]
	if !exists {
		return ""
	}

	v, _ := attr.Expr.Value(nil)
	jj := ctyJson.SimpleJSONValue{Value: v}

	result, err := json.MarshalIndent(jj, "", "  ") //jj.MarshalJSON()
	if err != nil {
		log.Fatal(err)
	}

	return string(result)
}

func getType(content *hcl.BodyContent, attrName string) string {
	attr, exists := content.Attributes[attrName]
	if !exists {
		return ""
	}
	val, diags := typeexpr.TypeConstraint(attr.Expr)
	if diags.HasErrors() {
		log.Fatal(diags.Errs())
	}

	return val.FriendlyName()
}

package markdown

import (
	"fmt"
	"io"
	"text/template"

	"github.com/mineiros-io/terradoc"
	"github.com/mineiros-io/terradoc/internal/entities"
	"github.com/mineiros-io/terradoc/internal/renderers"
)

const (
	templateName = "README.md"

	sectionTemplateName         = "section"
	referencesTemplateName      = "references"
	variableTemplateName        = "variable"
	attributeTemplateName       = "attribute"
	typeDescriptionTemplateName = "typeDescription"
	headerTemplateName          = "header"
	tocTemplateName             = "toc"
	outputTemplateName          = "output"

	varNestingLevel = 0
)

type markdownWriter struct {
	writer io.Writer
	templ  *template.Template
}

func newMarkdownWriter(writer io.Writer) (*markdownWriter, error) {
	const templatesPath = "templates/markdown/*"

	t, err := template.New(templateName).Funcs(renderers.TemplatesFuncMap).ParseFS(terradoc.TemplateFS, templatesPath)
	if err != nil {
		return nil, err
	}

	return &markdownWriter{writer: writer, templ: t}, nil
}

func (mw *markdownWriter) writeDefinition(definition entities.Doc) error {
	if err := mw.writeHeader(definition.Header); err != nil {
		return err
	}

	if err := mw.writeSections(definition.Sections); err != nil {
		return err
	}

	return mw.writeReferences(definition.References)
}

func (mw *markdownWriter) writeHeader(header entities.Header) error {
	// Prevent empty header from being rendered
	if len(header.Badges) > 0 || header.Image != "" {
		return mw.writeTemplate(headerTemplateName, header)
	}

	return nil
}

func (mw *markdownWriter) writeReferences(references []entities.Reference) error {
	if len(references) == 0 {
		return nil
	}

	return mw.writeTemplate(referencesTemplateName, references)
}

func (mw *markdownWriter) writeSections(sections []entities.Section) error {
	for _, section := range sections {
		if err := mw.writeSection(section); err != nil {
			return err
		}
	}

	return nil
}

func (mw *markdownWriter) writeSection(section entities.Section) error {
	if err := mw.writeTemplate(sectionTemplateName, section); err != nil {
		return err
	}

	if section.TOC {
		if err := mw.writeTOC(section.SubSections); err != nil {
			return err
		}
	}

	if err := mw.writeVariables(section.Variables); err != nil {
		return err
	}

	if err := mw.writeOutputs(section.Outputs); err != nil {
		return err
	}

	return mw.writeSections(section.SubSections)
}

func (mw *markdownWriter) writeTemplate(templateName string, v interface{}) error {
	return mw.templ.ExecuteTemplate(mw.writer, templateName, v)
}

func (mw *markdownWriter) writeVariables(variables []entities.Variable) error {
	for _, variable := range variables {
		if err := mw.writeVariable(variable); err != nil {
			return err
		}
	}

	return nil
}

func (mw *markdownWriter) writeVariable(variable entities.Variable) error {
	if err := mw.writeTemplate(variableTemplateName, variable); err != nil {
		return err
	}

	if len(variable.Attributes) > 0 {
		if err := mw.writeType(variable.Type, varNestingLevel); err != nil {
			return err
		}

		return mw.writeAttributes(variable.Attributes, variable.Name)
	}

	return nil
}

func (mw *markdownWriter) writeType(typeDefinition entities.Type, nestingLevel int) error {
	type typeRenderer struct {
		entities.Type
		IndentLevel int
	}

	indentLevel := renderers.GetIndent(nestingLevel)

	return mw.writeTemplate(
		typeDescriptionTemplateName,
		&typeRenderer{
			Type:        typeDefinition,
			IndentLevel: indentLevel,
		},
	)
}

func (mw *markdownWriter) writeAttributes(attributes []entities.Attribute, parentName string) error {
	for _, attribute := range attributes {
		if err := mw.writeAttribute(attribute, parentName); err != nil {
			return err
		}
	}

	return nil
}

func (mw *markdownWriter) writeAttribute(attribute entities.Attribute, parentName string) error {
	type attributeRenderer struct {
		entities.Attribute
		ParentName string
	}

	attrRenderer := attributeRenderer{Attribute: attribute, ParentName: parentName}

	if err := mw.writeTemplate(attributeTemplateName, attrRenderer); err != nil {
		return err
	}

	if len(attribute.Attributes) > 0 {
		if err := mw.writeType(attribute.Type, attribute.Level); err != nil {
			return err
		}

		nestedParentName := fmt.Sprintf("%s-%s", parentName, attribute.Name)

		return mw.writeAttributes(attribute.Attributes, nestedParentName)
	}

	return nil
}

type tocItemRenderer struct {
	Label       string
	IndentLevel int
}

func (mw *markdownWriter) writeTOC(sections []entities.Section) error {
	items := fetchTOCItems(sections, 0)

	return mw.writeTemplate(tocTemplateName, items)
}

func fetchTOCItems(sections []entities.Section, level int) (items []tocItemRenderer) {
	for _, section := range sections {

		tocItem := tocItemRenderer{Label: section.Title, IndentLevel: level}
		items = append(items, tocItem)

		nestedItems := fetchTOCItems(section.SubSections, level+2)

		items = append(items, nestedItems...)
	}

	return items
}

func (mw *markdownWriter) writeOutputs(outputs []entities.Output) error {
	for _, output := range outputs {
		if err := mw.writeOutput(output); err != nil {
			return err
		}
	}

	return nil
}

func (mw *markdownWriter) writeOutput(output entities.Output) error {
	if err := mw.writeTemplate(outputTemplateName, output); err != nil {
		return err
	}

	return nil
}

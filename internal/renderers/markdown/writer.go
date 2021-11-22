package markdown

import (
	"io"
	"log"
	"regexp"
	"strings"
	"text/template"

	"github.com/mineiros-io/terradoc"
	"github.com/mineiros-io/terradoc/internal/entities"
	"github.com/mineiros-io/terradoc/internal/renderers"
)

const (
	templateName = "README.md"

	headerTemplateName          = "header"
	badgeTemplateName           = "badge"
	referenceTemplateName       = "reference"
	sectionTemplateName         = "section"
	variableTemplateName        = "variable"
	attributeTemplateName       = "attribute"
	typeDescriptionTemplateName = "typeDescription"
	tocTemplateName             = "toc"
	rootSectionTemplateName     = "rootSection"

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

func (mw *markdownWriter) writeDefinition(definition entities.Definition) error {
	if err := mw.writeHeader(definition.Header); err != nil {
		return err
	}

	if err := mw.writeSections(definition.Sections); err != nil {
		return err
	}

	if err := mw.writeReferences(definition.References); err != nil {
		return err
	}

	return nil
}

func (mw *markdownWriter) writeRootSection(section entities.Section) error {
	if err := mw.writeTemplate(rootSectionTemplateName, section); err != nil {
		return err
	}

	if err := mw.writeTOC(section.SubSections); err != nil {
		return err
	}

	return mw.writeSections(section.SubSections)
}

func (mw *markdownWriter) writeReferences(references []entities.Reference) error {
	if len(references) == 0 {
		return nil
	}

	for _, reference := range references {
		if err := mw.writeTemplate(referenceTemplateName, reference); err != nil {
			return err
		}
	}

	return nil
}

func (mw *markdownWriter) writeHeader(header entities.Header) error {
	if err := mw.writeTemplate(headerTemplateName, header); err != nil {
		return err
	}

	return mw.writeBadges(header.Badges)
}

func (mw *markdownWriter) writeBadges(badges []entities.Badge) error {
	if len(badges) == 0 {
		return nil
	}

	for _, badge := range badges {
		if err := mw.writeTemplate(badgeTemplateName, badge); err != nil {
			return err
		}
	}

	return nil
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

	if err := mw.writeVariables(section.Variables); err != nil {
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

		return mw.writeAttributes(variable.Attributes)
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

func (mw *markdownWriter) writeAttributes(attributes []entities.Attribute) error {
	for _, attribute := range attributes {
		if err := mw.writeAttribute(attribute); err != nil {
			return err
		}
	}

	return nil
}

func (mw *markdownWriter) writeAttribute(attribute entities.Attribute) error {
	if err := mw.writeTemplate(attributeTemplateName, attribute); err != nil {
		return err
	}

	if len(attribute.Attributes) > 0 {
		if err := mw.writeType(attribute.Type, attribute.Level); err != nil {
			return err
		}

		return mw.writeAttributes(attribute.Attributes)
	}

	return nil
}

func (mw *markdownWriter) writeTOC(sections []entities.Section) error {
	items := fetchTOCItems(sections)

	return mw.writeTemplate(tocTemplateName, items)
}

type tocItemRenderer struct {
	Label       string
	Value       string
	IndentLevel int
}

func fetchTOCItems(sections []entities.Section) (items []tocItemRenderer) {
	for _, section := range sections {
		reg, err := regexp.Compile("[^a-zA-Z0-9 ]+")
		if err != nil {
			log.Fatal(err)
		}
		str := reg.ReplaceAllString(section.Title, "")

		value := strings.ReplaceAll(strings.ToLower(str), " ", "-")

		tocItem := tocItemRenderer{Label: section.Title, IndentLevel: section.Level, Value: value}

		nestedItems := fetchTOCItems(section.SubSections)

		items = append(items, append(nestedItems, tocItem)...)
	}

	return items
}

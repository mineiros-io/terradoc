package terradoc

import (
	"github.com/mineiros-io/terradoc/internal/entities"
	"github.com/mineiros-io/terradoc/internal/parser"
	"github.com/mineiros-io/terradoc/internal/renderers/readme_renderer"
	"github.com/mineiros-io/terradoc/internal/renderers/tfvariables_renderer"
	"github.com/mineiros-io/terradoc/internal/writer_factory"
)

const (
	readmeFileName    = "README.md"
	variablesFileName = "variables.tf"
)

func CreateTerraformDoc(wf writer_factory.WriterFactory, sf *entities.SourceFile) error {
	definition, err := parser.Parse(sf)
	if err != nil {
		return err
	}

	return createTerraformDoc(wf, definition)
}

func createTerraformDoc(wf writer_factory.WriterFactory, definition *entities.Definition) error {
	variablesFile, err := wf.NewWriter(variablesFileName)
	if err != nil {
		return err
	}
	defer variablesFile.Close()

	err = tfvariables_renderer.Render(variablesFile, definition)
	if err != nil {
		return err
	}

	return nil
}

func CreateReadmeDoc(wf writer_factory.WriterFactory, sf *entities.SourceFile) error {
	definition, err := parser.Parse(sf)
	if err != nil {
		return err
	}

	return createReadmeDoc(wf, definition)
}

func createReadmeDoc(wf writer_factory.WriterFactory, definition *entities.Definition) error {
	readmeFile, err := wf.NewWriter(readmeFileName)
	if err != nil {
		return err
	}
	defer readmeFile.Close()

	err = readme_renderer.Render(readmeFile, definition)
	if err != nil {
		return err
	}

	return nil
}

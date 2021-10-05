package terradoc

import (
	"github.com/mineiros-io/terradoc/internal/entities"
	"github.com/mineiros-io/terradoc/internal/writer_factory"
)

type Terradoc struct {
	doReadme      bool
	doVariables   bool
	writerFactory writer_factory.WriterFactory
}

func NewTerradoc(wf writer_factory.WriterFactory, doReadme bool, doVariables bool) *Terradoc {
	return &Terradoc{doReadme: doReadme, doVariables: doVariables, writerFactory: wf}
}

func (t *Terradoc) CreateDocumentation(sourceFile *entities.SourceFile) (err error) {
	if t.doReadme {
		err = CreateReadmeDoc(t.writerFactory, sourceFile)

		if err != nil {
			return err
		}
	}

	if t.doVariables {
		err = CreateTerraformDoc(t.writerFactory, sourceFile)

		if err != nil {
			return err
		}
	}

	return
}

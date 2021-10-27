package srcreader

import "github.com/mineiros-io/terradoc/internal/entities"

type SourceReader interface {
	ReadFromFile(string) (*entities.SourceFile, error)
	Read([]byte) (*entities.SourceFile, error)
}

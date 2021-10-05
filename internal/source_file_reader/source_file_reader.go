package source_file_reader

import "github.com/mineiros-io/terradoc/internal/entities"

type SourceFileReader interface {
	Read(interface{}) (*entities.SourceFile, error)
}

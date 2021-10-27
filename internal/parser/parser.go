package parser

import "github.com/mineiros-io/terradoc/internal/entities"

type Parser interface {
	Parse(*entities.SourceFile) (*entities.Definition, error)
}

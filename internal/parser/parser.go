package parser

import (
	"io"

	"github.com/mineiros-io/terradoc/internal/entities"
)

type Parser interface {
	Parse(io.Reader) (*entities.Definition, error)
}

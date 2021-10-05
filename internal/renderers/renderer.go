package renderers

import "github.com/mineiros-io/terradoc/internal/entities"

type Renderer interface {
	Render(*entities.Definition) error
}

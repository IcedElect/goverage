package renderer

import (
	"github.com/IcedElect/goverage/internal/structure/elements"
	"github.com/IcedElect/goverage/internal/structure/files"
	"github.com/IcedElect/goverage/internal/structure/tree"
)

type Renderer struct {
	elementsRegistry *elements.Registry
}

func NewRenderer(elementsRegistry *elements.Registry) *Renderer {
	return &Renderer{
		elementsRegistry: elementsRegistry,
	}
}

func (r *Renderer) RenderFile(file *files.File, outputDir string) error {
	return nil
}

func (r *Renderer) RenderDirectory(dir tree.Directory, outputDir string) error {
	return nil
}

package strategies

import (
	"github.com/IcedElect/goverage/internal/structure/elements"
	"github.com/IcedElect/goverage/internal/structure/files"
	"github.com/IcedElect/goverage/internal/structure/tree"
)

type Strategy interface {
	Name() string
	Execute(
		directories []tree.Directory,
		filesRegistry *files.Registry,
		elementsRegistry *elements.Registry,
		outputDir string,
	) error
}

type Registry struct {
	strategies map[string]Strategy
}

func NewRegistry(strategies ...Strategy) *Registry {
	strategiesMap := make(map[string]Strategy)
	for _, strategy := range strategies {
		strategiesMap[strategy.Name()] = strategy
	}

	return &Registry{
		strategies: strategiesMap,
	}
}

func (r *Registry) Get(name string) (Strategy, bool) {
	strategy, ok := r.strategies[name]
	return strategy, ok
}

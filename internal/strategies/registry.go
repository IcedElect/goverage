package strategies

import "golang.org/x/tools/cover"

type Strategy interface {
	Name() string
	Execute(profiles []*cover.Profile, outputDir string) error
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
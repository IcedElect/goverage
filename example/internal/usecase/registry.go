package usecase

import "github.com/IcedElect/goverage/example/internal/usecase/some"

type Registry struct {
	someUsecase some.Usecase
}

func NewRegistry() *Registry {
	return &Registry{
		someUsecase: some.NewUsecase(),
	}
}

func (r *Registry) SomeUsecase() some.Usecase {
	return r.someUsecase
}
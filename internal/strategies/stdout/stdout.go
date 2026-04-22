package stdout

import (
	"github.com/IcedElect/goverage/internal/structure/elements"
	"github.com/IcedElect/goverage/internal/structure/files"
	"github.com/IcedElect/goverage/internal/structure/tree"
)

type StdoutStrategy struct {
}

func NewStdoutStrategy() *StdoutStrategy {
	return &StdoutStrategy{}
}

func (s *StdoutStrategy) Name() string {
	return "stdout"
}

func (s *StdoutStrategy) Execute(
	directories []tree.Directory,
	filesRegistry *files.Registry,
	elementsRegistry *elements.Registry,
	threshold uint16,
	outputDir string,
) error {
	return nil
}

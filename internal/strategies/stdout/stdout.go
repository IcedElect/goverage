package stdout

import (
	"golang.org/x/tools/cover"
)

type StdoutStrategy struct {
}

func (s *StdoutStrategy) Name() string {
	return "Stdout"
}

func (s *StdoutStrategy) Execute(profiles []*cover.Profile, outputDir string) error {
	return nil
}

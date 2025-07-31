package profile

import (
	"fmt"

	"github.com/IcedElect/oh-my-cover-go/internal/strategies/html"
	"golang.org/x/tools/cover"
)

func ProcessProfile(profileFile string, outputDir string) (float64, error) {
	profiles, err := cover.ParseProfiles(profileFile)
	if err != nil {
		fmt.Printf("Error parsing cover profile: %v\n", err)
		return 0, nil
	}

	htmlStrategy := &html.HTMLStrategy{}
	coveragePercent, err := htmlStrategy.Execute(profiles, outputDir)
	if err != nil {
		return 0, err
	}
	return coveragePercent, nil
}

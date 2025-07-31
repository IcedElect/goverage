package profile

import (
	"fmt"

	"github.com/IcedElect/oh-my-cover-go/internal/strategies/html"
	"golang.org/x/tools/cover"
)

func ProcessProfile(profileFile string, outputDir string) error {
	profiles, err := cover.ParseProfiles(profileFile)
	if err != nil {
		fmt.Printf("Error parsing cover profile: %v\n", err)
		return nil
	}

	htmlStrategy := &html.HTMLStrategy{}
	err = htmlStrategy.Execute(profiles, outputDir)
	if err != nil {
		return err
	}
	return nil
}

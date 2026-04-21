package profile

import (
	"errors"
	"fmt"

	"github.com/IcedElect/goverage/internal/cli/ui"
	"github.com/IcedElect/goverage/internal/coverage"
	"github.com/IcedElect/goverage/internal/ignore"
	"github.com/IcedElect/goverage/internal/strategies/html"
	"github.com/IcedElect/goverage/internal/structure/elements"
	"github.com/IcedElect/goverage/internal/structure/files"
	"github.com/IcedElect/goverage/internal/structure/tree"
	"golang.org/x/tools/cover"
)

func ProcessProfile(profileFile string, outputDir string) (float64, error) {
	profiles, err := cover.ParseProfiles(profileFile)
	if err != nil {
		return 0, fmt.Errorf("error parsing cover profile: %w", err)
	}

	profiles = ignore.FilterProfiles(profiles)

	directories, err := tree.GetProfilesTree(profiles)
	if err != nil {
		return 0, fmt.Errorf("error building profiles tree: %w", err)
	}
	if len(directories) == 0 {
		return 0, errors.New("no profiles found")
	}

	filesRegistry, err := makeFilesRegistry(profiles)
	if err != nil {
		return 0, fmt.Errorf("error creating files registry: %w", err)
	}

	coverageCalculator := coverage.NewCalculator(filesRegistry)
	elementsRegistry := makeElementsRegistry(coverageCalculator, profiles, directories)

	totalCoverage := elementsRegistry.GetTotalCoverage()
	totalCoveragePercent := totalCoverage.Statements.Percent

	htmlStrategy := &html.HTMLStrategy{}
	err = htmlStrategy.Execute(directories, filesRegistry, elementsRegistry, outputDir)

	return totalCoveragePercent, err
}

func makeFilesRegistry(profiles []*cover.Profile) (*files.Registry, error) {
	filesRegistry, err := files.NewFilesRegistry(profiles)
	if err != nil {
		return nil, fmt.Errorf("error creating files registry: %w", err)
	}

	for _, profile := range profiles {
		if profileErr := filesRegistry.AddProfile(profile); profileErr != nil {
			ui.Errorlnf("Error adding profile [%s]: %v", profile.FileName, profileErr)
		}
	}

	return filesRegistry, nil
}

func makeElementsRegistry(
	coverageCalculator *coverage.Calculator,
	profiles []*cover.Profile,
	directories []tree.Directory,
) *elements.Registry {
	elementsRegistry := elements.NewElementsRegistry(coverageCalculator)

	for _, profile := range profiles {
		elementsRegistry.AddProfile(profile)
	}

	for _, dir := range directories {
		elementsRegistry.AddDirectory(dir, "")
	}

	return elementsRegistry
}

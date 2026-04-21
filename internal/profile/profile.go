package profile

import (
	"fmt"

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
		fmt.Printf("Error parsing cover profile: %v\n", err)
		return 0, nil
	}

	profiles = ignore.FilterProfiles(profiles)

	directories, err := tree.GetProfilesTree(profiles)
	if err != nil {
		return 0, fmt.Errorf("error building profiles tree: %v", err)
	}
	if len(directories) == 0 {
		return 0, fmt.Errorf("no profiles found")
	}

	filesRegistry, err := makeFilesRegistry(profiles)
	if err != nil {
		return 0, fmt.Errorf("error creating files registry: %v", err)
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
		return nil, fmt.Errorf("error creating files registry: %v", err)
	}

	// @TODO: use semaphore or workerpool for concurrent execution
	for _, profile := range profiles {
		if err := filesRegistry.AddProfile(profile); err != nil {
			return nil, fmt.Errorf("error adding profile [%s]: %v", profile.FileName, err)
		}
	}

	return filesRegistry, nil
}

func makeElementsRegistry(coverageCalculator *coverage.Calculator, profiles []*cover.Profile, directories []tree.Directory) *elements.Registry {
	elementsRegistry := elements.NewElementsRegistry(coverageCalculator)

	// @TODO: use semaphore or workerpool for concurrent execution
	for _, profile := range profiles {
		elementsRegistry.AddProfile(profile)
	}

	// @TODO: use semaphore or workerpool for concurrent execution
	for _, dir := range directories {
		elementsRegistry.AddDirectory(dir, "")
	}

	return elementsRegistry
}

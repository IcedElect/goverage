package html

import (
	"fmt"
	"os"
	"path"
	"path/filepath"
	"time"

	"github.com/IcedElect/goverage/internal/browser"
	"github.com/IcedElect/goverage/internal/utils"
	"golang.org/x/tools/cover"
)

type HTMLStrategy struct{}

func (s *HTMLStrategy) Name() string {
	return "HTML"
}

func (s *HTMLStrategy) Execute(profiles []*cover.Profile, outputDir string) (percent float64, err error) {
	outputPath := getOutputPath(outputDir)

	globalData = GlobalData{
		GeneratedTime: time.Now(),
	}

	coveragePercent, err := s.execute(profiles, outputPath)
	if err != nil {
		return 0, fmt.Errorf("error executing HTML strategy: %v", err)
	}

	if outputDir == "" && !browser.Open(path.Join("file://", outputPath, "index.html")) {
		fmt.Fprintf(os.Stderr, "HTML output written to %s\n", outputDir)
	}

	return coveragePercent, nil
}

func (s *HTMLStrategy) execute(profiles []*cover.Profile, outputDir string) (float64, error) {
	tree := utils.GetProfilesTree(profiles)
	if len(tree) == 0 {
		return 0, fmt.Errorf("no profiles found")
	}

	FilesRegistry, err := NewFilesRegistry(profiles)
	if err != nil {
		return 0, fmt.Errorf("error creating files registry: %v", err)
	}

	elementsRegistry := NewElementsRegistry(FilesRegistry)

	// @TODO: use semaphore or workerpool for concurrent execution
	for _, profile := range profiles {
		elementsRegistry.AddProfile(profile)
	}

	// @TODO: use semaphore or workerpool for concurrent execution
	for _, dir := range tree {
		elementsRegistry.AddDirectory(dir, "")
	}

	totalCoverage := elementsRegistry.GetTotalCoverage()
	totalCoveragePercent := totalCoverage.Statements.Percent
	globalData.TotalCoverage = totalCoverage

	err = s.executeAssets(filepath.Join(outputDir, "assets"))
	if err != nil {
		return totalCoveragePercent, fmt.Errorf("error executing assets: %v", err)
	}

	_, err = s.executeDirectory(outputDir, elementsRegistry, utils.Directory{
		Path:    "",
	})
	if err != nil {
		return totalCoveragePercent, fmt.Errorf("error executing root directory: %v", err)
	}

	err = s.executeDirectories(tree, elementsRegistry, outputDir)
	if err != nil {
		return totalCoveragePercent, fmt.Errorf("error executing directories: %v", err)
	}

	err = s.executeFiles(FilesRegistry.GetFiles(), outputDir)
	if err != nil {
		return totalCoveragePercent, fmt.Errorf("error executing files: %v", err)
	}

	return totalCoveragePercent, nil
}

func (s *HTMLStrategy) executeDirectories(dirs []utils.Directory, elementsRegistry *ElementsRegistry, outputDir string) error {
	for _, dir := range dirs {
		if _, err := s.executeDirectory(outputDir, elementsRegistry, dir); err != nil {
			return fmt.Errorf("error executing directory %s: %v", dir.Path, err)
		}
	}

	return nil
}

func (s *HTMLStrategy) executeFiles(files []*File, outputDir string) error {
	for _, file := range files {
		if _, err := s.executeFile(file, outputDir); err != nil {
			return fmt.Errorf("error executing file %s: %v", file.Name, err)
		}
	}

	return nil
}

func (s *HTMLStrategy) executeDirectory(outputDir string, elementsRegistry *ElementsRegistry, dir utils.Directory) (*os.File, error) {
	path := getPath(outputDir, dir.Path, "index.html")
	w, err := s.createFile(path)
	if err != nil {
		return w, fmt.Errorf("error creating index.html: %v", err)
	}
	defer w.Close()

	elements := elementsRegistry.GetElements(dir.Path)
	if len(elements) == 0 {
		return w, fmt.Errorf("no elements found for directory %s", dir.Path)
	}

	err = renderDirectory(w, dir, elements)
	if err != nil {
		return w, fmt.Errorf("error rendering directory: %v", err)
	}

	return w, nil
}

func (s *HTMLStrategy) executeFile(file *File, outputDir string) (*os.File, error) {
	path := getPath(outputDir, file.Path, file.Name+".html")
	w, err := s.createFile(path)
	if err != nil {
		return w, fmt.Errorf("error creating file %s: %v", path, err)
	}
	defer w.Close()

	err = renderFile(w, file)
	if err != nil {
		return w, fmt.Errorf("error rendering file: %v", err)
	}

	return w, nil
}

func (s *HTMLStrategy) executeAssets(outputPath string) error {
	files, err := assets.ReadDir("assets")
	if err != nil {
		return fmt.Errorf("error reading assets directory: %v", err)
	}
	for _, file := range files {
		src, err := assets.ReadFile("assets/" + file.Name())
		if err != nil {
			return fmt.Errorf("error reading asset file %s: %v", file.Name(), err)
		}
		// Create the destination directory if it doesn't exist
		if err := os.MkdirAll(outputPath, 0755); err != nil {
			return fmt.Errorf("error creating directory %s: %v", outputPath, err)
		}

		destPath := filepath.Join(outputPath, file.Name())
		if err := os.WriteFile(destPath, src, 0644); err != nil {
			return fmt.Errorf("error writing asset file %s: %v", file.Name(), err)
		}
	}

	return nil
}

func (s *HTMLStrategy) createFile(path string) (*os.File, error) {
	if err := os.MkdirAll(filepath.Dir(path), 0755); err != nil {
		return nil, fmt.Errorf("error creating directory %s: %v", filepath.Dir(path), err)
	}
	file, err := os.Create(path)
	if err != nil {
		return nil, fmt.Errorf("error creating file %s: %v", path, err)
	}
	return file, nil
}
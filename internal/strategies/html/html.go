package html

import (
	"fmt"
	"os"
	"path"
	"path/filepath"
	"time"

	"github.com/IcedElect/goverage/internal/cli/browser"
	"github.com/IcedElect/goverage/internal/cli/ui"
	"github.com/IcedElect/goverage/internal/coverage"
	"github.com/IcedElect/goverage/internal/structure/elements"
	"github.com/IcedElect/goverage/internal/structure/files"
	"github.com/IcedElect/goverage/internal/structure/tree"
	"github.com/IcedElect/goverage/internal/utils"
)

type HTMLStrategy struct{}

func (s *HTMLStrategy) Name() string {
	return "HTML"
}

func (s *HTMLStrategy) Execute(
	directories []tree.Directory,
	filesRegistry *files.Registry,
	elementsRegistry *elements.Registry,
	outputDir string,
) error {
	outputPath, err := utils.GetOutputPath(outputDir)
	if err != nil {
		return fmt.Errorf("error getting output path: %w", err)
	}

	globalData = GlobalData{
		GeneratedTime: time.Now(),
	}

	err = s.render(directories, filesRegistry, elementsRegistry, outputPath)
	if err != nil {
		return fmt.Errorf("error executing HTML strategy: %w", err)
	}

	if outputDir == "" && !browser.Open(path.Join("file://", outputPath, "index.html")) {
		ui.Infolnf("HTML output written to %s", outputPath)
	}

	return nil
}

func (s *HTMLStrategy) render(
	directories []tree.Directory,
	filesRegistry *files.Registry,
	elementsRegistry *elements.Registry,
	outputDir string,
) error {
	globalData.TotalCoverage = elementsRegistry.GetTotalCoverage()

	err := s.renderAssets(filepath.Join(outputDir, "assets"))
	if err != nil {
		return fmt.Errorf("error rendering assets: %w", err)
	}

	// Render root directory
	err = s.renderDirectory(outputDir, elementsRegistry, tree.Directory{Path: ""})
	if err != nil {
		return fmt.Errorf("error rendering root directory: %w", err)
	}

	s.renderDirectories(directories, elementsRegistry, outputDir)
	s.renderFiles(filesRegistry.GetFiles(), elementsRegistry, outputDir)

	return nil
}

func (s *HTMLStrategy) renderDirectories(
	dirs []tree.Directory,
	elementsRegistry *elements.Registry,
	outputDir string,
) {
	for _, dir := range dirs {
		if err := s.renderDirectory(outputDir, elementsRegistry, dir); err != nil {
			ui.Errorlnf("error rendering directory %s: %v", dir.Path, err)
		}
	}
}

func (s *HTMLStrategy) renderFiles(files []*files.File, elementsRegistry *elements.Registry, outputDir string) {
	for _, file := range files {
		if err := s.renderFile(file, elementsRegistry, outputDir); err != nil {
			ui.Errorlnf("error rendering file %s: %v", file.Name, err)
		}
	}
}

func (s *HTMLStrategy) renderDirectory(
	outputDir string,
	elementsRegistry *elements.Registry,
	dir tree.Directory,
) error {
	path := utils.GetPath(outputDir, dir.Path, "index.html")
	w, err := s.createFile(path)
	if err != nil {
		return fmt.Errorf("error creating index.html: %w", err)
	}
	defer w.Close()

	var coverage coverage.Coverage
	if directoryElement, ok := elementsRegistry.GetElement(dir.Path); ok {
		coverage = directoryElement.Coverage
	} else {
		coverage = elementsRegistry.GetTotalCoverage()
	}

	elements := elementsRegistry.GetElements(dir.Path)
	if len(elements) == 0 {
		return fmt.Errorf("no elements found for directory %s", dir.Path)
	}

	err = renderDirectory(w, dir, coverage, elements)
	if err != nil {
		return fmt.Errorf("error rendering directory: %w", err)
	}

	return nil
}

func (s *HTMLStrategy) renderFile(
	file *files.File,
	elementsRegistry *elements.Registry,
	outputDir string,
) error {
	path := utils.GetPath(outputDir, file.RelativePath, file.Name+".html")

	w, err := s.createFile(path)
	if err != nil {
		return fmt.Errorf("error creating file %s: %w", path, err)
	}
	defer w.Close()

	fileElement, ok := elementsRegistry.GetElement(file.Profile.FileName)
	if !ok {
		return fmt.Errorf("no element found for file %s", file.RelativePath)
	}

	err = renderFile(w, file, fileElement.Coverage)
	if err != nil {
		return fmt.Errorf("error rendering file: %w", err)
	}

	return nil
}

func (s *HTMLStrategy) renderAssets(outputPath string) error {
	files, err := assets.ReadDir("assets")
	if err != nil {
		return fmt.Errorf("error reading assets directory: %w", err)
	}
	for _, file := range files {
		src, readErr := assets.ReadFile("assets/" + file.Name())
		if readErr != nil {
			return fmt.Errorf("error reading asset file %s: %w", file.Name(), readErr)
		}
		// Create the destination directory if it doesn't exist
		if mkdirErr := os.MkdirAll(outputPath, 0750); mkdirErr != nil {
			return fmt.Errorf("error creating directory %s: %w", outputPath, mkdirErr)
		}

		destPath := filepath.Join(outputPath, file.Name())
		if writeErr := os.WriteFile(destPath, src, 0600); writeErr != nil {
			return fmt.Errorf("error writing asset file %s: %w", file.Name(), writeErr)
		}
	}

	return nil
}

func (s *HTMLStrategy) createFile(path string) (*os.File, error) {
	if err := os.MkdirAll(filepath.Dir(path), 0750); err != nil {
		return nil, fmt.Errorf("error creating directory [%s]: %w", filepath.Dir(path), err)
	}

	file, err := os.Create(path)
	if err != nil {
		return nil, fmt.Errorf("error creating file [%s]: %w", path, err)
	}

	return file, nil
}

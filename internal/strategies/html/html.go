package html

import (
	"fmt"
	"os"
	"path"
	"path/filepath"
	"time"

	"github.com/IcedElect/goverage/internal/browser"
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
	outputPath := utils.GetOutputPath(outputDir)

	globalData = GlobalData{
		GeneratedTime: time.Now(),
	}

	err := s.render(directories, filesRegistry, elementsRegistry, outputPath)
	if err != nil {
		return fmt.Errorf("error executing HTML strategy: %v", err)
	}

	if outputDir == "" && !browser.Open(path.Join("file://", outputPath, "index.html")) {
		fmt.Fprintf(os.Stderr, "HTML output written to %s\n", outputDir)
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
		return fmt.Errorf("error executing assets: %v", err)
	}

	// Render root directory
	err = s.renderDirectory(outputDir, elementsRegistry, tree.Directory{Path: ""})
	if err != nil {
		return fmt.Errorf("error executing root directory: %v", err)
	}

	// Render sub directories
	err = s.renderDirectories(directories, elementsRegistry, outputDir)
	if err != nil {
		return fmt.Errorf("error executing directories: %v", err)
	}

	// Render files
	err = s.renderFiles(filesRegistry.GetFiles(), outputDir)
	if err != nil {
		return fmt.Errorf("error executing files: %v", err)
	}

	return nil
}

func (s *HTMLStrategy) renderDirectories(dirs []tree.Directory, elementsRegistry *elements.Registry, outputDir string) error {
	for _, dir := range dirs {
		if err := s.renderDirectory(outputDir, elementsRegistry, dir); err != nil {
			return fmt.Errorf("error executing directory %s: %v", dir.Path, err)
		}
	}

	return nil
}

func (s *HTMLStrategy) renderFiles(files []*files.File, outputDir string) error {
	for _, file := range files {
		if _, err := s.renderFile(file, outputDir); err != nil {
			return fmt.Errorf("error executing file %s: %v", file.Name, err)
		}
	}

	return nil
}

func (s *HTMLStrategy) renderDirectory(outputDir string, elementsRegistry *elements.Registry, dir tree.Directory) error {
	path := utils.GetPath(outputDir, dir.Path, "index.html")
	w, err := s.createFile(path)
	if err != nil {
		return fmt.Errorf("error creating index.html: %v", err)
	}
	defer w.Close()

	elements := elementsRegistry.GetElements(dir.Path)
	if len(elements) == 0 {
		return fmt.Errorf("no elements found for directory %s", dir.Path)
	}

	err = renderDirectory(w, dir, elements)
	if err != nil {
		return fmt.Errorf("error rendering directory: %v", err)
	}

	return nil
}

func (s *HTMLStrategy) renderFile(file *files.File, outputDir string) (*os.File, error) {
	path := utils.GetPath(outputDir, file.RelativePath, file.Name+".html")

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

func (s *HTMLStrategy) renderAssets(outputPath string) error {
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
		return nil, fmt.Errorf("error creating directory [%s]: %v", filepath.Dir(path), err)
	}

	file, err := os.Create(path)
	if err != nil {
		return nil, fmt.Errorf("error creating file [%s]: %v", path, err)
	}

	return file, nil
}

package utils

import (
	"fmt"
	"os"
	"path/filepath"
)

func GetOutputPath(outputDir string) (string, error) {
	var outputPath string
	var err error

	if outputDir == "" {
		outputPath, err = os.MkdirTemp("", "cover")
		if err != nil {
			return "", fmt.Errorf("error creating temporary directory: %w", err)
		}
		return outputPath, nil
	}

	outputPath, err = filepath.Abs(outputDir)
	if err != nil {
		return "", fmt.Errorf("error getting absolute path for output directory: %w", err)
	}

	return outputPath, nil
}

func GetPath(basePath, path, fileName string) string {
	if path == "" {
		return filepath.Join(basePath, fileName)
	}
	return filepath.Join(basePath, path, fileName)
}

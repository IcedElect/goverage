package html

import (
	"fmt"
	"os"
	"path/filepath"
)

func getOutputPath(outputDir string) string {
	var outputPath string
	var err error

	if outputDir == "" {
		outputPath, err = os.MkdirTemp("", "cover")
		if err != nil {
			fmt.Printf("Error creating temporary directory: %v\n", err)
			return ""
		}
		return outputPath
	}

	outputPath, err = filepath.Abs(outputDir)
	if err != nil {
		fmt.Printf("Error getting absolute path for output directory: %v\n", err)
		return ""
	}

	return outputPath
}

func getPath(basePath, path, fileName string) string {
	if path == "" {
		return filepath.Join(basePath, fileName)
	}
	return filepath.Join(basePath, path, fileName)
}
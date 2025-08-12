package utils

import (
	"errors"
	"fmt"
	"log"
	"os"
	"path/filepath"

	"golang.org/x/mod/modfile"
)

var modulePath string

func GetModulePath() (string, error) {
	if modulePath != "" {
		return modulePath, nil
	}

	pwd, err := os.Getwd()
	if err != nil {
		return "", fmt.Errorf("failed to get current working directory: %w", err)
	}

	goModPath, err := findGoModPath(pwd)
	if err != nil {
		return "", fmt.Errorf("failed to find go.mod: %w", err)
	}

	fmt.Println(goModPath)
	
	data, err := os.ReadFile(goModPath)
	if err != nil {
		log.Fatalf("Failed to read go.mod: %v", err)
	}

	f, err := modfile.Parse("go.mod", data, nil)
	if err != nil {
		log.Fatalf("Failed to parse go.mod: %v", err)
	}

	modulePath = f.Module.Mod.Path
	return modulePath, nil
}

func findGoModPath(startDir string) (string, error) {
	dir := startDir

	for {
		goModPath := filepath.Join(dir, "go.mod")
		fmt.Println("Checking for go.mod at:", goModPath)
		if _, err := os.Stat(goModPath); err == nil {
			return goModPath, nil
		}

		parent := filepath.Dir(dir)
		if parent == dir {
			return "", errors.New("go.mod not found in any parent directory")
		}

		dir = parent
	}
}

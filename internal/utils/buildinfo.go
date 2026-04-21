package utils

import (
	"errors"
	"fmt"
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

	data, err := os.ReadFile(goModPath)
	if err != nil {
		return "", fmt.Errorf("failed to read go.mod: %w", err)
	}

	f, err := modfile.Parse("go.mod", data, nil)
	if err != nil {
		return "", fmt.Errorf("failed to parse go.mod: %w", err)
	}

	modulePath = f.Module.Mod.Path
	return modulePath, nil
}

func findGoModPath(startDir string) (string, error) {
	dir := startDir

	for {
		goModPath := filepath.Join(dir, "go.mod")
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

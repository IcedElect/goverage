package ignore

import (
	"bufio"
	"io/fs"
	"os"
	"path/filepath"
	"strings"

	"golang.org/x/tools/cover"
)

const (
	IgnoreFileName = ".coverignore"
)

func FilterProfiles(profiles []*cover.Profile) []*cover.Profile {
	lines := GetIgnoreLines(".")
	ignores, err := NewFromLines(lines)
	if err != nil {
		return nil
	}

	var filtered []*cover.Profile
	for _, profile := range profiles {
		if !ignores.Match(profile.FileName) {
			filtered = append(filtered, profile)
		}
	}

	return filtered
}

func GetIgnoreLines(root string) []string {
	var lines []string

	filepath.WalkDir(root, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		if d.IsDir() {
			return nil
		}

		if filepath.Base(path) != IgnoreFileName {
			return nil
		}

		file, err := os.Open(path)
		if err != nil {
			return err
		}
		defer file.Close()

		scanner := bufio.NewScanner(file)
		for scanner.Scan() {
			line := strings.TrimSpace(scanner.Text())
			if line == "" || strings.HasPrefix(line, "#") {
				continue
			}
			lines = append(lines, appendLinePath(line, filepath.Dir(path)))
		}

		return nil
	})

	return lines
}

func appendLinePath(pattern string, path string) string {
	if strings.HasPrefix(pattern, "!") {
		return "!" + filepath.Join(path, pattern[1:])
	}

	if strings.HasPrefix(pattern, "/") {
		return filepath.Join(path, pattern[1:])
	}

	return filepath.Join(path, pattern)
}
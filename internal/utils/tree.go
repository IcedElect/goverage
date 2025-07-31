package utils

import (
	"path/filepath"
	"strings"

	"golang.org/x/tools/cover"
)

type Directory struct {
	Path        string
	Profiles    []*cover.Profile
}

func GetProfilesTree(profiles []*cover.Profile) []*Directory {
	if len(profiles) == 0 {
		return nil
	}

	tree := make(map[string]*Directory)
	for _, profile := range profiles {
		fileName := strings.TrimPrefix(profile.FileName, GetModulePath())
		fileName = filepath.ToSlash(fileName)
		if strings.HasPrefix(profile.FileName, ".") || filepath.IsAbs(profile.FileName) {
			// Relative or absolute path.
			continue
		}
		dirPath := filepath.Dir(fileName)
		if _, ok := tree[dirPath]; !ok {
			tree[dirPath] = &Directory{
				Path:     dirPath,
				Profiles: []*cover.Profile{},
			}
		}
		tree[dirPath].Profiles = append(tree[dirPath].Profiles, profile)
	}

	var result []*Directory
	for _, dir := range tree {
		result = append(result, dir)
	}
	return result
}

package tree

import (
	"fmt"
	"path/filepath"
	"strings"

	"github.com/IcedElect/goverage/internal/utils"
	"golang.org/x/tools/cover"
)

type Directory struct {
	Path     string
	Profiles []*cover.Profile
}

func GetProfilesTree(profiles []*cover.Profile) ([]Directory, error) {
	if len(profiles) == 0 {
		return nil, nil
	}

	modulePath, err := utils.GetModulePath()
	if err != nil {
		return nil, fmt.Errorf("error to detect module path: %w", err)
	}

	tree := make(map[string]Directory)
	for _, profile := range profiles {
		fileName := strings.TrimPrefix(profile.FileName, modulePath)
		fileName = filepath.ToSlash(fileName)
		if strings.HasPrefix(profile.FileName, ".") || filepath.IsAbs(profile.FileName) {
			// it's relative or absolute path.
			continue
		}

		dirPath := filepath.Dir(fileName)

		// Add slash to prefix if not has
		if !strings.HasPrefix(dirPath, "/") {
			dirPath = "/" + dirPath
		}

		dir, ok := tree[dirPath]
		if !ok {
			dir = Directory{
				Path:     dirPath,
				Profiles: make([]*cover.Profile, 0, len(dir.Profiles)),
			}
		}

		dir.Profiles = append(dir.Profiles, profile)

		tree[dirPath] = dir
	}

	var result []Directory
	for _, dir := range tree {
		result = append(result, dir)
	}
	return result, nil
}

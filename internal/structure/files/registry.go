package files

import (
	"fmt"
	"path"
	"strings"

	"github.com/IcedElect/goverage/internal/utils"
	"golang.org/x/tools/cover"
)

type File struct {
	Path         string
	RelativePath string
	Name         string
	Funcs        []*utils.FuncExtent
	Profile      *cover.Profile
}

type Registry struct {
	files map[string]*File
	dirs  map[string]*utils.Pkg
}

func NewFilesRegistry(profiles []*cover.Profile) (*Registry, error) {
	dirs, err := utils.FindPkgs(profiles)
	if err != nil {
		return nil, fmt.Errorf("failed to find pkgs: %w", err)
	}

	registry := &Registry{
		files: make(map[string]*File),
		dirs:  dirs,
	}

	return registry, nil
}

func (r *Registry) GetFiles() []*File {
	var files []*File
	for _, file := range r.files {
		files = append(files, file)
	}
	return files
}

func (r *Registry) GetFile(fileName string) (*File, bool) {
	file, ok := r.files[fileName]
	return file, ok
}

func (r *Registry) AddProfile(profile *cover.Profile) error {
	if _, ok := r.files[profile.FileName]; ok {
		return nil
	}

	filePath, err := utils.FindFile(r.dirs, profile.FileName)
	if err != nil {
		return fmt.Errorf("error finding file %s: %v", profile.FileName, err)
	}

	funcs, err := utils.FindFuncs(filePath)
	if err != nil {
		return fmt.Errorf("error finding functions in %s: %v", profile.FileName, err)
	}

	modulePath, err := utils.GetModulePath()
	if err != nil {
		return fmt.Errorf("error getting module path: %v", err)
	}

	name := path.Base(filePath)
	relativePath := path.Dir(profile.FileName)
	relativePath = strings.TrimPrefix(relativePath, modulePath)
	relativePath = strings.TrimPrefix(relativePath, "/")

	r.files[profile.FileName] = &File{
		Path:         filePath,
		RelativePath: relativePath,
		Name:         name,
		Funcs:        funcs,
		Profile:      profile,
	}

	return nil
}

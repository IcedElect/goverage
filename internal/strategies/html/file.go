package html

import (
	"bufio"
	"fmt"
	"html/template"
	"io"
	"math"
	"os"
	"path"
	"strings"

	"github.com/IcedElect/goverage/internal/utils"
	"golang.org/x/tools/cover"
)

type File struct {
	Path  string
	Name  string
	Funcs []*utils.FuncExtent
	Code template.HTML
}

type FilesRegistry struct {
	files map[string]*File
	dirs  map[string]*utils.Pkg
}

func NewFilesRegistry(profiles []*cover.Profile) (*FilesRegistry, error) {
	dirs, err := utils.FindPkgs(profiles)
	if err != nil {
		return nil, err
	}

	registry := &FilesRegistry{
		files: make(map[string]*File),
		dirs: dirs,
	}

	// @TODO: use semaphore or workerpool for concurrent execution
	for _, profile := range profiles {
		registry.AddProfile(profile)
	}

	return registry, nil
}

func (r *FilesRegistry) GetFiles() []*File {
	var files []*File
	for _, file := range r.files {
		files = append(files, file)
	}
	return files
}

func (r *FilesRegistry) GetFile(fileName string) (*File, bool) {
	file, ok := r.files[fileName]
	return file, ok
}

func (r *FilesRegistry) AddProfile(profile *cover.Profile) error {
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

	code, err := r.getFileHtml(filePath, profile)
	if err != nil {
		return fmt.Errorf("error generating HTML for %s: %v", profile.FileName, err)
	}

	modulePath, err := utils.GetModulePath()
	if err != nil {
		return fmt.Errorf("error getting module path: %v", err)
	}

	name := path.Base(filePath)
	path := path.Dir(profile.FileName)
	path = strings.TrimPrefix(path, modulePath)
	path = strings.TrimPrefix(path, "/")

	r.files[profile.FileName] = &File{
		Path:  path,
		Name:  name,
		Funcs: funcs,
		Code:  template.HTML(code),
	}

	return nil
}

func (r *FilesRegistry) getFileHtml(path string, profile *cover.Profile) (string, error) {
	src, err := os.ReadFile(path)
	if err != nil {
		return "", fmt.Errorf("can't read %q: %v", path, err)
	}
		
	var buf strings.Builder
	boundaries := profile.Boundaries(src)
	err = htmlGen(&buf, src, boundaries)
	if err != nil {
		return "", fmt.Errorf("error generating HTML for %s: %v", path, err)
	}

	return buf.String(), nil
}

// htmlGen generates an HTML coverage report with the provided filename,
// source code, and tokens, and writes it to the given Writer.
func htmlGen(w io.Writer, src []byte, boundaries []cover.Boundary) error {
	dst := bufio.NewWriter(w)
	for i := range src {
		for len(boundaries) > 0 && boundaries[0].Offset == i {
			b := boundaries[0]
			if b.Start {
				n := 0
				if b.Count > 0 {
					n = int(math.Floor(b.Norm*9)) + 1
				}
				fmt.Fprintf(dst, `<span class="cov%v" title="%v">`, n, b.Count)
			} else {
				dst.WriteString("</span>")
			}
			boundaries = boundaries[1:]
		}

		switch b := src[i]; b {
		case '>':
			dst.WriteString("&gt;")
		case '<':
			dst.WriteString("&lt;")
		case '&':
			dst.WriteString("&amp;")
		case '\t':
			dst.WriteString("        ")
		default:
			dst.WriteByte(b)
		}
	}
	return dst.Flush()
}
package html

import (
	"bufio"
	"fmt"
	"html/template"
	"io"
	"math"
	"os"
	"strings"

	"github.com/IcedElect/goverage/internal/coverage"
	"github.com/IcedElect/goverage/internal/structure/elements"
	"github.com/IcedElect/goverage/internal/structure/files"
	"github.com/IcedElect/goverage/internal/structure/tree"
	"golang.org/x/tools/cover"
)

func renderDirectory(w io.Writer, dir tree.Directory, coverage coverage.Coverage, elements []*elements.Element) error {
	tmplParsed, err := template.New("layout").
		Funcs(templateFuncs).
		ParseFS(templates, "templates/layout.html", "templates/directory.page.html")
	if err != nil {
		return fmt.Errorf("error parsing templates: %v", err)
	}

	err = tmplParsed.Execute(w, TemplateData{
		CurrentPath: dir.Path,
		Global:      globalData,
		Directory:   &dir,
		Coverage:    coverage,
		Elements:    elements,
	})
	if err != nil {
		return fmt.Errorf("error executing template: %v", err)
	}

	return nil
}

func renderFile(w io.Writer, file *files.File, coverage coverage.Coverage) error {
	tmplParsed, err := template.New("layout").
		Funcs(templateFuncs).
		ParseFS(templates, "templates/layout.html", "templates/file.page.html")
	if err != nil {
		return fmt.Errorf("error parsing templates: %v", err)
	}

	code, err := getFileCode(file.Path, file.Profile)
	if err != nil {
		return fmt.Errorf("error generating code HTML for %s: %v", file.Path, err)
	}

	err = tmplParsed.Execute(w, TemplateData{
		CurrentPath: file.RelativePath,
		Global:      globalData,
		File:        file,
		FileCode:    template.HTML(code),
		Coverage:    coverage,
	})
	if err != nil {
		return fmt.Errorf("error executing template: %v", err)
	}

	return nil
}

func getFileCode(path string, profile *cover.Profile) (string, error) {
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

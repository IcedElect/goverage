package html

import (
	"fmt"
	"html/template"
	"io"

	"github.com/IcedElect/oh-my-cover-go/internal/utils"
)

func renderDirectory(w io.Writer, dir utils.Directory, elements []*Element) error {
	tmplParsed, err := template.New("layout").
		Funcs(templateFuncs).
		ParseFS(templates, "templates/layout.html", "templates/directory.page.html")
	if err != nil {
		return fmt.Errorf("error parsing templates: %v", err)
	}

	err = tmplParsed.Execute(w, TemplateData{
		CurrentPath: dir.Path,
		Global:   globalData,
		Directory: &dir,
		Elements:  elements,
	})
	if err != nil {
		return fmt.Errorf("error executing template: %v", err)
    }

	return nil
}

func renderFile(w io.Writer, file *File) error {
	tmplParsed, err := template.New("layout").
		Funcs(templateFuncs).
		ParseFS(templates, "templates/layout.html", "templates/file.page.html")
	if err != nil {
		return fmt.Errorf("error parsing templates: %v", err)
	}

	err = tmplParsed.Execute(w, TemplateData{
		CurrentPath: file.Path,
		Global: globalData,
		File:   file,
	})
	if err != nil {
		return fmt.Errorf("error executing template: %v", err)
	}

	return nil
}
package html

import (
	"embed"
	"html/template"

	"github.com/IcedElect/oh-my-cover-go/internal/utils"
)

var (
	//go:embed templates/*
	// templates contains the HTML templates for the HTML strategy.
	templates embed.FS

	//go:embed assets/*
	// assets contains the static assets (like CSS and JS) for the HTML strategy.
	assets embed.FS

	globalData GlobalData

	templateFuncs = template.FuncMap{
		"level": level,
	}
)

type GlobalData struct {
	BaseUrl       template.URL
	TotalCoverage Coverage
}

type TemplateData struct {
	Global    GlobalData
	File      *File
	Directory *utils.Directory
	Elements  []*Element
}

func level(percent float64) string {
	if percent < 40 {
		return "low"
	} else if percent < 80 {
		return "medium"
	}
	return "high"
}
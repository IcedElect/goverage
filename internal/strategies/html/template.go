package html

import (
	"embed"
	"html/template"
	"strings"
	"time"

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
		"baseurl": baseurl,
		"timeformat": timeformat,
	}
)

type GlobalData struct {
	GeneratedTime time.Time
	TotalCoverage Coverage
}

type TemplateData struct {
	CurrentPath string
	Global      GlobalData
	File        *File
	Directory   *utils.Directory
	Elements    []*Element
}

func level(percent float64) string {
	if percent < 40 {
		return "low"
	} else if percent < 80 {
		return "medium"
	}
	return "high"
}

func baseurl(path string) template.URL {
	// Clean and trim any trailing slashes
	path = strings.Trim(path, "/")
	if path == "" {
		return "."
	}

	// Count the number of path segments
	segments := strings.Split(path, "/")
	depth := len(segments)

	// Return "../" for each segment
	return template.URL(strings.Repeat("../", depth))
}

func timeformat(t time.Time) string {
	// Format the time in a human-readable format
	return t.Format("2006-01-02 15:04:05 MST -0700")
}

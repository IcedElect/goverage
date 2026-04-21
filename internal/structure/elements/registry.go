package elements

import (
	"path"
	"strings"

	"github.com/IcedElect/goverage/internal/coverage"
	"github.com/IcedElect/goverage/internal/structure/tree"
	"github.com/IcedElect/goverage/internal/utils"
	"golang.org/x/tools/cover"
)

type CoverageCalculator interface {
	CoverageByProfile(profile *cover.Profile) coverage.Coverage
	CoverageByDirectory(dir tree.Directory) coverage.Coverage
}

type Element struct {
	Name     string
	Path     string
	Url      string
	Coverage coverage.Coverage
}

type Registry struct {
	elements           map[string]*Element
	coverageCalculator CoverageCalculator
}

func NewElementsRegistry(coverageCalculator CoverageCalculator) *Registry {
	return &Registry{
		elements:           make(map[string]*Element),
		coverageCalculator: coverageCalculator,
	}
}

func (r *Registry) GetTotalCoverage() coverage.Coverage {
	elements := r.GetElements("")
	totalStatements := 0
	coveredStatements := 0
	totalLines := 0
	coveredLines := 0
	totalFuncs := 0
	coveredFuncs := 0

	for _, element := range elements {
		totalStatements += element.Coverage.Statements.Total
		coveredStatements += element.Coverage.Statements.Covered
		totalLines += element.Coverage.Lines.Total
		coveredLines += element.Coverage.Lines.Covered
		totalFuncs += element.Coverage.Functions.Total
		coveredFuncs += element.Coverage.Functions.Covered
	}

	return coverage.Coverage{
		Statements:   coverage.NewCoverageItem(totalStatements, coveredStatements),
		Functions:    coverage.NewCoverageItem(totalFuncs, coveredFuncs),
		Lines:        coverage.NewCoverageItem(totalLines, coveredLines),
		TotalPercent: utils.Percent(int64(coveredStatements+coveredLines+coveredFuncs), int64(totalStatements+totalLines+totalFuncs)),
	}
}

func (r *Registry) GetElement(path string) (*Element, bool) {
	element, ok := r.elements[path]
	return element, ok
}

func (r *Registry) GetElements(path string) []*Element {
	elements := make([]*Element, 0)

	for _, element := range r.elements {
		if element.Path == path {
			elements = append(elements, element)
		}
	}

	return elements
}

func (r *Registry) AddProfile(profile *cover.Profile) *Element {
	if element, exists := r.elements[profile.FileName]; exists {
		return element
	}

	modulePath, err := utils.GetModulePath()
	if err != nil {
		return nil
	}

	filePath := path.Dir(profile.FileName)
	filePath = strings.TrimPrefix(filePath, modulePath)

	fileName := path.Base(profile.FileName)

	element := &Element{
		Name:     fileName,
		Path:     filePath,
		Url:      fileName + ".html",
		Coverage: r.coverageCalculator.CoverageByProfile(profile),
	}
	r.elements[profile.FileName] = element

	return element
}

func (r *Registry) AddDirectory(dir tree.Directory, path string) *Element {
	modulePath, err := utils.GetModulePath()
	if err != nil {
		return nil
	}

	element := &Element{
		Name:     dir.Path,
		Path:     strings.TrimPrefix(path, modulePath),
		Url:      strings.TrimPrefix(dir.Path, "/") + "/index.html",
		Coverage: r.coverageCalculator.CoverageByDirectory(dir),
	}
	r.elements[dir.Path] = element

	return element
}

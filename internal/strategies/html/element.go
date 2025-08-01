package html

import (
	"path"
	"strings"

	"github.com/IcedElect/goverage/internal/utils"
	"golang.org/x/tools/cover"
)

type Coverage struct {
	// Coverage summary by statements.
	Statements CoverageItem

	// Coverage summary by lines.
	Lines CoverageItem

	// Coverage summary by functions.
	Functions CoverageItem

	// Coverage total summary
	TotalPercent float64
}

type CoverageItem struct {
	Total	 int
	Covered  int
	Percent  float64
}

func makeCoverageItem(total, covered int) CoverageItem {
	return CoverageItem{
		Total:  total,
		Covered: covered,
		Percent: utils.Percent(int64(covered), int64(total)),
	}
}

type Element struct {
	Name     string
	Path     string
	Url	     string
	Coverage Coverage
}

type ElementsRegistry struct {
	elements      map[string]*Element
	filesRegistry *FilesRegistry
}

func NewElementsRegistry(filesRegistry *FilesRegistry) *ElementsRegistry {
	return &ElementsRegistry{
		elements:      make(map[string]*Element),
		filesRegistry: filesRegistry,
	}
}

func (r *ElementsRegistry) GetTotalCoverage() Coverage {
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

	return Coverage{
		Statements: makeCoverageItem(totalStatements, coveredStatements),
		Functions:  makeCoverageItem(totalFuncs, coveredFuncs),
		Lines:      makeCoverageItem(totalLines, coveredLines),
		TotalPercent: utils.Percent(int64(coveredStatements+coveredLines+coveredFuncs), int64(totalStatements+totalLines+totalFuncs)),
	}
}

func (r *ElementsRegistry) GetElements(path string) []*Element {
	elements := make([]*Element, 0)

	for _, element := range r.elements {
		if element.Path == path {
			elements = append(elements, element)
		}
	}

	return elements
}

func (r *ElementsRegistry) AddProfile(profile *cover.Profile) *Element {
	if element, exists := r.elements[profile.FileName]; exists {
		return element
	}

	// name := path.Base(profile.FileName)
	path := path.Dir(profile.FileName)
	path = strings.TrimPrefix(path, utils.GetModulePath())

	file, exists := r.filesRegistry.GetFile(profile.FileName)
	if !exists {
		return nil
	}

	coverage := r.coverageByProfile(profile)
	element := &Element{
		Name:     file.Name,
		Path:     path,
		Url:      file.Name + ".html",
		Coverage: coverage,
	}
	r.elements[profile.FileName] = element

	return element
}

func (r *ElementsRegistry) AddDirectory(dir utils.Directory, path string) *Element {
	coverage := r.coverageByDirectory(dir)
	element := &Element{
		Name:     dir.Path,
		Path:     strings.TrimPrefix(path, utils.GetModulePath()),
		Url:      strings.TrimPrefix(dir.Path, "/") + "/index.html",
		Coverage: coverage,
	}
	r.elements[dir.Path] = element

	return element
}

func (r *ElementsRegistry) coverageByProfile(profile *cover.Profile) Coverage {
	totalStatements := 0
	coveredStatements := 0
	totalLines := 0
	coveredLines := 0
	totalFuncs := 0
	coveredFuncs := 0

	for _, block := range profile.Blocks {
		totalStatements += block.NumStmt
		if block.Count > 0 {
			coveredStatements += block.NumStmt
		}
	}

	for _, block := range profile.Blocks {
		if block.StartLine == block.EndLine {
			totalLines++
			if block.Count > 0 {
				coveredLines++
			}
		} else {
			totalLines += (block.EndLine - block.StartLine + 1)
			if block.Count > 0 {
				coveredLines += (block.EndLine - block.StartLine + 1)
			}
		}
	}

	file, exists := r.filesRegistry.GetFile(profile.FileName)
	if !exists {
		return Coverage{
			Statements: makeCoverageItem(totalStatements, coveredStatements),
			Lines:      makeCoverageItem(totalLines, coveredLines),
			Functions:  makeCoverageItem(totalFuncs, coveredFuncs),
			TotalPercent: utils.Percent(int64(coveredStatements+coveredLines+coveredFuncs), int64(totalStatements+totalLines+totalFuncs)),
		}
	}

	totalFuncs = len(file.Funcs)
	for _, f := range file.Funcs {
		num, _ := f.Coverage(profile)
		if num > 0 {
			coveredFuncs += 1
		}
	}

	return Coverage{
		Statements: makeCoverageItem(totalStatements, coveredStatements),
		Functions:  makeCoverageItem(totalFuncs, coveredFuncs),
		Lines:      makeCoverageItem(totalLines, coveredLines),
		TotalPercent: utils.Percent(int64(coveredStatements+coveredLines+coveredFuncs), int64(totalStatements+totalLines+totalFuncs)),
	}
}

func (r *ElementsRegistry) coverageByDirectory(dir utils.Directory) Coverage {
	totalStatements := 0
	coveredStatements := 0
	totalLines := 0
	coveredLines := 0
	totalFuncs := 0
	coveredFuncs := 0

	for _, profile := range dir.Profiles {
		element, ok := r.elements[profile.FileName]
		if !ok {
			element = r.AddProfile(profile)
		}

		totalStatements += element.Coverage.Statements.Total
		coveredStatements += element.Coverage.Statements.Covered
		totalLines += element.Coverage.Lines.Total
		coveredLines += element.Coverage.Lines.Covered
		totalFuncs += element.Coverage.Functions.Total
		coveredFuncs += element.Coverage.Functions.Covered
	}

	return Coverage{
		Statements: makeCoverageItem(totalStatements, coveredStatements),
		Functions:      makeCoverageItem(totalFuncs, coveredFuncs),
		Lines:      makeCoverageItem(totalLines, coveredLines),
		TotalPercent: utils.Percent(int64(coveredStatements+coveredLines+coveredFuncs), int64(totalStatements+totalLines+totalFuncs)),
	}
}

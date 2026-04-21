package coverage

//go:generate echo $PWD - $GOPACKAGE - $GOFILE
//go:generate mockgen -source=./calculator.go -destination ./mocks.go -package $GOPACKAGE

import (
	"github.com/IcedElect/goverage/internal/structure/files"
	"github.com/IcedElect/goverage/internal/structure/tree"
	"github.com/IcedElect/goverage/internal/utils"
	"golang.org/x/tools/cover"
)

type FilesRegistry interface {
	GetFile(fileName string) (*files.File, bool)
}

type Cache interface {
	Get(profileName string) (Coverage, bool)
	Set(profileName string, coverage Coverage)
}

type Calculator struct {
	filesRegistry FilesRegistry
	cache         Cache
}

func NewCalculator(filesRegistry FilesRegistry) *Calculator {
	return &Calculator{
		filesRegistry: filesRegistry,
		cache:         newCache(),
	}
}

func (c *Calculator) CoverageByProfile(profile *cover.Profile) Coverage {
	if coverage, ok := c.cache.Get(profile.FileName); ok {
		return coverage
	}

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
		totalLines += (block.EndLine - block.StartLine + 1)
		if block.Count > 0 {
			coveredLines += (block.EndLine - block.StartLine + 1)
		}
	}

	if file, exists := c.filesRegistry.GetFile(profile.FileName); exists {
		totalFuncs = len(file.Funcs)
		for _, f := range file.Funcs {
			num, _ := f.Coverage(profile)
			if num > 0 {
				coveredFuncs += 1
			}
		}
	}

	coverage := Coverage{
		Statements: NewCoverageItem(totalStatements, coveredStatements),
		Functions:  NewCoverageItem(totalFuncs, coveredFuncs),
		Lines:      NewCoverageItem(totalLines, coveredLines),
		TotalPercent: utils.Percent(
			int64(coveredStatements+coveredLines+coveredFuncs),
			int64(totalStatements+totalLines+totalFuncs),
		),
	}

	c.cache.Set(profile.FileName, coverage)

	return coverage
}

func (c *Calculator) CoverageByDirectory(dir tree.Directory) Coverage {
	totalStatements := 0
	coveredStatements := 0
	totalLines := 0
	coveredLines := 0
	totalFuncs := 0
	coveredFuncs := 0

	for _, profile := range dir.Profiles {
		profileCoverage := c.CoverageByProfile(profile)

		totalStatements += profileCoverage.Statements.Total
		coveredStatements += profileCoverage.Statements.Covered
		totalLines += profileCoverage.Lines.Total
		coveredLines += profileCoverage.Lines.Covered
		totalFuncs += profileCoverage.Functions.Total
		coveredFuncs += profileCoverage.Functions.Covered
	}

	return Coverage{
		Statements: NewCoverageItem(totalStatements, coveredStatements),
		Functions:  NewCoverageItem(totalFuncs, coveredFuncs),
		Lines:      NewCoverageItem(totalLines, coveredLines),
		TotalPercent: utils.Percent(
			int64(coveredStatements+coveredLines+coveredFuncs),
			int64(totalStatements+totalLines+totalFuncs),
		),
	}
}

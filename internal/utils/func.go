// This file implements the visitor that computes the (line, column)-(line-column) range for each function.

package utils

import (
	"go/ast"
	"go/parser"
	"go/token"

	"golang.org/x/tools/cover"
)

// FuncExtent describes a function's extent in the source by file and position.
type FuncExtent struct {
	Name      string
	StartLine int
	StartCol  int
	EndLine   int
	EndCol    int
}

// FindFuncs parses the file and returns a slice of FuncExtent descriptors.
func FindFuncs(name string) ([]*FuncExtent, error) {
	fset := token.NewFileSet()
	parsedFile, err := parser.ParseFile(fset, name, nil, 0)
	if err != nil {
		return nil, err
	}
	visitor := &FuncVisitor{
		FSet:    fset,
		Name:    name,
		AstFile: parsedFile,
	}
	ast.Walk(visitor, visitor.AstFile)
	return visitor.Funcs, nil
}

// FuncVisitor implements the visitor that builds the function position list for a file.
type FuncVisitor struct {
	FSet    *token.FileSet
	Name    string // Name of file.
	AstFile *ast.File
	Funcs   []*FuncExtent
}

// Visit implements the ast.Visitor interface.
func (v *FuncVisitor) Visit(node ast.Node) ast.Visitor {
	switch n := node.(type) {
	case *ast.FuncDecl:
		if n.Body == nil {
			// Do not count declarations of assembly functions.
			break
		}
		start := v.FSet.Position(n.Pos())
		end := v.FSet.Position(n.End())
		fe := &FuncExtent{
			Name:      n.Name.Name,
			StartLine: start.Line,
			StartCol:  start.Column,
			EndLine:   end.Line,
			EndCol:    end.Column,
		}
		v.Funcs = append(v.Funcs, fe)
	}
	return v
}

// Coverage returns the fraction of the statements in the function that were covered, as a numerator and denominator.
func (f *FuncExtent) Coverage(profile *cover.Profile) (num, den int64) {
	// We could avoid making this n^2 overall by doing a single scan and annotating the functions,
	// but the sizes of the data structures is never very large and the scan is almost instantaneous.
	var covered, total int64
	// The blocks are sorted, so we can stop counting as soon as we reach the end of the relevant block.
	for _, b := range profile.Blocks {
		if b.StartLine > f.EndLine || (b.StartLine == f.EndLine && b.StartCol >= f.EndCol) {
			// Past the end of the function.
			break
		}
		if b.EndLine < f.StartLine || (b.EndLine == f.StartLine && b.EndCol <= f.StartCol) {
			// Before the beginning of the function
			continue
		}
		total += int64(b.NumStmt)
		if b.Count > 0 {
			covered += int64(b.NumStmt)
		}
	}
	return covered, total
}
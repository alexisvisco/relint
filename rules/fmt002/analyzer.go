package fmt002

import (
	"go/ast"
	"go/token"

	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/analysis/passes/inspect"
	"golang.org/x/tools/go/ast/inspector"
)

var Analyzer = &analysis.Analyzer{
	Name:     "fmt002",
	Doc:      "FMT-002: top-level declarations must be in order: type, const, var, func",
	Run:      run,
	Requires: []*analysis.Analyzer{inspect.Analyzer},
}

func declOrder(d ast.Decl) int {
	gd, ok := d.(*ast.GenDecl)
	if !ok {
		return 3
	}
	switch gd.Tok {
	case token.TYPE:
		return 0
	case token.CONST:
		return 1
	case token.VAR:
		return 2
	}
	return 4
}

func run(pass *analysis.Pass) (interface{}, error) {
	insp := pass.ResultOf[inspect.Analyzer].(*inspector.Inspector)
	_ = insp

	for _, f := range pass.Files {
		checkFile(pass, f)
	}
	return nil, nil
}

func checkFile(pass *analysis.Pass, f *ast.File) {
	// Exclude import decls — they always precede other declarations in valid Go
	// and their declOrder (4) would cause false positives against type (0).
	var decls []ast.Decl
	for _, d := range f.Decls {
		gd, ok := d.(*ast.GenDecl)
		if ok && gd.Tok == token.IMPORT {
			continue
		}
		decls = append(decls, d)
	}

	for i := 1; i < len(decls); i++ {
		if declOrder(decls[i]) < declOrder(decls[i-1]) {
			pass.Report(analysis.Diagnostic{
				Pos:     decls[i].Pos(),
				Message: "FMT-002: declarations must be in order: type, const, var, func",
			})
			return
		}
	}
}
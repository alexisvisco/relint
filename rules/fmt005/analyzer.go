package fmt005

import (
	"go/ast"
	"go/token"

	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/analysis/passes/inspect"
	"golang.org/x/tools/go/ast/inspector"
)

var Analyzer = &analysis.Analyzer{
	Name:     "fmt005",
	Doc:      "FMT-005: type specs in type blocks must be separated by exactly one blank line",
	Run:      run,
	Requires: []*analysis.Analyzer{inspect.Analyzer},
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
	for _, decl := range f.Decls {
		gd, ok := decl.(*ast.GenDecl)
		if !ok || gd.Tok != token.TYPE || gd.Lparen == token.NoPos || len(gd.Specs) < 2 {
			continue
		}
		for i := 1; i < len(gd.Specs); i++ {
			prev := gd.Specs[i-1]
			curr := gd.Specs[i]
			if blankLinesBetween(pass, f, prev.End(), curr.Pos()) != 1 {
				pass.Reportf(curr.Pos(), "FMT-005: type specs in a type block must be separated by exactly one blank line")
			}
		}
	}
}

func blankLinesBetween(pass *analysis.Pass, file *ast.File, from, to token.Pos) int {
	fset := pass.Fset
	fromLine := fset.Position(from).Line
	toLine := fset.Position(to).Line

	if toLine <= fromLine+1 {
		return 0
	}

	commentLines := make(map[int]bool)
	for _, cg := range file.Comments {
		for _, c := range cg.List {
			line := fset.Position(c.Slash).Line
			if line > fromLine && line < toLine {
				commentLines[line] = true
			}
		}
	}

	blank := 0
	for line := fromLine + 1; line < toLine; line++ {
		if !commentLines[line] {
			blank++
		}
	}
	return blank
}

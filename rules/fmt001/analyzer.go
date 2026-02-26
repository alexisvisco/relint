package fmt001

import (
	"go/ast"
	"go/token"

	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/analysis/passes/inspect"
	"golang.org/x/tools/go/ast/inspector"
)

var Analyzer = &analysis.Analyzer{
	Name:     "fmt001",
	Doc:      "FMT-001: consecutive type declarations must be merged into a type block",
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
	decls := f.Decls
	i := 0
	for i < len(decls) {
		gd, ok := decls[i].(*ast.GenDecl)
		if !ok || gd.Tok != token.TYPE || gd.Lparen != token.NoPos {
			i++
			continue
		}

		start := i
		j := i + 1
		for j < len(decls) {
			next, ok := decls[j].(*ast.GenDecl)
			if !ok || next.Tok != token.TYPE || next.Lparen != token.NoPos {
				break
			}
			j++
		}

		if j-start < 2 {
			i = j
			continue
		}

		run := decls[start:j]
		pass.Report(analysis.Diagnostic{
			Pos:     run[1].Pos(),
			End:     run[len(run)-1].End(),
			Message: "FMT-001: consecutive type declarations should be merged into a type block",
		})
		i = j
	}
}
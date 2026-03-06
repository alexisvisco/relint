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
	Doc:      "FMT-001: consecutive type/const/var declarations must be merged into declaration blocks",
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
		if !ok || !isMergeableGenDecl(gd) {
			i++
			continue
		}

		start := i
		j := i + 1
		for j < len(decls) {
			next, ok := decls[j].(*ast.GenDecl)
			if !ok || !isMergeableGenDecl(next) || next.Tok != gd.Tok {
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
			Message: "FMT-001: consecutive declarations should be merged into a declaration block",
		})
		i = j
	}
}

func isMergeableGenDecl(gd *ast.GenDecl) bool {
	if gd == nil || gd.Lparen != token.NoPos {
		return false
	}
	return gd.Tok == token.TYPE || gd.Tok == token.CONST || gd.Tok == token.VAR
}

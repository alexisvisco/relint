package fmt004

import (
	"go/ast"
	"go/token"

	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/analysis/passes/inspect"
	"golang.org/x/tools/go/ast/inspector"
)

var Analyzer = &analysis.Analyzer{
	Name:     "fmt004",
	Doc:      "FMT-004: interface method signatures must be separated by exactly one blank line",
	Run:      run,
	Requires: []*analysis.Analyzer{inspect.Analyzer},
}

func run(pass *analysis.Pass) (interface{}, error) {
	insp := pass.ResultOf[inspect.Analyzer].(*inspector.Inspector)

	insp.Preorder([]ast.Node{(*ast.InterfaceType)(nil)}, func(n ast.Node) {
		iface := n.(*ast.InterfaceType)
		if iface.Methods == nil || len(iface.Methods.List) < 2 {
			return
		}

		fset := pass.Fset
		methods := iface.Methods.List

		for i := 1; i < len(methods); i++ {
			prev := methods[i-1]
			curr := methods[i]
			prevEnd := fset.Position(prev.End())
			currStart := fset.Position(curr.Pos())

			gap := currStart.Line - prevEnd.Line - 1

			if gap != 1 {
				if gap == 0 {
					pass.Reportf(curr.Pos(), "FMT-004: interface methods must be separated by exactly one blank line")
				} else {
					pass.Reportf(curr.Pos(), "FMT-004: interface methods must be separated by exactly one blank line")
				}
			}
		}
	})

	_ = token.NoPos
	return nil, nil
}

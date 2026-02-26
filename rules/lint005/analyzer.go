package lint005

import (
	"go/ast"

	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/analysis/passes/inspect"
	"golang.org/x/tools/go/ast/inspector"
)

var Analyzer = &analysis.Analyzer{
	Name:     "lint005",
	Doc:      "LINT-005: functions with more than 4 parameters should use a params struct",
	Run:      run,
	Requires: []*analysis.Analyzer{inspect.Analyzer},
}

func run(pass *analysis.Pass) (interface{}, error) {
	insp := pass.ResultOf[inspect.Analyzer].(*inspector.Inspector)

	nodeFilter := []ast.Node{(*ast.FuncDecl)(nil)}

	insp.Preorder(nodeFilter, func(n ast.Node) {
		fn := n.(*ast.FuncDecl)
		if fn.Type.Params == nil {
			return
		}

		count := 0
		for _, field := range fn.Type.Params.List {
			if len(field.Names) == 0 {
				count++
			} else {
				count += len(field.Names)
			}
		}

		if count > 4 {
			pass.Reportf(fn.Name.Pos(), "LINT-005: function %q has %d parameters, consider using a %sParams struct", fn.Name.Name, count, fn.Name.Name)
		}
	})

	return nil, nil
}

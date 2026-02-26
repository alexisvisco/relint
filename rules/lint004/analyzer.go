package lint004

import (
	"go/ast"
	"go/types"

	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/analysis/passes/inspect"
	"golang.org/x/tools/go/ast/inspector"
)

var Analyzer = &analysis.Analyzer{
	Name:     "lint004",
	Doc:      "LINT-004: context.Context must be the first parameter",
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

		params := fn.Type.Params.List
		if len(params) == 0 {
			return
		}

		// Flatten all params with their positions
		type param struct {
			typ ast.Expr
			idx int // absolute index in flattened list
		}
		var flat []param
		for _, field := range params {
			if len(field.Names) == 0 {
				flat = append(flat, param{field.Type, len(flat)})
			} else {
				for range field.Names {
					flat = append(flat, param{field.Type, len(flat)})
				}
			}
		}

		// Find if any param after index 0 is context.Context
		for i := 1; i < len(flat); i++ {
			if isContextType(pass, flat[i].typ) {
				pass.Reportf(fn.Name.Pos(), "LINT-004: context.Context must be the first parameter")
				return
			}
		}
	})

	return nil, nil
}

func isContextType(pass *analysis.Pass, expr ast.Expr) bool {
	t := pass.TypesInfo.TypeOf(expr)
	if t == nil {
		return false
	}
	named, ok := t.(*types.Named)
	if !ok {
		return false
	}
	obj := named.Obj()
	return obj.Pkg() != nil && obj.Pkg().Path() == "context" && obj.Name() == "Context"
}

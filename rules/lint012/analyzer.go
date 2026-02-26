package lint012

import (
	"go/ast"
	"go/types"
	"strings"

	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/analysis/passes/inspect"
	"golang.org/x/tools/go/ast/inspector"
)

var Analyzer = &analysis.Analyzer{
	Name:     "lint012",
	Doc:      "LINT-012: store functions must not return core/model types",
	Run:      run,
	Requires: []*analysis.Analyzer{inspect.Analyzer},
}

func run(pass *analysis.Pass) (interface{}, error) {
	pkgName := pass.Pkg.Name()
	if !strings.Contains(pkgName, "store") {
		return nil, nil
	}

	insp := pass.ResultOf[inspect.Analyzer].(*inspector.Inspector)

	insp.Preorder([]ast.Node{(*ast.FuncDecl)(nil)}, func(n ast.Node) {
		fn := n.(*ast.FuncDecl)
		// Must be a method on a *XStore receiver
		if fn.Recv == nil || len(fn.Recv.List) == 0 {
			return
		}
		recv := fn.Recv.List[0].Type
		// Dereference pointer
		if star, ok := recv.(*ast.StarExpr); ok {
			recv = star.X
		}
		recvIdent, ok := recv.(*ast.Ident)
		if !ok {
			return
		}
		if !strings.HasSuffix(recvIdent.Name, "Store") {
			return
		}

		if fn.Type.Results == nil {
			return
		}

		for _, field := range fn.Type.Results.List {
			t := pass.TypesInfo.TypeOf(field.Type)
			if containsCoreModel(t) {
				pass.Reportf(fn.Name.Pos(), "LINT-012: store function %q must not return core/model types", fn.Name.Name)
				return
			}
		}
	})

	return nil, nil
}

func containsCoreModel(t types.Type) bool {
	if t == nil {
		return false
	}
	// Unwrap pointer
	if ptr, ok := t.(*types.Pointer); ok {
		return containsCoreModel(ptr.Elem())
	}
	// Unwrap slice
	if sl, ok := t.(*types.Slice); ok {
		return containsCoreModel(sl.Elem())
	}
	named, ok := t.(*types.Named)
	if !ok {
		return false
	}
	obj := named.Obj()
	if obj.Pkg() != nil && strings.Contains(obj.Pkg().Path(), "core/model") {
		return true
	}
	return false
}

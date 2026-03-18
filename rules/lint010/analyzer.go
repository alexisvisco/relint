package lint010

import (
	"go/ast"
	"strings"

	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/analysis/passes/inspect"
	"golang.org/x/tools/go/ast/inspector"
)

var Analyzer = &analysis.Analyzer{
	Name:     "lint010",
	Doc:      "LINT-010: *Service and *Store interfaces must be declared in a types package, except under core",
	Run:      run,
	Requires: []*analysis.Analyzer{inspect.Analyzer},
}

func run(pass *analysis.Pass) (interface{}, error) {
	if pass.Pkg.Name() == "types" || isCorePackage(pass.Pkg.Path()) {
		return nil, nil
	}

	insp := pass.ResultOf[inspect.Analyzer].(*inspector.Inspector)

	insp.Preorder([]ast.Node{(*ast.GenDecl)(nil)}, func(n ast.Node) {
		gd := n.(*ast.GenDecl)
		for _, spec := range gd.Specs {
			ts, ok := spec.(*ast.TypeSpec)
			if !ok {
				continue
			}
			if _, isIface := ts.Type.(*ast.InterfaceType); !isIface {
				continue
			}
			if !strings.HasSuffix(ts.Name.Name, "Service") && !strings.HasSuffix(ts.Name.Name, "Store") {
				continue
			}
			pass.Reportf(ts.Name.Pos(), "LINT-010: interface %q must be declared in a types package", ts.Name.Name)
		}
	})

	return nil, nil
}

func isCorePackage(pkgPath string) bool {
	return strings.HasSuffix(pkgPath, "/core") || strings.Contains(pkgPath, "/core/")
}

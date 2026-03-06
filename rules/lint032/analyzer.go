package lint032

import (
	"go/ast"
	"strings"

	"golang.org/x/tools/go/analysis"
)

var Analyzer = &analysis.Analyzer{
	Name: "lint032",
	Doc:  "LINT-032: layer packages must expose a single constructor named New",
	Run:  run,
}

func run(pass *analysis.Pass) (interface{}, error) {
	pkgName := pass.Pkg.Name()
	if !isLayerPackage(pkgName) {
		return nil, nil
	}

	newFuncs := make([]*ast.FuncDecl, 0)
	for _, f := range pass.Files {
		for _, decl := range f.Decls {
			fn, ok := decl.(*ast.FuncDecl)
			if !ok || fn.Recv != nil || !fn.Name.IsExported() {
				continue
			}
			if !strings.HasPrefix(fn.Name.Name, "New") {
				continue
			}
			newFuncs = append(newFuncs, fn)
			if fn.Name.Name != "New" {
				pass.Reportf(fn.Name.Pos(), "LINT-032: constructor %q in package %q must be named \"New\"", fn.Name.Name, pkgName)
			}
		}
	}

	if len(newFuncs) > 1 {
		pass.Reportf(newFuncs[0].Name.Pos(), "LINT-032: package %q must expose only one constructor matching New*; found %d", pkgName, len(newFuncs))
	}

	return nil, nil
}

func isLayerPackage(name string) bool {
	if strings.HasSuffix(name, "store") || strings.HasSuffix(name, "service") {
		return true
	}
	return strings.HasSuffix(name, "handler") && name != "handler"
}

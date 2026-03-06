package lint019

import (
	"go/ast"
	"go/token"
	"path/filepath"
	"strings"

	"golang.org/x/tools/go/analysis"
)

var Analyzer = &analysis.Analyzer{
	Name: "lint019",
	Doc:  "LINT-019: FxModule must be declared in handler.go/service.go/store.go depending on package suffix",
	Run:  run,
}

func run(pass *analysis.Pass) (interface{}, error) {
	pkgName := pass.Pkg.Name()
	expectedFile := expectedFxModuleFile(pkgName)
	if expectedFile == "" {
		return nil, nil
	}

	for _, f := range pass.Files {
		base := filepath.Base(pass.Fset.File(f.Pos()).Name())

		for _, decl := range f.Decls {
			gd, ok := decl.(*ast.GenDecl)
			if !ok || gd.Tok != token.VAR {
				continue
			}
			for _, spec := range gd.Specs {
				vs, ok := spec.(*ast.ValueSpec)
				if !ok {
					continue
				}
				for _, n := range vs.Names {
					if n.Name != "FxModule" {
						continue
					}
					if base != expectedFile {
						pass.Reportf(n.Pos(), "LINT-019: FxModule in package %q must be declared in file %q", pkgName, expectedFile)
					}
				}
			}
		}
	}

	return nil, nil
}

func expectedFxModuleFile(pkgName string) string {
	switch {
	case strings.HasSuffix(pkgName, "handler"):
		return "handler.go"
	case strings.HasSuffix(pkgName, "service"):
		return "service.go"
	case strings.HasSuffix(pkgName, "store"):
		return "store.go"
	default:
		return ""
	}
}

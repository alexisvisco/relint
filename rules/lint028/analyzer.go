package lint028

import (
	"go/ast"
	"reflect"
	"strconv"

	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/analysis/passes/inspect"
	"golang.org/x/tools/go/ast/inspector"
)

var Analyzer = &analysis.Analyzer{
	Name:     "lint028",
	Doc:      "LINT-028: exported fields in model structs must declare a gorm tag",
	Run:      run,
	Requires: []*analysis.Analyzer{inspect.Analyzer},
}

func run(pass *analysis.Pass) (interface{}, error) {
	if pass.Pkg.Name() != "model" {
		return nil, nil
	}

	insp := pass.ResultOf[inspect.Analyzer].(*inspector.Inspector)
	insp.Preorder([]ast.Node{(*ast.TypeSpec)(nil)}, func(n ast.Node) {
		ts := n.(*ast.TypeSpec)
		st, ok := ts.Type.(*ast.StructType)
		if !ok {
			return
		}

		for _, field := range st.Fields.List {
			if len(field.Names) == 0 {
				continue
			}

			hasGorm := hasTagKey(field.Tag, "gorm")
			for _, name := range field.Names {
				if !name.IsExported() {
					continue
				}
				if hasGorm {
					continue
				}
				pass.Reportf(name.Pos(), "LINT-028: exported model field %q must declare a gorm tag", name.Name)
			}
		}
	})

	return nil, nil
}

func hasTagKey(tag *ast.BasicLit, key string) bool {
	if tag == nil {
		return false
	}
	if len(tag.Value) < 2 {
		return false
	}
	unquoted, err := strconv.Unquote(tag.Value)
	if err != nil {
		return false
	}
	_, ok := reflect.StructTag(unquoted).Lookup(key)
	return ok
}

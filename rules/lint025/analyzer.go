package lint025

import (
	"fmt"
	"go/ast"
	"go/token"
	"strings"
	"unicode"

	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/analysis/passes/inspect"
	"golang.org/x/tools/go/ast/inspector"

	"github.com/alexisvisco/relint/analysisutil"
)

var Analyzer = &analysis.Analyzer{
	Name:     "lint025",
	Doc:      "LINT-025: handler structs must be declared in {name}.go",
	Run:      run,
	Requires: []*analysis.Analyzer{inspect.Analyzer},
}

func run(pass *analysis.Pass) (interface{}, error) {
	if pass.Pkg.Name() != "handler" {
		return nil, nil
	}

	insp := pass.ResultOf[inspect.Analyzer].(*inspector.Inspector)
	insp.Preorder([]ast.Node{(*ast.GenDecl)(nil)}, func(n ast.Node) {
		gd := n.(*ast.GenDecl)
		if gd.Tok != token.TYPE {
			return
		}
		for _, spec := range gd.Specs {
			ts, ok := spec.(*ast.TypeSpec)
			if !ok {
				continue
			}
			if _, isStruct := ts.Type.(*ast.StructType); !isStruct {
				continue
			}
			name := ts.Name.Name
			if !strings.HasSuffix(name, "Handler") {
				continue
			}

			baseName := strings.TrimSuffix(name, "Handler")
			expectedFile := fmt.Sprintf("%s.go", toSnake(baseName))
			actualFile := analysisutil.FileBasename(pass, ts.Name.Pos())
			if actualFile != expectedFile {
				pass.Reportf(ts.Name.Pos(), "LINT-025: handler struct %q must be declared in file %q", name, expectedFile)
			}
		}
	})

	return nil, nil
}

func toSnake(s string) string {
	var result []rune
	for i, r := range s {
		if unicode.IsUpper(r) && i > 0 {
			result = append(result, '_')
		}
		result = append(result, unicode.ToLower(r))
	}
	return string(result)
}

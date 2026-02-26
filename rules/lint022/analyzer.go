package lint022

import (
	"fmt"
	"go/ast"
	"strings"
	"unicode"

	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/analysis/passes/inspect"
	"golang.org/x/tools/go/ast/inspector"

	"github.com/alexisvisco/relint/analysisutil"
)

var Analyzer = &analysis.Analyzer{
	Name:     "lint022",
	Doc:      "LINT-022: handler route methods must be in route handler files",
	Run:      run,
	Requires: []*analysis.Analyzer{inspect.Analyzer},
}

func run(pass *analysis.Pass) (interface{}, error) {
	if pass.Pkg.Name() != "handler" {
		return nil, nil
	}

	insp := pass.ResultOf[inspect.Analyzer].(*inspector.Inspector)

	insp.Preorder([]ast.Node{(*ast.FuncDecl)(nil)}, func(n ast.Node) {
		fn := n.(*ast.FuncDecl)
		if fn.Recv == nil || len(fn.Recv.List) == 0 {
			return
		}

		recv := fn.Recv.List[0].Type
		if star, ok := recv.(*ast.StarExpr); ok {
			recv = star.X
		}
		recvIdent, ok := recv.(*ast.Ident)
		if !ok {
			return
		}
		recvName := recvIdent.Name
		if !strings.HasSuffix(recvName, "Handler") {
			return
		}

		methodName := fn.Name.Name
		if !fn.Name.IsExported() {
			return
		}

		// Compute expected file name
		// handlerName without "Handler" suffix, lowercased
		baseName := strings.TrimSuffix(recvName, "Handler")
		expectedFile := expectedRouteFile(baseName, methodName)
		actualFile := analysisutil.FileBasename(pass, fn.Name.Pos())

		if actualFile != expectedFile {
			pass.Reportf(fn.Name.Pos(), "LINT-022: route handler %q on %q must be in file %q", methodName, recvName, expectedFile)
		}
	})

	return nil, nil
}

func expectedRouteFile(handlerName, routeName string) string {
	handlerSnake := toSnake(handlerName)
	routeSnake := toSnake(routeName)
	if handlerName == routeName {
		return fmt.Sprintf("%s_handler.go", handlerSnake)
	}
	return fmt.Sprintf("%s_%s_handler.go", handlerSnake, routeSnake)
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

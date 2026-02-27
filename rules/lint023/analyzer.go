package lint023

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
	Name:     "lint023",
	Doc:      "LINT-023: Route Input/Output types must be in their route handler file",
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
			name := ts.Name.Name
			// Check if name ends with Input or Output (route-specific types)
			if !strings.HasSuffix(name, "Input") && !strings.HasSuffix(name, "Output") {
				continue
			}

			// Determine the route name from type name
			// e.g. LoginInput -> Login, RegisterOutput -> Register
			var routeName string
			if strings.HasSuffix(name, "Input") {
				routeName = strings.TrimSuffix(name, "Input")
			} else {
				routeName = strings.TrimSuffix(name, "Output")
			}

			if routeName == "" {
				continue
			}

			// Find possible handler name patterns by looking at methods in the package
			// For now, we need to determine the struct name
			// LoginInput suggests it belongs to *XHandler.Login → x_login_handler.go
			// We look for any *XHandler struct that has a Login method
			expectedFiles := findExpectedFiles(pass, routeName)
			actualFile := analysisutil.FileBasename(pass, ts.Name.Pos())

			// If actual file matches none of the expected files, flag
			matched := false
			for _, ef := range expectedFiles {
				if actualFile == ef {
					matched = true
					break
				}
			}

			if !matched && len(expectedFiles) > 0 {
				pass.Reportf(ts.Name.Pos(), "LINT-023: type %q must be declared in route file %q", name, expectedFiles[0])
			}
		}
	})

	return nil, nil
}

func findExpectedFiles(pass *analysis.Pass, routeName string) []string {
	var result []string
	seen := map[string]bool{}
	for _, f := range pass.Files {
		for _, decl := range f.Decls {
			fn, ok := decl.(*ast.FuncDecl)
			if !ok || fn.Recv == nil || len(fn.Recv.List) == 0 {
				continue
			}
			if fn.Name.Name != routeName {
				continue
			}
			recv := fn.Recv.List[0].Type
			if star, ok := recv.(*ast.StarExpr); ok {
				recv = star.X
			}
			recvIdent, ok := recv.(*ast.Ident)
			if !ok {
				continue
			}
			handlerName := strings.TrimSuffix(recvIdent.Name, "Handler")
			routeFile := expectedRouteFile(handlerName, routeName)
			if !seen[routeFile] {
				result = append(result, routeFile)
				seen[routeFile] = true
			}
			// Also allow shared handler file (e.g. auth.go, tenant.go).
			sharedFile := fmt.Sprintf("%s.go", toSnake(handlerName))
			if !seen[sharedFile] {
				result = append(result, sharedFile)
				seen[sharedFile] = true
			}
		}
	}
	return result
}

func expectedRouteFile(handlerName, routeName string) string {
	handlerSnake := toSnake(handlerName)
	routeSnake := toSnake(routeName)
	routePart := normalizeRoutePart(handlerSnake, routeSnake)
	if routePart == "" {
		return fmt.Sprintf("%s_handler.go", handlerSnake)
	}
	return fmt.Sprintf("%s_%s_handler.go", handlerSnake, routePart)
}

func normalizeRoutePart(handlerSnake, routeSnake string) string {
	routePart := routeSnake
	aliases := []string{handlerSnake, pluralize(handlerSnake)}
	for _, alias := range aliases {
		routePart = strings.TrimPrefix(routePart, alias+"_")
		routePart = strings.TrimSuffix(routePart, "_"+alias)
	}
	routePart = strings.Trim(routePart, "_")
	if routePart == handlerSnake || routePart == pluralize(handlerSnake) {
		return ""
	}
	return routePart
}

func pluralize(s string) string {
	if strings.HasSuffix(s, "y") && len(s) > 1 {
		prev := s[len(s)-2]
		if !strings.ContainsRune("aeiou", rune(prev)) {
			return s[:len(s)-1] + "ies"
		}
	}
	if strings.HasSuffix(s, "s") || strings.HasSuffix(s, "x") || strings.HasSuffix(s, "z") ||
		strings.HasSuffix(s, "ch") || strings.HasSuffix(s, "sh") {
		return s + "es"
	}
	return s + "s"
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

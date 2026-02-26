package lint021

import (
	"go/ast"
	"go/types"
	"strings"

	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/analysis/passes/inspect"
	"golang.org/x/tools/go/ast/inspector"
)

var Analyzer = &analysis.Analyzer{
	Name:     "lint021",
	Doc:      "LINT-021: store functions must not return not-found sentinel errors directly",
	Run:      run,
	Requires: []*analysis.Analyzer{inspect.Analyzer},
}

func run(pass *analysis.Pass) (interface{}, error) {
	pkgName := pass.Pkg.Name()
	if !strings.Contains(pkgName, "store") {
		return nil, nil
	}

	insp := pass.ResultOf[inspect.Analyzer].(*inspector.Inspector)

	// Check return statements for known not-found sentinels.
	insp.Preorder([]ast.Node{(*ast.ReturnStmt)(nil)}, func(n ast.Node) {
		ret := n.(*ast.ReturnStmt)
		for _, result := range ret.Results {
			if forbidden, ok := forbiddenNotFoundSelector(pass, result); ok {
				pass.Reportf(result.Pos(), "LINT-021: store function must not return %s directly, wrap it in a domain error", forbidden)
			}
		}
	})

	return nil, nil
}

func forbiddenNotFoundSelector(pass *analysis.Pass, expr ast.Expr) (string, bool) {
	sel, ok := expr.(*ast.SelectorExpr)
	if !ok {
		return "", false
	}
	ident, ok := sel.X.(*ast.Ident)
	if !ok {
		return "", false
	}
	obj := pass.TypesInfo.Uses[ident]
	if obj == nil {
		return "", false
	}
	pkgName, ok := obj.(*types.PkgName)
	if !ok {
		return "", false
	}
	path := pkgName.Imported().Path()
	localName := pkgName.Name()
	switch sel.Sel.Name {
	case "ErrNoRows":
		if path == "database/sql" || path == "github.com/jackc/pgx/v5" || strings.HasSuffix(path, "pgx") {
			return pathToSelector(path, "ErrNoRows"), true
		}
	case "ErrRecordNotFound":
		if path == "gorm.io/gorm" || path == "github.com/jinzhu/gorm" || localName == "gorm" {
			if localName == "gorm" {
				return "gorm.ErrRecordNotFound", true
			}
			return pathToSelector(path, "ErrRecordNotFound"), true
		}
	}
	return "", false
}

func pathToSelector(path, member string) string {
	switch path {
	case "database/sql":
		return "sql." + member
	case "github.com/jackc/pgx/v5":
		return "pgx." + member
	case "gorm.io/gorm":
		return "gorm." + member
	case "github.com/jinzhu/gorm":
		return "gorm." + member
	default:
		parts := strings.Split(path, "/")
		return parts[len(parts)-1] + "." + member
	}
}

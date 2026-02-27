package lint029

import (
	"go/ast"
	"reflect"
	"strconv"
	"strings"

	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/analysis/passes/inspect"
	"golang.org/x/tools/go/ast/inspector"
)

var Analyzer = &analysis.Analyzer{
	Name:     "lint029",
	Doc:      "LINT-029: model relation fields must be pointer types or slices of pointers",
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
			gormTag, ok := lookupTag(field.Tag, "gorm")
			if !ok {
				continue
			}

			relationKey, isRelation := relationAttribute(gormTag)
			if !isRelation {
				continue
			}

			if isPointerType(field.Type) || isSliceOfPointers(field.Type) {
				continue
			}

			for _, name := range fieldNames(field) {
				pass.Reportf(name.Pos(), "LINT-029: relation field %q with gorm tag %q must be a pointer or a slice of pointers ([]*Type)", name.Name, relationKey)
			}
		}
	})

	return nil, nil
}

func lookupTag(tag *ast.BasicLit, key string) (string, bool) {
	if tag == nil || len(tag.Value) < 2 {
		return "", false
	}
	unquoted, err := strconv.Unquote(tag.Value)
	if err != nil {
		return "", false
	}
	value, ok := reflect.StructTag(unquoted).Lookup(key)
	return value, ok
}

func relationAttribute(gormTag string) (string, bool) {
	for _, segment := range strings.Split(gormTag, ";") {
		segment = strings.TrimSpace(segment)
		if segment == "" {
			continue
		}

		key := segment
		if idx := strings.IndexByte(segment, ':'); idx >= 0 {
			key = segment[:idx]
		}
		key = strings.TrimSpace(strings.ToLower(key))

		switch key {
		case "foreignkey":
			return "foreignKey", true
		case "many2many":
			return "many2many", true
		case "polymorphictype":
			return "polymorphicType", true
		}
	}

	return "", false
}

func isPointerType(expr ast.Expr) bool {
	_, ok := expr.(*ast.StarExpr)
	return ok
}

func isSliceOfPointers(expr ast.Expr) bool {
	slice, ok := expr.(*ast.ArrayType)
	if !ok || slice.Len != nil {
		return false
	}
	_, ok = slice.Elt.(*ast.StarExpr)
	return ok
}

func fieldNames(field *ast.Field) []*ast.Ident {
	if len(field.Names) > 0 {
		return field.Names
	}
	if ident, ok := field.Type.(*ast.Ident); ok {
		return []*ast.Ident{ident}
	}
	return nil
}

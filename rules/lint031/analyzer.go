package lint031

import (
	"go/ast"
	"go/token"
	"reflect"
	"regexp"
	"strconv"
	"strings"
	"unicode"

	"golang.org/x/tools/go/analysis"
)

const humaImportPath = "github.com/danielgtaylor/huma/v2"

var (
	pathParamPattern = regexp.MustCompile(`\{([^{}]+)\}`)
	routeFuncs       = map[string]bool{
		"Get":     true,
		"Post":    true,
		"Put":     true,
		"Patch":   true,
		"Delete":  true,
		"Head":    true,
		"Options": true,
	}
)

var Analyzer = &analysis.Analyzer{
	Name: "lint031",
	Doc:  "LINT-031: huma path params must be lowerCamelCase",
	Run:  run,
}

func run(pass *analysis.Pass) (interface{}, error) {
	for _, f := range pass.Files {
		checkPathTags(pass, f)
		checkHumaRoutePathParams(pass, f)
	}

	return nil, nil
}

func checkHumaRoutePathParams(pass *analysis.Pass, f *ast.File) {
	aliases := humaAliases(f)
	if len(aliases) == 0 {
		return
	}

	ast.Inspect(f, func(n ast.Node) bool {
		call, ok := n.(*ast.CallExpr)
		if !ok {
			return true
		}
		sel, ok := call.Fun.(*ast.SelectorExpr)
		if !ok {
			return true
		}
		pkgIdent, ok := sel.X.(*ast.Ident)
		if !ok || !aliases[pkgIdent.Name] || !routeFuncs[sel.Sel.Name] {
			return true
		}
		if len(call.Args) < 2 {
			return true
		}
		pathLit, ok := call.Args[1].(*ast.BasicLit)
		if !ok || pathLit.Kind != token.STRING {
			return true
		}
		path, err := strconv.Unquote(pathLit.Value)
		if err != nil {
			return true
		}

		matches := pathParamPattern.FindAllStringSubmatch(path, -1)
		for _, m := range matches {
			if len(m) < 2 {
				continue
			}
			param := m[1]
			if isLowerCamel(param) {
				continue
			}
			pass.Reportf(pathLit.Pos(), "LINT-031: huma path param %q must be lowerCamelCase", param)
		}
		return true
	})
}

func checkPathTags(pass *analysis.Pass, f *ast.File) {
	ast.Inspect(f, func(n ast.Node) bool {
		field, ok := n.(*ast.Field)
		if !ok || field.Tag == nil || field.Tag.Kind != token.STRING {
			return true
		}

		rawTag, err := strconv.Unquote(field.Tag.Value)
		if err != nil {
			return true
		}
		pathTag := reflect.StructTag(rawTag).Get("path")
		if pathTag == "" {
			return true
		}
		pathName := strings.Split(pathTag, ",")[0]
		if pathName == "" || pathName == "-" {
			return true
		}
		if isLowerCamel(pathName) {
			return true
		}
		pass.Reportf(field.Tag.Pos(), "LINT-031: path tag %q must be lowerCamelCase", pathName)
		return true
	})
}

func humaAliases(f *ast.File) map[string]bool {
	aliases := make(map[string]bool)
	for _, imp := range f.Imports {
		p, err := strconv.Unquote(imp.Path.Value)
		if err != nil || !isHumaImportPath(p) {
			continue
		}
		alias := "huma"
		if imp.Name != nil && imp.Name.Name != "" && imp.Name.Name != "_" && imp.Name.Name != "." {
			alias = imp.Name.Name
		}
		aliases[alias] = true
	}
	return aliases
}

func isHumaImportPath(path string) bool {
	if path == humaImportPath {
		return true
	}
	return strings.HasSuffix(path, "/"+humaImportPath)
}

func isLowerCamel(s string) bool {
	if s == "" {
		return false
	}
	for i, r := range s {
		if i == 0 {
			if !unicode.IsLower(r) {
				return false
			}
			continue
		}
		if unicode.IsLetter(r) || unicode.IsDigit(r) {
			continue
		}
		return false
	}
	return true
}

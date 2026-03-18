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

var pathParamPattern = regexp.MustCompile(`\{([^{}]+)\}`)

var Analyzer = &analysis.Analyzer{
	Name: "lint031",
	Doc:  "LINT-031: httpapi path params must be lowerCamelCase",
	Run:  run,
}

func run(pass *analysis.Pass) (interface{}, error) {
	for _, f := range pass.Files {
		checkPathTags(pass, f)
		checkHTTPAPIRoutePathParams(pass, f)
	}

	return nil, nil
}

func checkHTTPAPIRoutePathParams(pass *analysis.Pass, f *ast.File) {
	ast.Inspect(f, func(n ast.Node) bool {
		call, ok := n.(*ast.CallExpr)
		if !ok {
			return true
		}
		sel, ok := call.Fun.(*ast.SelectorExpr)
		if !ok || sel.Sel.Name != "WithPattern" {
			return true
		}
		if len(call.Args) == 0 {
			return true
		}
		patternLit, ok := call.Args[0].(*ast.BasicLit)
		if !ok || patternLit.Kind != token.STRING {
			return true
		}
		pattern, err := strconv.Unquote(patternLit.Value)
		if err != nil {
			return true
		}
		_, path, hasPath := strings.Cut(strings.TrimSpace(pattern), " ")
		if !hasPath || path == "" {
			path = strings.TrimSpace(pattern)
		}

		matches := pathParamPattern.FindAllStringSubmatch(path, -1)
		for _, m := range matches {
			if len(m) < 2 {
				continue
			}
			param := normalizePathParam(m[1])
			if isLowerCamel(param) {
				continue
			}
			pass.Reportf(patternLit.Pos(), "LINT-031: httpapi path param %q must be lowerCamelCase", m[1])
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

func normalizePathParam(param string) string {
	return strings.TrimSuffix(param, "...")
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

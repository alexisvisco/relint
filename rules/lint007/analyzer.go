package lint007

import (
	"go/ast"
	"go/token"
	"go/types"
	"strings"

	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/analysis/passes/inspect"
	"golang.org/x/tools/go/ast/inspector"
)

var exceptionsFlag string

var Analyzer = &analysis.Analyzer{
	Name:     "lint007",
	Doc:      "LINT-007: enum const values must be prefixed with the type name",
	Run:      run,
	Requires: []*analysis.Analyzer{inspect.Analyzer},
}

func init() {
	Analyzer.Flags.StringVar(
		&exceptionsFlag,
		"exceptions",
		"environment.Environment",
		"comma-separated list of package.Type exceptions for enum prefix checks",
	)
}

func parseExceptions(flag string) map[string]bool {
	m := make(map[string]bool)
	for _, s := range strings.Split(flag, ",") {
		s = strings.TrimSpace(s)
		if s == "" {
			continue
		}
		m[s] = true
	}
	return m
}

func run(pass *analysis.Pass) (interface{}, error) {
	insp := pass.ResultOf[inspect.Analyzer].(*inspector.Inspector)
	exceptions := parseExceptions(exceptionsFlag)
	pkgName := pass.Pkg.Name()

	// Collect named primitive types in this package
	primitiveTypes := map[string]bool{}

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
			// Only named primitive underlying types
			t := pass.TypesInfo.TypeOf(ts.Type)
			if t == nil {
				continue
			}
			underlying := t.Underlying()
			switch underlying.(type) {
			case *types.Basic:
				primitiveTypes[ts.Name.Name] = true
			}
		}
	})

	if len(primitiveTypes) == 0 {
		return nil, nil
	}

	// Check const declarations
	insp.Preorder([]ast.Node{(*ast.GenDecl)(nil)}, func(n ast.Node) {
		gd := n.(*ast.GenDecl)
		if gd.Tok != token.CONST {
			return
		}
		for _, spec := range gd.Specs {
			vs, ok := spec.(*ast.ValueSpec)
			if !ok {
				continue
			}
			if vs.Type == nil {
				continue
			}
			// Get the type name
			typeIdent, ok := vs.Type.(*ast.Ident)
			if !ok {
				continue
			}
			typeName := typeIdent.Name
			if !primitiveTypes[typeName] {
				continue
			}
			if exceptions[pkgName+"."+typeName] {
				continue
			}
			// Check each const name
			for _, name := range vs.Names {
				if !strings.HasPrefix(name.Name, typeName) {
					pass.Reportf(name.Pos(), "LINT-007: const %q must be prefixed with type name %q", name.Name, typeName)
				}
			}
		}
	})

	return nil, nil
}

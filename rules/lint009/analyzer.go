package lint009

import (
	"go/ast"
	"strings"

	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/analysis/passes/inspect"
	"golang.org/x/tools/go/ast/inspector"
)

var (
	pluralPackagesFlag string
	exceptionsFlag     string
)

var Analyzer = &analysis.Analyzer{
	Name:     "lint009",
	Doc:      "LINT-009: package name must not be plural",
	Run:      run,
	Requires: []*analysis.Analyzer{inspect.Analyzer},
}

func init() {
	Analyzer.Flags.StringVar(
		&pluralPackagesFlag,
		"plural-packages",
		"handlers,services,stores,types,models",
		"comma-separated list of package names considered plural and therefore forbidden",
	)
	Analyzer.Flags.StringVar(
		&exceptionsFlag,
		"exceptions",
		"types",
		"comma-separated list of package names exempt from the plural check",
	)
}

func parseSet(flag string) map[string]bool {
	m := make(map[string]bool)
	for _, s := range strings.Split(flag, ",") {
		s = strings.TrimSpace(s)
		if s != "" {
			m[s] = true
		}
	}
	return m
}

func run(pass *analysis.Pass) (interface{}, error) {
	plural := parseSet(pluralPackagesFlag)
	exceptions := parseSet(exceptionsFlag)

	pkgName := pass.Pkg.Name()
	if !plural[pkgName] || exceptions[pkgName] {
		return nil, nil
	}

	insp := pass.ResultOf[inspect.Analyzer].(*inspector.Inspector)
	reported := false
	insp.Preorder([]ast.Node{(*ast.File)(nil)}, func(n ast.Node) {
		if reported {
			return
		}
		f := n.(*ast.File)
		pass.Reportf(f.Name.Pos(), "LINT-009: package name %q must not be plural", pkgName)
		reported = true
	})

	return nil, nil
}

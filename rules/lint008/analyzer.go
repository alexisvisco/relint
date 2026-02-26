package lint008

import (
	"go/ast"
	"strings"

	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/analysis/passes/inspect"
	"golang.org/x/tools/go/ast/inspector"
)

var excludedSuffixesFlag string

var Analyzer = &analysis.Analyzer{
	Name:     "lint008",
	Doc:      "LINT-008: package name must not contain underscores (with configurable excluded suffixes)",
	Run:      run,
	Requires: []*analysis.Analyzer{inspect.Analyzer},
}

func init() {
	Analyzer.Flags.StringVar(
		&excludedSuffixesFlag,
		"excluded-suffixes",
		"_test",
		"comma-separated package-name suffixes excluded from underscore checks",
	)
}

func parseExcludedSuffixes(flag string) []string {
	suffixes := make([]string, 0)
	for _, s := range strings.Split(flag, ",") {
		s = strings.TrimSpace(s)
		if s == "" {
			continue
		}
		suffixes = append(suffixes, s)
	}
	return suffixes
}

func run(pass *analysis.Pass) (interface{}, error) {
	pkgName := pass.Pkg.Name()
	if !strings.Contains(pkgName, "_") {
		return nil, nil
	}
	for _, suffix := range parseExcludedSuffixes(excludedSuffixesFlag) {
		if strings.HasSuffix(pkgName, suffix) {
			return nil, nil
		}
	}

	// Report on the package clause of the first file
	insp := pass.ResultOf[inspect.Analyzer].(*inspector.Inspector)
	reported := false
	insp.Preorder([]ast.Node{(*ast.File)(nil)}, func(n ast.Node) {
		if reported {
			return
		}
		f := n.(*ast.File)
		pass.Reportf(f.Name.Pos(), "LINT-008: package name %q must not contain underscores", pkgName)
		reported = true
	})

	return nil, nil
}

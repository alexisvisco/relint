package lint019

import (
	"strings"

	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/analysis/passes/inspect"

	"github.com/alexisvisco/relint/analysisutil"
)

var Analyzer = &analysis.Analyzer{
	Name:     "lint019",
	Doc:      "LINT-019: *store/*service/*handler packages must contain fx_module.go with FxModule variable",
	Run:      run,
	Requires: []*analysis.Analyzer{inspect.Analyzer},
}

func run(pass *analysis.Pass) (interface{}, error) {
	pkgName := pass.Pkg.Name()
	if !strings.HasSuffix(pkgName, "store") &&
		!strings.HasSuffix(pkgName, "service") &&
		!strings.HasSuffix(pkgName, "handler") {
		return nil, nil
	}

	if !analysisutil.PackageContainsFile(pass, "fx_module.go") {
		// Report on the first file
		if len(pass.Files) > 0 {
			pass.Reportf(pass.Files[0].Pos(), "LINT-019: package %q must contain fx_module.go with FxModule variable", pkgName)
		}
		return nil, nil
	}

	// Check that fx_module.go has FxModule var
	for _, f := range pass.Files {
		if analysisutil.FileBasename(pass, f.Pos()) == "fx_module.go" {
			if !analysisutil.ContainsExportedVar(f, "FxModule") {
				pass.Reportf(f.Pos(), "LINT-019: fx_module.go in package %q must declare a FxModule variable", pkgName)
			}
			break
		}
	}

	return nil, nil
}

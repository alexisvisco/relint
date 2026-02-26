package lint024

import (
	"go/ast"
	"go/token"
	"regexp"
	"strings"

	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/analysis/passes/inspect"
	"golang.org/x/tools/go/ast/inspector"

	"github.com/alexisvisco/relint/analysisutil"
)

var Analyzer = &analysis.Analyzer{
	Name:     "lint024",
	Doc:      "LINT-024: body types in handler files must be named {Name}BodyInput or {Name}BodyOutput",
	Run:      run,
	Requires: []*analysis.Analyzer{inspect.Analyzer},
}

// validBodyPattern matches XBodyInput or XBodyOutput
var validBodyPattern = regexp.MustCompile(`^[A-Z][a-zA-Z0-9]*Body(Input|Output)$`)

func run(pass *analysis.Pass) (interface{}, error) {
	if pass.Pkg.Name() != "handler" {
		return nil, nil
	}

	bodyUsage := analysisutil.AnalyzeBodyTypeUsage(pass)
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
			if !strings.Contains(name, "Body") {
				continue
			}
			if bodyUsage.BodyOnlyStructs[name] {
				// Nested body-only helper structs are validated by LINT-026.
				continue
			}

			// Only check in non-route files (files not ending with _handler.go)
			filename := analysisutil.FileBasename(pass, ts.Name.Pos())
			if strings.HasSuffix(filename, "_handler.go") {
				continue
			}

			if !validBodyPattern.MatchString(name) {
				pass.Reportf(ts.Name.Pos(), "LINT-024: body type %q must be named \"{Name}BodyInput\" or \"{Name}BodyOutput\"", name)
			}
		}
	})

	return nil, nil
}

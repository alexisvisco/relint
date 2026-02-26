package fmt003

import (
	"go/ast"
	"go/token"

	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/analysis/passes/inspect"
	"golang.org/x/tools/go/ast/inspector"
)

var Analyzer = &analysis.Analyzer{
	Name:     "fmt003",
	Doc:      "FMT-003: functions must not start or end with blank lines, and must not have consecutive blank lines inside",
	Run:      run,
	Requires: []*analysis.Analyzer{inspect.Analyzer},
}

func run(pass *analysis.Pass) (interface{}, error) {
	insp := pass.ResultOf[inspect.Analyzer].(*inspector.Inspector)

	insp.Preorder([]ast.Node{(*ast.FuncDecl)(nil)}, func(n ast.Node) {
		fn := n.(*ast.FuncDecl)
		if fn.Body == nil {
			return
		}
		file := fileOf(pass, fn.Pos())
		if file == nil {
			return
		}
		if ast.IsGenerated(file) {
			return
		}
		checkFuncBody(pass, file, fn)
	})

	return nil, nil
}

func fileOf(pass *analysis.Pass, pos token.Pos) *ast.File {
	for _, f := range pass.Files {
		if f.Pos() <= pos && pos <= f.End() {
			return f
		}
	}
	return nil
}

// blankLinesBetween counts lines that contain neither code nor a comment
// in the half-open interval (from, to). Comments between statements are not
// blank lines and must not inflate the gap count.
func blankLinesBetween(pass *analysis.Pass, file *ast.File, from, to token.Pos) int {
	fset := pass.Fset
	fromLine := fset.Position(from).Line
	toLine := fset.Position(to).Line

	if toLine <= fromLine+1 {
		return 0
	}

	// Mark which lines carry a comment.
	commentLines := make(map[int]bool)
	for _, cg := range file.Comments {
		for _, c := range cg.List {
			line := fset.Position(c.Slash).Line
			if line > fromLine && line < toLine {
				commentLines[line] = true
			}
		}
	}

	blank := 0
	for line := fromLine + 1; line < toLine; line++ {
		if !commentLines[line] {
			blank++
		}
	}
	return blank
}

func checkFuncBody(pass *analysis.Pass, file *ast.File, fn *ast.FuncDecl) {
	body := fn.Body
	if len(body.List) == 0 {
		return
	}

	firstStmt := body.List[0]
	lastStmt := body.List[len(body.List)-1]

	// Leading blank line check (comments at the top of a func body are fine).
	if blankLinesBetween(pass, file, body.Lbrace, firstStmt.Pos()) > 0 {
		pass.Reportf(body.Lbrace, "FMT-003: function %q body must not start with a blank line", fn.Name.Name)
	}

	// Trailing blank line check.
	if blankLinesBetween(pass, file, lastStmt.End(), body.Rbrace) > 0 {
		pass.Reportf(body.Rbrace, "FMT-003: function %q body must not end with a blank line", fn.Name.Name)
	}

	// Consecutive blank lines inside the body (comments between statements are allowed).
	for i := 1; i < len(body.List); i++ {
		prev := body.List[i-1]
		curr := body.List[i]
		if blankLinesBetween(pass, file, prev.End(), curr.Pos()) > 1 {
			pass.Reportf(curr.Pos(), "FMT-003: function %q body must not have consecutive blank lines", fn.Name.Name)
		}
	}

	_ = token.NoPos
}

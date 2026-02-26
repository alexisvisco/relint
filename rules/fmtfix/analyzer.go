package fmtfix

import (
	"bytes"
	"go/ast"
	"go/format"
	"go/token"
	"sort"
	"strings"

	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/analysis/passes/inspect"
	"golang.org/x/tools/go/ast/inspector"
)

var Analyzer = &analysis.Analyzer{
	Name:     "fmtfix",
	Doc:      "FMTFIX: merge consecutive type declarations, normalize type-block spacing, and reorder top-level declarations (type, const, var, func)",
	Run:      run,
	Requires: []*analysis.Analyzer{inspect.Analyzer},
}

func run(pass *analysis.Pass) (interface{}, error) {
	insp := pass.ResultOf[inspect.Analyzer].(*inspector.Inspector)
	_ = insp

	for _, f := range pass.Files {
		checkFile(pass, f)
	}
	return nil, nil
}

// declItem is either a group of consecutive single type decls (to be merged)
// or a single other decl (non-type, or already-grouped type block).
type declItem struct {
	order     int
	typeDecls []*ast.GenDecl // non-nil: consecutive single type decls
	decl      ast.Decl       // non-nil: everything else
}

func declOrder(d ast.Decl) int {
	gd, ok := d.(*ast.GenDecl)
	if !ok {
		return 3
	}
	switch gd.Tok {
	case token.TYPE:
		return 0
	case token.CONST:
		return 1
	case token.VAR:
		return 2
	}
	return 4
}

// declStart returns the start position of a declaration including its doc comment.
// Using d.Pos() alone would leave the doc comment orphaned before the TextEdit range.
func declStart(d ast.Decl) token.Pos {
	switch v := d.(type) {
	case *ast.GenDecl:
		if v.Doc != nil {
			return v.Doc.Pos()
		}
	case *ast.FuncDecl:
		if v.Doc != nil {
			return v.Doc.Pos()
		}
	}
	return d.Pos()
}

// buildDeclItems groups consecutive single (non-parenthesised) type decls together;
// everything else is an individual item. Import decls must be excluded before calling.
func buildDeclItems(decls []ast.Decl) []declItem {
	var items []declItem
	i := 0
	for i < len(decls) {
		gd, ok := decls[i].(*ast.GenDecl)
		if ok && gd.Tok == token.TYPE && gd.Lparen == token.NoPos {
			j := i + 1
			for j < len(decls) {
				next, ok2 := decls[j].(*ast.GenDecl)
				if !ok2 || next.Tok != token.TYPE || next.Lparen != token.NoPos {
					break
				}
				j++
			}
			group := make([]*ast.GenDecl, j-i)
			for k := range group {
				group[k] = decls[i+k].(*ast.GenDecl)
			}
			items = append(items, declItem{order: 0, typeDecls: group})
			i = j
		} else {
			items = append(items, declItem{order: declOrder(decls[i]), decl: decls[i]})
			i++
		}
	}
	return items
}

// mergeAdjacentTypeGroups coalesces type groups that end up adjacent after sorting.
func mergeAdjacentTypeGroups(items []declItem) []declItem {
	merged := make([]declItem, 0, len(items))
	for _, item := range items {
		if item.typeDecls != nil && len(merged) > 0 && merged[len(merged)-1].typeDecls != nil {
			merged[len(merged)-1].typeDecls = append(merged[len(merged)-1].typeDecls, item.typeDecls...)
		} else {
			merged = append(merged, item)
		}
	}
	return merged
}

func checkFile(pass *analysis.Pass, f *ast.File) {
	// Exclude import decls — they are always first and must not be reordered.
	var decls []ast.Decl
	for _, d := range f.Decls {
		gd, ok := d.(*ast.GenDecl)
		if ok && gd.Tok == token.IMPORT {
			continue
		}
		decls = append(decls, d)
	}

	if len(decls) == 0 {
		return
	}

	items := buildDeclItems(decls)
	needsTypeSpacing := hasTypeBlockSpacingIssue(pass, f)

	needsMerge := false
	for _, item := range items {
		if len(item.typeDecls) > 1 {
			needsMerge = true
			break
		}
	}

	ordered := true
	for i := 1; i < len(items); i++ {
		if items[i].order < items[i-1].order {
			ordered = false
			break
		}
	}

	if !needsMerge && ordered && !needsTypeSpacing {
		return
	}

	sorted := make([]declItem, len(items))
	copy(sorted, items)
	sort.SliceStable(sorted, func(i, j int) bool {
		return sorted[i].order < sorted[j].order
	})

	// After reordering, type groups that are now adjacent should be merged too.
	sorted = mergeAdjacentTypeGroups(sorted)

	// Re-check: if nothing actually changed, skip.
	needsMerge = false
	for _, item := range sorted {
		if len(item.typeDecls) > 1 {
			needsMerge = true
			break
		}
	}
	if !needsMerge && ordered && !needsTypeSpacing {
		return
	}

	newSrc := generateText(pass, sorted)
	if newSrc == nil {
		return
	}

	// TextEdit.Pos must include the doc comment of the first decl so it is not
	// left orphaned before the edit range (which would misplace directives like //go:embed).
	startPos := declStart(decls[0])

	pass.Report(analysis.Diagnostic{
		Pos:     decls[0].Pos(),
		Message: "FMTFIX: apply format fixes (merge type blocks, reorder declarations)",
		SuggestedFixes: []analysis.SuggestedFix{{
			Message: "Apply format fixes",
			TextEdits: []analysis.TextEdit{{
				Pos:     startPos,
				End:     decls[len(decls)-1].End(),
				NewText: newSrc,
			}},
		}},
	})
}

func generateText(pass *analysis.Pass, items []declItem) []byte {
	var result bytes.Buffer
	for i, item := range items {
		if i > 0 {
			result.WriteString("\n\n")
		}
		if item.decl != nil {
			if gd, ok := item.decl.(*ast.GenDecl); ok && gd.Tok == token.TYPE && gd.Lparen != token.NoPos && len(gd.Specs) > 1 {
				block := buildTypeBlockWithSpacing(pass, gd)
				if block == nil {
					return nil
				}
				result.Write(block)
				continue
			}
			if err := format.Node(&result, pass.Fset, item.decl); err != nil {
				return nil
			}
		} else if len(item.typeDecls) == 1 {
			if err := format.Node(&result, pass.Fset, item.typeDecls[0]); err != nil {
				return nil
			}
		} else {
			merged := buildMergedTypeBlock(pass, item.typeDecls)
			if merged == nil {
				return nil
			}
			result.Write(merged)
		}
	}
	return result.Bytes()
}

// buildMergedTypeBlock produces a `type ( ... )` block from a slice of single-spec GenDecls.
// It preserves each spec's doc comment and uses format.Node for correct indentation.
func buildMergedTypeBlock(pass *analysis.Pass, decls []*ast.GenDecl) []byte {
	var buf bytes.Buffer
	buf.WriteString("type (\n")
	for i, d := range decls {
		// Include the GenDecl doc comment (e.g. // Foo is ...) indented by one tab.
		if d.Doc != nil {
			for _, c := range d.Doc.List {
				buf.WriteString("\t")
				buf.WriteString(c.Text)
				buf.WriteString("\n")
			}
		}
		// Format the TypeSpec and indent every line by one tab.
		ts := d.Specs[0].(*ast.TypeSpec)
		var specBuf bytes.Buffer
		if err := format.Node(&specBuf, pass.Fset, ts); err != nil {
			return nil
		}
		lines := strings.Split(specBuf.String(), "\n")
		for j, line := range lines {
			if j == len(lines)-1 && line == "" {
				continue // trailing newline from format.Node — skip it
			}
			buf.WriteString("\t")
			buf.WriteString(line)
			buf.WriteString("\n")
		}
		if i < len(decls)-1 {
			buf.WriteString("\n")
		}
	}
	buf.WriteString(")")
	return buf.Bytes()
}

// buildTypeBlockWithSpacing re-renders an existing type (...) block
// with exactly one blank line between each type spec.
func buildTypeBlockWithSpacing(pass *analysis.Pass, gd *ast.GenDecl) []byte {
	var buf bytes.Buffer
	buf.WriteString("type (\n")
	for i, spec := range gd.Specs {
		var specBuf bytes.Buffer
		if err := format.Node(&specBuf, pass.Fset, spec); err != nil {
			return nil
		}
		lines := strings.Split(specBuf.String(), "\n")
		for j, line := range lines {
			if j == len(lines)-1 && line == "" {
				continue
			}
			buf.WriteString("\t")
			buf.WriteString(line)
			buf.WriteString("\n")
		}
		if i < len(gd.Specs)-1 {
			buf.WriteString("\n")
		}
	}
	buf.WriteString(")")
	return buf.Bytes()
}

func hasTypeBlockSpacingIssue(pass *analysis.Pass, f *ast.File) bool {
	for _, decl := range f.Decls {
		gd, ok := decl.(*ast.GenDecl)
		if !ok || gd.Tok != token.TYPE || gd.Lparen == token.NoPos || len(gd.Specs) < 2 {
			continue
		}
		for i := 1; i < len(gd.Specs); i++ {
			if blankLinesBetween(pass, f, gd.Specs[i-1].End(), gd.Specs[i].Pos()) != 1 {
				return true
			}
		}
	}
	return false
}

// blankLinesBetween counts lines without code/comments between two positions.
func blankLinesBetween(pass *analysis.Pass, file *ast.File, from, to token.Pos) int {
	fset := pass.Fset
	fromLine := fset.Position(from).Line
	toLine := fset.Position(to).Line

	if toLine <= fromLine+1 {
		return 0
	}

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

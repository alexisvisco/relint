package lint027

import (
	"go/ast"
	"strconv"
	"strings"

	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/analysis/passes/inspect"
	"golang.org/x/tools/go/ast/inspector"
)

var Analyzer = &analysis.Analyzer{
	Name:     "lint027",
	Doc:      "LINT-027: structs in model packages must not declare json tags",
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
			if field.Tag == nil {
				continue
			}

			fixedTag, removeEntireTag, changed := removeJSONTag(field.Tag.Value)
			if !changed {
				continue
			}

			edit := analysis.TextEdit{
				Pos:     field.Tag.Pos(),
				End:     field.Tag.End(),
				NewText: []byte(fixedTag),
			}
			if removeEntireTag {
				edit.Pos = field.Type.End()
				edit.NewText = []byte("")
			}

			pass.Report(analysis.Diagnostic{
				Pos:     field.Tag.Pos(),
				Message: "LINT-027: model struct fields must not declare json tags",
				SuggestedFixes: []analysis.SuggestedFix{{
					Message:   "Remove json tag",
					TextEdits: []analysis.TextEdit{edit},
				}},
			})
		}
	})

	return nil, nil
}

func removeJSONTag(tagLiteral string) (newTagLiteral string, removeEntireTag bool, changed bool) {
	if len(tagLiteral) < 2 {
		return "", false, false
	}

	quote := tagLiteral[0]
	if quote != '`' && quote != '"' {
		return "", false, false
	}

	unquoted, err := strconv.Unquote(tagLiteral)
	if err != nil {
		return "", false, false
	}

	updated, removed := removeTagKey(unquoted, "json")
	if !removed {
		return "", false, false
	}

	if strings.TrimSpace(updated) == "" {
		return "", true, true
	}

	if quote == '`' && !strings.Contains(updated, "`") {
		return "`" + updated + "`", false, true
	}

	return strconv.Quote(updated), false, true
}

func removeTagKey(tag, keyToRemove string) (string, bool) {
	segments := make([]string, 0)
	removed := false

	for i := 0; i < len(tag); {
		for i < len(tag) && tag[i] == ' ' {
			i++
		}
		if i >= len(tag) {
			break
		}

		segStart := i
		keyStart := i
		for i < len(tag) && tag[i] > ' ' && tag[i] != ':' && tag[i] != '"' {
			i++
		}
		if keyStart == i || i >= len(tag) || tag[i] != ':' {
			return tag, false
		}

		key := tag[keyStart:i]
		i++ // skip ':'
		if i >= len(tag) || tag[i] != '"' {
			return tag, false
		}
		i++ // skip opening '"'

		closed := false
		for i < len(tag) {
			switch tag[i] {
			case '\\':
				i += 2
			case '"':
				i++
				closed = true
			default:
				i++
			}
			if closed {
				break
			}
		}
		if !closed {
			return tag, false
		}

		segment := tag[segStart:i]
		if key == keyToRemove {
			removed = true
			continue
		}
		segments = append(segments, segment)
	}

	if !removed {
		return tag, false
	}

	return strings.Join(segments, " "), true
}

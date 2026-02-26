package analysisutil

import (
	"go/ast"
	"go/token"
	"regexp"
	"strings"

	"golang.org/x/tools/go/analysis"
)

var bodyStructPattern = regexp.MustCompile(`^[A-Z][a-zA-Z0-9]*Body(Input|Output)$`)

type BodyStructMeta struct {
	Prefix string
	Suffix string
}

type BodyTypeUsage struct {
	BodyStructs      map[string]BodyStructMeta
	BodyOnlyStructs  map[string]bool
	UsedByBodyStruct map[string]map[string]bool
	DeclPos          map[string]token.Pos
}

// AnalyzeBodyTypeUsage computes relationships between body structs and local struct types.
// A "body-only struct" is a local struct type referenced by at least one body struct field
// and not referenced by non-body type declarations, function signatures, or typed vars.
func AnalyzeBodyTypeUsage(pass *analysis.Pass) *BodyTypeUsage {
	localTypes := make(map[string]bool)
	localStructs := make(map[string]bool)
	declPos := make(map[string]token.Pos)
	bodyStructs := make(map[string]BodyStructMeta)

	for _, f := range pass.Files {
		for _, decl := range f.Decls {
			gd, ok := decl.(*ast.GenDecl)
			if !ok || gd.Tok != token.TYPE {
				continue
			}
			for _, spec := range gd.Specs {
				ts, ok := spec.(*ast.TypeSpec)
				if !ok {
					continue
				}
				name := ts.Name.Name
				localTypes[name] = true
				declPos[name] = ts.Name.Pos()
				if _, isStruct := ts.Type.(*ast.StructType); isStruct {
					localStructs[name] = true
				}

				if suffix, ok := bodyStructSuffix(name); ok {
					bodyStructs[name] = BodyStructMeta{
						Prefix: strings.TrimSuffix(name, suffix),
						Suffix: suffix,
					}
				}
			}
		}
	}

	bodyUsage := make(map[string]int)
	nonBodyUsage := make(map[string]int)
	usedByBodyStruct := make(map[string]map[string]bool)

	for _, f := range pass.Files {
		for _, decl := range f.Decls {
			switch d := decl.(type) {
			case *ast.GenDecl:
				if d.Tok == token.TYPE {
					for _, spec := range d.Specs {
						ts, ok := spec.(*ast.TypeSpec)
						if !ok {
							continue
						}
						st, ok := ts.Type.(*ast.StructType)
						if !ok || st.Fields == nil {
							continue
						}
						_, isBodyStruct := bodyStructs[ts.Name.Name]
						for _, field := range st.Fields.List {
							for _, ref := range collectReferencedLocalTypes(field.Type, localTypes) {
								if isBodyStruct {
									bodyUsage[ref]++
									if usedByBodyStruct[ref] == nil {
										usedByBodyStruct[ref] = make(map[string]bool)
									}
									usedByBodyStruct[ref][ts.Name.Name] = true
								} else {
									nonBodyUsage[ref]++
								}
							}
						}
					}
				}

				if d.Tok == token.VAR {
					for _, spec := range d.Specs {
						vs, ok := spec.(*ast.ValueSpec)
						if !ok || vs.Type == nil {
							continue
						}
						for _, ref := range collectReferencedLocalTypes(vs.Type, localTypes) {
							nonBodyUsage[ref]++
						}
					}
				}

			case *ast.FuncDecl:
				if d.Recv != nil {
					for _, field := range d.Recv.List {
						for _, ref := range collectReferencedLocalTypes(field.Type, localTypes) {
							nonBodyUsage[ref]++
						}
					}
				}
				if d.Type.Params != nil {
					for _, field := range d.Type.Params.List {
						for _, ref := range collectReferencedLocalTypes(field.Type, localTypes) {
							nonBodyUsage[ref]++
						}
					}
				}
				if d.Type.Results != nil {
					for _, field := range d.Type.Results.List {
						for _, ref := range collectReferencedLocalTypes(field.Type, localTypes) {
							nonBodyUsage[ref]++
						}
					}
				}
			}
		}
	}

	bodyOnly := make(map[string]bool)
	for name := range localStructs {
		if bodyUsage[name] > 0 && nonBodyUsage[name] == 0 {
			bodyOnly[name] = true
		}
	}

	return &BodyTypeUsage{
		BodyStructs:      bodyStructs,
		BodyOnlyStructs:  bodyOnly,
		UsedByBodyStruct: usedByBodyStruct,
		DeclPos:          declPos,
	}
}

func bodyStructSuffix(name string) (string, bool) {
	if !bodyStructPattern.MatchString(name) {
		return "", false
	}
	if strings.HasSuffix(name, "Input") {
		return "Input", true
	}
	if strings.HasSuffix(name, "Output") {
		return "Output", true
	}
	return "", false
}

func collectReferencedLocalTypes(expr ast.Expr, localTypes map[string]bool) []string {
	refs := make(map[string]bool)
	var visit func(ast.Expr)
	visit = func(e ast.Expr) {
		switch t := e.(type) {
		case *ast.Ident:
			if localTypes[t.Name] {
				refs[t.Name] = true
			}
		case *ast.StarExpr:
			visit(t.X)
		case *ast.ArrayType:
			visit(t.Elt)
		case *ast.Ellipsis:
			visit(t.Elt)
		case *ast.MapType:
			visit(t.Key)
			visit(t.Value)
		case *ast.ChanType:
			visit(t.Value)
		case *ast.ParenExpr:
			visit(t.X)
		case *ast.IndexExpr:
			visit(t.X)
			visit(t.Index)
		case *ast.IndexListExpr:
			visit(t.X)
			for _, idx := range t.Indices {
				visit(idx)
			}
		case *ast.StructType:
			if t.Fields == nil {
				return
			}
			for _, f := range t.Fields.List {
				visit(f.Type)
			}
		case *ast.InterfaceType:
			if t.Methods == nil {
				return
			}
			for _, m := range t.Methods.List {
				visit(m.Type)
			}
		}
	}
	visit(expr)

	result := make([]string, 0, len(refs))
	for r := range refs {
		result = append(result, r)
	}
	return result
}

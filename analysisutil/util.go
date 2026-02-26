package analysisutil

import (
	"go/ast"
	"go/token"
	"go/types"
	"path/filepath"
	"strings"

	"golang.org/x/tools/go/analysis"
)

// IsSlogCall reports whether call is a call to slog.X or (*slog.Logger).X.
func IsSlogCall(pass *analysis.Pass, call *ast.CallExpr) bool {
	sel, ok := call.Fun.(*ast.SelectorExpr)
	if !ok {
		return false
	}
	// slog.Info, slog.Error, etc. (package-level)
	if ident, ok := sel.X.(*ast.Ident); ok {
		obj := pass.TypesInfo.Uses[ident]
		if obj != nil {
			if pkgName, ok := obj.(*types.PkgName); ok {
				return pkgName.Imported().Path() == "log/slog"
			}
		}
	}
	// (*slog.Logger).Info, etc. (method call)
	t := pass.TypesInfo.TypeOf(sel.X)
	if t != nil {
		// deref pointer
		if ptr, ok := t.(*types.Pointer); ok {
			t = ptr.Elem()
		}
		if named, ok := t.(*types.Named); ok {
			obj := named.Obj()
			if obj.Pkg() != nil && obj.Pkg().Path() == "log/slog" && obj.Name() == "Logger" {
				return true
			}
		}
	}
	return false
}

// IsInPackage reports whether pass.Pkg.Name() equals pkgName.
func IsInPackage(pass *analysis.Pass, pkgName string) bool {
	return pass.Pkg.Name() == pkgName
}

// FileBasename returns the base filename (without directory) for a given Pos.
func FileBasename(pass *analysis.Pass, pos token.Pos) string {
	return filepath.Base(pass.Fset.File(pos).Name())
}

// FilePath returns the full file path for a given Pos.
func FilePath(pass *analysis.Pass, pos token.Pos) string {
	return pass.Fset.File(pos).Name()
}

// IsSnakeLower reports whether s is a valid lowercase_snake_case identifier
// (including dot-notation like error.message).
func IsSnakeLower(s string) bool {
	if s == "" {
		return false
	}
	// Allow dot-notation segments
	parts := strings.Split(s, ".")
	for _, p := range parts {
		if !isSnakeSegment(p) {
			return false
		}
	}
	return true
}

func isSnakeSegment(s string) bool {
	if s == "" {
		return false
	}
	// Must start with lowercase letter
	if s[0] < 'a' || s[0] > 'z' {
		return false
	}
	for _, c := range s[1:] {
		if !((c >= 'a' && c <= 'z') || (c >= '0' && c <= '9') || c == '_') {
			return false
		}
	}
	return true
}

// PackageContainsFile reports whether any file in pass has the given base name.
func PackageContainsFile(pass *analysis.Pass, filename string) bool {
	for _, f := range pass.Files {
		if filepath.Base(pass.Fset.File(f.Pos()).Name()) == filename {
			return true
		}
	}
	return false
}

// ExportedFuncCount returns the count of exported top-level functions in a file.
func ExportedFuncCount(file *ast.File) int {
	count := 0
	for _, decl := range file.Decls {
		fn, ok := decl.(*ast.FuncDecl)
		if !ok {
			continue
		}
		if fn.Name.IsExported() {
			count++
		}
	}
	return count
}

// ContainsExportedVar reports whether a file has a top-level var declaration
// with the given name.
func ContainsExportedVar(file *ast.File, name string) bool {
	for _, decl := range file.Decls {
		gd, ok := decl.(*ast.GenDecl)
		if !ok {
			continue
		}
		if gd.Tok != token.VAR {
			continue
		}
		for _, spec := range gd.Specs {
			vs, ok := spec.(*ast.ValueSpec)
			if !ok {
				continue
			}
			for _, n := range vs.Names {
				if n.Name == name {
					return true
				}
			}
		}
	}
	return false
}

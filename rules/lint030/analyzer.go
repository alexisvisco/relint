package lint030

import (
	"go/ast"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"golang.org/x/mod/modfile"
	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/analysis/passes/inspect"
	"golang.org/x/tools/go/ast/inspector"
)

var (
	rootsFlag string
)

var Analyzer = &analysis.Analyzer{
	Name:     "lint030",
	Doc:      "LINT-030: packages under protected roots must not import sibling module roots",
	Run:      run,
	Requires: []*analysis.Analyzer{inspect.Analyzer},
}

func init() {
	Analyzer.Flags.StringVar(
		&rootsFlag,
		"roots",
		"core",
		"comma-separated list of protected root directories (for example: core,shared)",
	)
}

func run(pass *analysis.Pass) (interface{}, error) {
	protectedRoots := parseSet(rootsFlag)
	if len(protectedRoots) == 0 {
		return nil, nil
	}

	modulePath := resolveModulePath(pass)
	currentRoot, allowBareLocalImports := packageRoot(pass.Pkg.Path(), modulePath)
	if currentRoot == "" || !protectedRoots[currentRoot] {
		return nil, nil
	}

	insp := pass.ResultOf[inspect.Analyzer].(*inspector.Inspector)
	insp.Preorder([]ast.Node{(*ast.ImportSpec)(nil)}, func(n ast.Node) {
		spec := n.(*ast.ImportSpec)
		importPath, err := strconv.Unquote(spec.Path.Value)
		if err != nil || importPath == "" {
			return
		}

		importedRoot, isLocal := importedRootForPath(importPath, modulePath, allowBareLocalImports)
		if !isLocal || importedRoot == "" || importedRoot == currentRoot {
			return
		}

		pass.Reportf(
			spec.Path.Pos(),
			"LINT-030: package under root %q must not import sibling root %q via %q",
			currentRoot,
			importedRoot,
			importPath,
		)
	})

	return nil, nil
}

func parseSet(v string) map[string]bool {
	out := map[string]bool{}
	for _, s := range strings.Split(v, ",") {
		s = strings.TrimSpace(s)
		if s != "" {
			out[s] = true
		}
	}
	return out
}

func packageRoot(pkgPath, modulePath string) (root string, allowBareLocalImports bool) {
	relPath := pkgPath
	if modulePath != "" {
		if pkgPath == modulePath {
			return "", false
		}
		prefix := modulePath + "/"
		if strings.HasPrefix(pkgPath, prefix) {
			relPath = strings.TrimPrefix(pkgPath, prefix)
		}
	}

	normalized, srcStyle := stripSrcPrefix(relPath)
	return firstSegment(normalized), srcStyle
}

func importedRootForPath(importPath, modulePath string, allowBareLocalImports bool) (string, bool) {
	if modulePath != "" {
		prefix := modulePath + "/"
		if strings.HasPrefix(importPath, prefix) {
			rel := strings.TrimPrefix(importPath, prefix)
			normalized, _ := stripSrcPrefix(rel)
			return firstSegment(normalized), true
		}
		if !allowBareLocalImports {
			return "", false
		}
	}

	// Fallback for GOPATH-style testdata (analysistest):
	// treat slash-containing, dotless-root imports as local module imports.
	root := firstSegment(importPath)
	if root == "" || !strings.Contains(importPath, "/") || strings.Contains(root, ".") {
		return "", false
	}
	return root, true
}

func stripSrcPrefix(path string) (string, bool) {
	if i := strings.LastIndex(path, "/src/"); i >= 0 {
		return path[i+len("/src/"):], true
	}
	if strings.HasPrefix(path, "src/") {
		return strings.TrimPrefix(path, "src/"), true
	}
	return path, false
}

func firstSegment(s string) string {
	s = strings.TrimPrefix(s, "/")
	if s == "" {
		return ""
	}
	if i := strings.IndexByte(s, '/'); i >= 0 {
		return s[:i]
	}
	return s
}

func resolveModulePath(pass *analysis.Pass) string {
	if len(pass.Files) == 0 {
		return ""
	}

	file := pass.Fset.File(pass.Files[0].Pos())
	if file == nil {
		return ""
	}

	goModPath := findGoMod(filepath.Dir(file.Name()))
	if goModPath == "" {
		return ""
	}

	data, err := os.ReadFile(goModPath)
	if err != nil {
		return ""
	}

	parsed, err := modfile.Parse(goModPath, data, nil)
	if err != nil || parsed == nil || parsed.Module == nil {
		return ""
	}

	return strings.TrimSpace(parsed.Module.Mod.Path)
}

func findGoMod(startDir string) string {
	dir := startDir
	for {
		candidate := filepath.Join(dir, "go.mod")
		if _, err := os.Stat(candidate); err == nil {
			return candidate
		}
		parent := filepath.Dir(dir)
		if parent == dir {
			return ""
		}
		dir = parent
	}
}

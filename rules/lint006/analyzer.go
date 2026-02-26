package lint006

import (
	"go/ast"
	"go/parser"
	"go/token"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
	"sync"

	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/analysis/passes/inspect"
	"golang.org/x/tools/go/ast/inspector"
)

var Analyzer = &analysis.Analyzer{
	Name:     "lint006",
	Doc:      "LINT-006: functions with more than 2 return values should use a result struct",
	Run:      run,
	Requires: []*analysis.Analyzer{inspect.Analyzer},
}

type moduleFxIndex struct {
	byPkg map[string]map[string]bool
}

var moduleFxIndexCache sync.Map // module root path -> *moduleFxIndex

func run(pass *analysis.Pass) (interface{}, error) {
	insp := pass.ResultOf[inspect.Analyzer].(*inspector.Inspector)
	exemptProvided := exemptFunctionsProvidedViaFx(pass)

	nodeFilter := []ast.Node{(*ast.FuncDecl)(nil)}

	insp.Preorder(nodeFilter, func(n ast.Node) {
		fn := n.(*ast.FuncDecl)
		if fn.Type.Results == nil {
			return
		}

		count := 0
		for _, field := range fn.Type.Results.List {
			if len(field.Names) == 0 {
				count++
			} else {
				count += len(field.Names)
			}
		}

		if count > 2 {
			if fn.Recv == nil && exemptProvided[fn.Name.Name] {
				return
			}
			pass.Reportf(fn.Name.Pos(), "LINT-006: function %q has %d return values, consider using a %sResult struct", fn.Name.Name, count, fn.Name.Name)
		}
	})

	return nil, nil
}

func exemptFunctionsProvidedViaFx(pass *analysis.Pass) map[string]bool {
	exempt := providedInCurrentPackage(pass)

	moduleRoot := findModuleRoot(pass)
	if moduleRoot == "" {
		return exempt
	}
	index := getOrBuildModuleFxIndex(moduleRoot)
	for fn := range index.byPkg[pass.Pkg.Path()] {
		exempt[fn] = true
	}
	return exempt
}

func providedInCurrentPackage(pass *analysis.Pass) map[string]bool {
	exempt := make(map[string]bool)
	for _, f := range pass.Files {
		ast.Inspect(f, func(n ast.Node) bool {
			call, ok := n.(*ast.CallExpr)
			if !ok || !isFxProvideCall(call) {
				return true
			}
			for _, arg := range call.Args {
				ident, ok := arg.(*ast.Ident)
				if ok {
					exempt[ident.Name] = true
				}
			}
			return true
		})
	}
	return exempt
}

func getOrBuildModuleFxIndex(moduleRoot string) *moduleFxIndex {
	if cached, ok := moduleFxIndexCache.Load(moduleRoot); ok {
		return cached.(*moduleFxIndex)
	}
	built := buildModuleFxIndex(moduleRoot)
	actual, _ := moduleFxIndexCache.LoadOrStore(moduleRoot, built)
	return actual.(*moduleFxIndex)
}

func buildModuleFxIndex(moduleRoot string) *moduleFxIndex {
	index := &moduleFxIndex{byPkg: make(map[string]map[string]bool)}

	fset := token.NewFileSet()
	_ = filepath.WalkDir(moduleRoot, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return nil
		}
		if d.IsDir() {
			base := d.Name()
			if base == ".git" || base == ".idea" || base == "vendor" || base == "node_modules" {
				return filepath.SkipDir
			}
			return nil
		}
		if !strings.HasSuffix(path, ".go") || strings.HasSuffix(path, "_test.go") {
			return nil
		}
		file, err := parser.ParseFile(fset, path, nil, 0)
		if err != nil {
			return nil
		}

		imports := importsByAlias(file)
		ast.Inspect(file, func(n ast.Node) bool {
			call, ok := n.(*ast.CallExpr)
			if !ok || !isFxProvideCall(call) {
				return true
			}
			for _, arg := range call.Args {
				sel, ok := arg.(*ast.SelectorExpr)
				if !ok {
					continue
				}
				xid, ok := sel.X.(*ast.Ident)
				if !ok {
					continue
				}
				pkgPath, ok := imports[xid.Name]
				if !ok {
					continue
				}
				if index.byPkg[pkgPath] == nil {
					index.byPkg[pkgPath] = make(map[string]bool)
				}
				index.byPkg[pkgPath][sel.Sel.Name] = true
			}
			return true
		})
		return nil
	})

	return index
}

func importsByAlias(file *ast.File) map[string]string {
	m := make(map[string]string)
	for _, imp := range file.Imports {
		path := strings.Trim(imp.Path.Value, `"`)
		if path == "" {
			continue
		}

		var alias string
		if imp.Name != nil {
			alias = imp.Name.Name
		} else {
			parts := strings.Split(path, "/")
			alias = parts[len(parts)-1]
		}
		if alias == "_" || alias == "." {
			continue
		}
		m[alias] = path
	}
	return m
}

func isFxProvideCall(call *ast.CallExpr) bool {
	sel, ok := call.Fun.(*ast.SelectorExpr)
	if !ok {
		return false
	}
	xid, ok := sel.X.(*ast.Ident)
	if !ok {
		return false
	}
	return xid.Name == "fx" && sel.Sel.Name == "Provide"
}

func findModuleRoot(pass *analysis.Pass) string {
	if len(pass.Files) == 0 {
		return ""
	}
	start := filepath.Dir(pass.Fset.File(pass.Files[0].Pos()).Name())
	dir := start
	for {
		if _, err := os.Stat(filepath.Join(dir, "go.mod")); err == nil {
			return dir
		}
		parent := filepath.Dir(dir)
		if parent == dir {
			return ""
		}
		dir = parent
	}
}

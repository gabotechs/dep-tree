package golang

import (
	"go/ast"
	"path"
	"path/filepath"
	"strings"

	"github.com/gabotechs/dep-tree/internal/language"
)

// ParseImports TODO: refactor this.
//
//nolint:gocyclo
func (l *Language) ParseImports(file *language.FileInfo) (*language.ImportsResult, error) {
	content := file.Content.(*File)
	result := language.ImportsResult{}

	// 1. Load all the exported symbols in the current package. They can
	//    be referenced without the package name prefix in this file.
	thisPackageSymbols := make(map[string]string)
	thisDir := filepath.Dir(file.AbsPath)
	thisPackage, err := NewPackageFromDir(thisDir)
	if err != nil {
		return nil, err
	}
	for filePath, packageFile := range thisPackage.Files {
		for symbol := range packageFile.Scope.Objects {
			thisPackageSymbols[symbol] = filePath
		}
	}

	// 2. Walk all the unresolved symbols, and try to match them either with
	//    the standard library or the symbols for the current local package.
	//    if neither of those match, then the symbol is referencing an imported
	//    package. The imported package can be either a third party lib or another
	//    local package.
	resolved := map[string]struct{}{}
	for _, unresolved := range content.Unresolved {
		if _, ok := stdLibSymbols[unresolved.Name]; ok {
			// we don't care about std lib symbols.
			continue
		}

		if _, ok := resolved[unresolved.Name]; ok {
			// this was already resolved.
			continue
		}

		if filePath, ok := thisPackageSymbols[unresolved.Name]; ok {
			// The symbol comes from a file in the same dir.
			result.Imports = append(result.Imports, language.SymbolsImport([]string{unresolved.Name}, filePath))
			resolved[unresolved.Name] = struct{}{}
		}
	}

	// 3. Load all the local packages imported by the file that are not
	//    third party libraries, and that in fact are part of the codebase.
	importedPackages := make(map[string]*Package)
	for _, imp := range content.Imports {
		// TODO: what about dot imports?

		// For some reason, import statements come surrounded by quotes ('"path/filepath"').
		importPath := imp.Path.Value[1 : len(imp.Path.Value)-1]

		packagePath := l.importToPath(importPath)
		if packagePath == "" {
			continue
		}
		pkg, err := NewPackageFromDir(filepath.Join(l.Root.AbsDir, packagePath))
		if err != nil {
			result.Errors = append(result.Errors, err)
			continue
		}
		var alias string
		if imp.Name != nil {
			alias = imp.Name.Name
		} else {
			alias = importToAlias(importPath)
		}
		importedPackages[alias] = pkg
	}

	// 4. Walk the ast looking for references to imported packages.
	resolved2 := map[[2]string]struct{}{}
	for _, decl := range content.Decls {
		ast.Inspect(decl, func(node ast.Node) bool {
			selectorExpr, ok := node.(*ast.SelectorExpr)
			// 4.1 the node needs to be a `selectorExpr`.
			if !ok || selectorExpr.Sel == nil {
				return true
			}
			// 4.2 the selector element needs to be an identifier.
			libAlias, ok := selectorExpr.X.(*ast.Ident)
			if !ok {
				return true
			}

			if _, ok = resolved2[[2]string{libAlias.Name, selectorExpr.Sel.Name}]; ok {
				return true
			}

			// 4.3 the selected lib must be in the list of imported packages.
			pkg, ok := importedPackages[libAlias.Name]
			if !ok {
				return true
			}
			// 4.4 the selector identifier must be in the list of exported symbols.
			f, ok := pkg.SymbolToFile[selectorExpr.Sel.Name]
			if !ok {
				return true
			}

			result.Imports = append(result.Imports, language.SymbolsImport(
				[]string{selectorExpr.Sel.Name},
				f.AbsPath,
			))
			resolved2[[2]string{libAlias.Name, selectorExpr.Sel.Name}] = struct{}{}
			return true
		})
	}

	return &result, nil
}

// / importToPath receives an import string and attempts to retrieve the
// / folder where the package that is being imported lives. If the folder
// / is not found, it returns an empty string. For third party libraries,
// / this will always be empty. It returns the path relative to the root.
func (l *Language) importToPath(imp string) string {
	pref := l.GoMod.Module + "/"
	if strings.HasPrefix(imp, pref) {
		p := strings.TrimPrefix(imp, pref)
		return filepath.Join(strings.Split(p, "/")...)
	}
	return ""
}

// / importToAlias receives an import string and returns the alias that
// / is expected to be used in the code for accessing the contents of
// / the package.
func importToAlias(imp string) string {
	base := path.Base(imp)
	for _, split := range []string{".", "-"} {
		baseSplit := strings.Split(base, split)
		base = baseSplit[len(baseSplit)-1]
	}
	return base
}

var stdLibSymbols = map[string]bool{
	"string":     true,
	"bool":       true,
	"byte":       true,
	"rune":       true,
	"int":        true,
	"int8":       true,
	"int16":      true,
	"int32":      true,
	"int64":      true,
	"uint":       true,
	"uint8":      true,
	"uint16":     true,
	"uint32":     true,
	"uint64":     true,
	"uintptr":    true,
	"float32":    true,
	"float64":    true,
	"complex64":  true,
	"complex128": true,
	"nil":        true,
	"make":       true,
	"panic":      true,
	// TODO: there's more...
}

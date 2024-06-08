package golang

import (
	"go/ast"
	"path"
	"path/filepath"
	"strings"

	"github.com/gabotechs/dep-tree/internal/language"
)

//nolint:gocyclo
func (l *Language) ParseImports(file *language.FileInfo) (*language.ImportsResult, error) {
	content := file.Content.(*File)
	result := language.ImportsResult{}

	// 1. Load all the local packages imported by the file that are not
	//    third party libraries, and that in fact are part of the codebase.
	importedPackages := make(map[string][]*Package)
	thisModule := l.GoMod.Module + "/"
	for _, importSpec := range content.AstFile.Imports {
		importStmt := NewImportStmt(importSpec)

		if !importStmt.IsLocal(thisModule) {
			continue
		}
		pkgs, err := PackagesInDir(filepath.Join(l.Root.AbsDir, importStmt.RelPath(thisModule)))
		if err != nil {
			result.Errors = append(result.Errors, err)
			continue
		}
		importedPackages[importStmt.Alias()] = pkgs
	}

	// 2. Walk all the unresolved symbols, and try to match them with the ones
	//    exported by the current package. An unresolved symbol might be:
	//    1. a standard library identifier (string, panic, make, ...)
	//    2. a type of function declared in this same package
	//    3. a reference to an imported package (e.g. this file: `ast`, `path`, `filepath`, ...)
	//    This step resolves only symbols from 2.
	localResolutions := map[string]struct{}{}
	for _, unresolved := range content.AstFile.Unresolved {
		if _, ok := localResolutions[unresolved.Name]; ok {
			continue
		}

		pkgLookup := []*Package{content.Package}
		// Dot imports make a package's symbols available as if they were
		// local imports.
		pkgLookup = append(pkgLookup, importedPackages["."]...)

		for _, pkg := range pkgLookup {
			if f, ok := pkg.SymbolToFile[unresolved.Name]; ok {
				result.Imports = append(result.Imports, language.SymbolsImport([]string{unresolved.Name}, f.AbsPath))
				localResolutions[unresolved.Name] = struct{}{}
				break
			}
		}
	}

	// 3. Walk the ast looking for references to imported packages.
	otherPackageResolutions := map[[2]string]struct{}{}
	for _, decl := range content.AstFile.Decls {
		ast.Inspect(decl, func(node ast.Node) bool {
			selectorExpr, ok := node.(*ast.SelectorExpr)

			// 3.1 the node needs to be a `selectorExpr`.
			if !ok || selectorExpr.Sel == nil {
				return true
			}

			// 3.2 the selector element needs to be an identifier.
			libAlias, ok := selectorExpr.X.(*ast.Ident)
			if !ok {
				return true
			}

			// 3.3 this was already resolved before.
			key := [2]string{libAlias.Name, selectorExpr.Sel.Name}
			if _, ok = otherPackageResolutions[key]; ok {
				return true
			}

			// 3.4 the selected lib must be in the list of imported packages,
			//     otherwise it might be a third party library.
			pkgs, ok := importedPackages[libAlias.Name]
			if !ok {
				return true
			}

			// 3.5 the selector identifier (the right side of the dot) must be in the
			//     list symbols exported from that package.
			var absPath string
			for _, pkg := range pkgs {
				if f, ok := pkg.SymbolToFile[selectorExpr.Sel.Name]; ok {
					absPath = f.AbsPath
				}
			}

			if absPath == "" {
				return true
			}

			result.Imports = append(result.Imports, language.SymbolsImport(
				[]string{selectorExpr.Sel.Name},
				absPath,
			))
			otherPackageResolutions[key] = struct{}{}
			return true
		})
	}

	return &result, nil
}

type ImportStmt struct {
	ImportPath string
	ImportName string
}

func NewImportStmt(imp *ast.ImportSpec) ImportStmt {
	var importName string
	if imp.Name != nil {
		importName = imp.Name.Name
	}
	return ImportStmt{
		ImportPath: imp.Path.Value[1 : len(imp.Path.Value)-1],
		ImportName: importName,
	}
}

func (i *ImportStmt) IsLocal(moduleName string) bool {
	if !strings.HasSuffix(moduleName, "/") {
		moduleName += "/"
	}
	return strings.HasPrefix(i.ImportPath, moduleName)
}

func (i *ImportStmt) RelPath(moduleName string) string {
	if !strings.HasSuffix(moduleName, "/") {
		moduleName += "/"
	}
	if !i.IsLocal(moduleName) {
		return ""
	}
	return strings.TrimPrefix(i.ImportPath, moduleName)
}

func (i *ImportStmt) Alias() string {
	if i.ImportName != "" {
		return i.ImportName
	}

	base := path.Base(i.ImportPath)
	for _, split := range []string{".", "-"} {
		baseSplit := strings.Split(base, split)
		base = baseSplit[len(baseSplit)-1]
	}
	return base
}

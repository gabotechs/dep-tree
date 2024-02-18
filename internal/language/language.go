package language

// FileInfo gathers all the information related to a source file.
type FileInfo struct {
	// Content is a bucket for language implementations to inject some language-specific data in it,
	// like parsed statements.
	Content any
	// AbsPath is the absolute path of the source file.
	AbsPath string
	// RelPath is the path relative to the root of the project. Different programming languages might
	// choose to decide what is the "root of the project". For example, JS might be where the nearest
	// package.json is located.
	RelPath string
	// Package is the name of the package/module/workspace where the source file is located. Each
	// language implementation is in charge of deciding what is a package. For example, for JS/TS
	// Package might be the "name" field of the closest package.json file, for rust the name of the
	// cargo workspace where the file belongs to.
	Package string
	// Loc is the amount of lines of code a file has.
	Loc int
	// Size is the size in bytes of the file.
	Size int
}

// ImportEntry represents an import statement in a programming language.
type ImportEntry struct {
	// All is true if all the symbols from another source file are imported. Some programming languages
	// allow importing all the symbols:
	// JS     -> import * from './foo'
	// Python -> from .foo import *
	// Rust   -> use crate::foo::*;
	All bool
	// Symbols are the specific symbols that are imported from another source file. Some programming languages
	// allow to import only specific symbols:
	// JS     -> import { bar } from './foo'
	// Python -> from .foo import bar
	// Rust   -> use crate::foo::bar;
	Symbols []string
	// AbsPath is the absolute path of the source file from where are the symbols are import imported.
	// For example, having file /foo/bar.ts with the following content
	//   import { baz } from './baz'
	// will result in an ImportEntry with AbsPath = /foo/baz.ts
	AbsPath string
}

// AllImport builds an ImportEntry where all the symbols are imported.
func AllImport(absPath string) ImportEntry {
	return ImportEntry{All: true, AbsPath: absPath}
}

// EmptyImport builds an ImportEntry where nothing specific is imported, like a side effect import.
func EmptyImport(absPath string) ImportEntry {
	return ImportEntry{AbsPath: absPath}
}

// SymbolsImport builds an ImportEntry where only specific symbols are imported.
func SymbolsImport(symbols []string, absPath string) ImportEntry {
	return ImportEntry{Symbols: symbols, AbsPath: absPath}
}

// ImportsResult is the result of gathering all the import statements from a source file.
type ImportsResult struct {
	// Imports is the list of ImportEntry for the source file.
	Imports []ImportEntry
	// Errors are the non-fatal errors that occurred while parsing imports. These
	// might be rendered nicely in a UI.
	Errors []error
}

// ExportSymbol represents a symbol that it's being exported.
// Exported symbols might be aliased in certain programming languages, for example, in JS:
//
// export { foo as bar } from './file'
//
// would export something that looks like this:
//
// ExportSymbol{ Original: "foo", Alias: "bar" }
type ExportSymbol struct {
	// Original is the original value of the symbol.
	Original string
	// Alias is the alias under which the symbol is exported.
	Alias string
}

// ExportEntry represents an import statement in a programming language.
type ExportEntry struct {
	// All means that all the symbols from AbsPath are exported.
	// This is specially useful for languages that allow re-exporting all the symbols
	// from a certain location, for example, in JS:
	//
	// export * from './file'
	All bool
	// Symbols are symbols that are exported from AbsPath.
	Symbols []ExportSymbol
	// AbsPath is absolute path from where they are exported, it might be from the same file or from another.
	// Typically, this value will be the same as FileInfo.AbsPath, but in some language, you can re-export
	// symbols that are declared in different source files, for example:
	//
	// export function foo() {} // <- exports symbol `foo` that is declared in that same file.
	//
	// export foo from './file' // <= exports symbol `foo` that is declared in another file.
	AbsPath string
}

// ExportsResult is the result of gathering all the export statements from a source file.
type ExportsResult struct {
	// Exports is the list of ExportEntry in the source file.
	Exports []ExportEntry
	// Errors are the non-fatal errors that occurred while parsing exports. These
	// might be rendered nicely in a UI.
	Errors []error
}

type Language interface {
	// ParseFile receives an absolute file path and returns F, where F is the specific file implementation
	//  defined by the language. This file object F will be used as input for parsing imports and exports.
	ParseFile(path string) (*FileInfo, error)
	// ParseImports receives the file F parsed by the ParseFile method and gathers the imports that the file
	//  F contains.
	ParseImports(file *FileInfo) (*ImportsResult, error)
	// ParseExports receives the file F parsed by the ParseFile method and gathers the exports that the file
	//  F contains.
	ParseExports(file *FileInfo) (*ExportsResult, error)
}

package js

import (
	"fmt"
	"path"
	"path/filepath"
	"strings"

	"dep-tree/internal/js/grammar"
	"dep-tree/internal/utils"
)

type Import struct {
	Names   []string
	AbsPath string
}

func (p *Parser) ParseImport(unparsed []byte, dir string) (*Import, error) {
	result := Import{
		Names: make([]string, 0),
	}
	matches := grammar.ParsePathFromImport(unparsed)
	if len(matches) == 0 {
		return nil, fmt.Errorf("could not parse import importFrom from '%s'", string(unparsed))
	}
	importFrom := strings.Trim(string(matches[0]), " \n\"'")
	// 1. If import is relative.
	if importFrom[0] == '.' {
		result.AbsPath = getFileAbsPath(path.Join(dir, importFrom))
		if result.AbsPath == "" {
			return nil, fmt.Errorf("could not perform relative import for '%s' because the file or dir was not found", string(unparsed))
		}
		return &result, nil
	}
	// 2. If is imported from baseUrl.
	baseUrl := p.TsConfig.CompilerOptions.BaseUrl
	importFromBaseUrl := path.Join(p.ProjectRoot, baseUrl, importFrom)
	result.AbsPath = getFileAbsPath(importFromBaseUrl)
	if result.AbsPath != "" {
		return &result, nil
	}
	// 3. If imported from a path override.
	pathOverrides := p.TsConfig.CompilerOptions.Paths
	if pathOverrides == nil {
		return nil, nil
	}
	for pathOverride, searchPaths := range pathOverrides {
		pathOverride = strings.ReplaceAll(pathOverride, "*", "")
		if strings.HasPrefix(importFrom, pathOverride) {
			for _, searchPath := range searchPaths {
				searchPath = strings.ReplaceAll(searchPath, "*", "")
				newImportFrom := strings.ReplaceAll(importFrom, pathOverride, searchPath)
				importFromBaseUrlAndPaths := path.Join(p.ProjectRoot, baseUrl, newImportFrom)
				result.AbsPath = getFileAbsPath(importFromBaseUrlAndPaths)
				if result.AbsPath != "" {
					return &result, nil
				}
			}
			return nil, fmt.Errorf("import '%s' was matched to path '%s' in tscofing's paths option, but the resolved path did not match an existing file", importFrom, pathOverride)
		}
	}
	return nil, nil
}

type Export struct {
	SourceMap map[string]string
}

func (p *Parser) ParseExport(unparsed []byte, dir string) (*Export, error) {
	result := Export{}
	return &result, nil
}

type FileInfo struct {
	imports []*Import
	exports []*Export
}

func (p *Parser) ParseFileInfo(content []byte, dir string) (*FileInfo, error) {
	fileInfo := FileInfo{}
	for _, importMatch := range grammar.ParseImport(content) {
		parsedImport, err := p.ParseImport(importMatch, dir)
		if err != nil {
			return nil, err
		} else if parsedImport != nil {
			fileInfo.imports = append(fileInfo.imports, parsedImport)
		}
	}

	for _, exportMatch := range grammar.ParseExport(content) {
		parsedExport, err := p.ParseExport(exportMatch, dir)
		if err != nil {
			return nil, err
		} else if parsedExport != nil {
			fileInfo.exports = append(fileInfo.exports, parsedExport)
		}
	}
	return &fileInfo, nil
}

func retrieveWithExt(absPath string) string {
	for _, ext := range Extensions {
		if strings.HasSuffix(absPath, "."+ext) {
			return absPath
		}
	}
	for _, ext := range Extensions {
		withExtPath := absPath + "." + ext
		if utils.FileExists(withExtPath) {
			return withExtPath
		}
	}
	return ""
}

func getFileAbsPath(id string) string {
	absPath, err := filepath.Abs(id)
	switch {
	case err != nil:
		return ""
	case utils.DirExists(id):
		return retrieveWithExt(path.Join(absPath, "index"))
	default:
		return retrieveWithExt(absPath)
	}
}

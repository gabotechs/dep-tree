package golang

import (
	"strings"

	"github.com/gabotechs/dep-tree/internal/language"
)

func (l *Language) ParseExports(file *language.FileInfo) (*language.ExportsResult, error) {
	content := file.Content.(*File)
	results := language.ExportsResult{}
	for symbol := range content.Scope.Objects {
		if len(symbol) == 0 {
			continue
		}
		if symbol[:1] == strings.ToUpper(symbol[:1]) {
			results.Exports = append(results.Exports, language.ExportEntry{
				Symbols: []language.ExportSymbol{{
					Original: symbol,
				}},
				AbsPath: file.AbsPath,
			})
		}
	}
	return &results, nil
}

package grammar

import "regexp"

var importRegex = regexp.MustCompile(
	"import\\s+?((([\\w*\\s{},]*)\\s+from\\s+?)|)((\".*?\")|('.*?'))\\s*?(?:;|$|)",
)

var dynImportRegex = regexp.MustCompile(
	"import\\s*?\\(\\s*?((\".*?\")|('.*?'))\\s*?\\)\\s*?(?:;|$|)",
)

func ParseImport(unparsed []byte) [][]byte {
	common := importRegex.FindAll(unparsed, -1)
	dynamic := dynImportRegex.FindAll(unparsed, -1)
	return append(common, dynamic...)
}

var exportRegex = regexp.MustCompile(
	"export\\s+?((([\\w*\\s{},]*)\\s+from\\s+?)|)((\".*?\")|('.*?'))\\s*?(?:;|$|)",
)

func ParseExport(unparsed []byte) [][]byte {
	return exportRegex.FindAll(unparsed, -1)
}

var importPathRegex = regexp.MustCompile(
	"(\".*?\")|('.*?')",
)

func ParsePathFromImport(unparsed []byte) [][]byte {
	return importPathRegex.FindAll(unparsed, -1)
}

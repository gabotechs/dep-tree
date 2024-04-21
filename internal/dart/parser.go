package dart

import (
	"bufio"
	"os"
	"regexp"
	"strings"
)

var packageRegex = regexp.MustCompile(`package:[^/]+\/`)
var importRegex = regexp.MustCompile(`import\s+(['"])(.*?\.dart)`)
var exportRegex = regexp.MustCompile(`export\s+(['"])(.*?\.dart)`)

type ImportStatement struct {
	From string
}

type ExportStatement struct {
	From string
}

type Statement struct {
	Import *ImportStatement
	Export *ExportStatement
}

type File struct {
	Statements []Statement
}

func ParseFile(path string) (*File, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var fileData File
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		// Remove package patterns from the line
		line := packageRegex.ReplaceAllString(scanner.Text(), "")
		line = strings.TrimSpace(line)

		if importMatch := importRegex.FindStringSubmatch(line); importMatch != nil {
			fileData.Statements = append(fileData.Statements, Statement{
				Import: &ImportStatement{From: importMatch[2]},
			})
		} else if exportMatch := exportRegex.FindStringSubmatch(line); exportMatch != nil {
			fileData.Statements = append(fileData.Statements, Statement{
				Export: &ExportStatement{From: exportMatch[2]},
			})
		}
	}
	if err := scanner.Err(); err != nil {
		return nil, err
	}
	return &fileData, nil
}

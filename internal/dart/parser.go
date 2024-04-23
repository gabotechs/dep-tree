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
	From       string
	IsAbsolute bool
}

type ExportStatement struct {
	From       string
	IsAbsolute bool
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
		line := scanner.Text()

		// Remove comments
		if idx := strings.Index(line, "//"); idx != -1 {
			line = line[:idx]
		}

		line = strings.TrimSpace(line)

		// Remove package patterns from the line and determine if the import is absolute
		originalLine := line // Keep the original line to check for package later
		line = packageRegex.ReplaceAllString(line, "")

		// Check if the package pattern was matched to set IsAbsolute
		isAbsolute := line != originalLine

		if importMatch := importRegex.FindStringSubmatch(line); importMatch != nil {
			fileData.Statements = append(fileData.Statements, Statement{
				Import: &ImportStatement{
					From:       importMatch[2],
					IsAbsolute: isAbsolute,
				},
			})
		} else if exportMatch := exportRegex.FindStringSubmatch(line); exportMatch != nil {
			fileData.Statements = append(fileData.Statements, Statement{
				Import: &ImportStatement{ // Treat exports like imports!
					From:       exportMatch[2],
					IsAbsolute: isAbsolute,
				},
				Export: &ExportStatement{
					From:       exportMatch[2],
					IsAbsolute: isAbsolute,
				},
			})
		}
	}
	if err := scanner.Err(); err != nil {
		return nil, err
	}
	return &fileData, nil
}

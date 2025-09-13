package cpp

import (
	"bufio"
	"regexp"
	"strings"
)

// IncludeStatement represents a C++ #include directive
type IncludeStatement struct {
	// IsSystem is true for #include <header> (system headers)
	IsSystem bool
	// Header is the included header name
	Header string
	// OriginalLine is the full original include line
	OriginalLine string
}

// File represents a parsed C++ file with include statements
type File struct {
	Includes []IncludeStatement
}

var (
	// Regex patterns for matching C++ include statements
	systemIncludePattern = regexp.MustCompile(`^\s*#\s*include\s*<([^>]+)>\s*`)
	localIncludePattern  = regexp.MustCompile(`^\s*#\s*include\s*"([^"]+)"\s*`)
)

// ParseCppFile parses a C++ file and extracts include statements
func ParseCppFile(content string) (*File, error) {
	file := &File{
		Includes: make([]IncludeStatement, 0),
	}

	scanner := bufio.NewScanner(strings.NewReader(content))
	for scanner.Scan() {
		line := scanner.Text()
		
		// Check for system includes (#include <header>)
		if matches := systemIncludePattern.FindStringSubmatch(line); matches != nil {
			file.Includes = append(file.Includes, IncludeStatement{
				IsSystem:     true,
				Header:       matches[1],
				OriginalLine: line,
			})
			continue
		}
		
		// Check for local includes (#include "header")
		if matches := localIncludePattern.FindStringSubmatch(line); matches != nil {
			file.Includes = append(file.Includes, IncludeStatement{
				IsSystem:     false,
				Header:       matches[1],
				OriginalLine: line,
			})
		}
	}

	return file, scanner.Err()
}

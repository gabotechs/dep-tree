package cpp

import (
	"testing"
)

func TestParseCppFile(t *testing.T) {
	tests := []struct {
		name     string
		content  string
		expected int // expected number of includes
	}{
		{
			name: "simple local include",
			content: `#include "header.h"
int main() { return 0; }`,
			expected: 1,
		},
		{
			name: "system include",
			content: `#include <iostream>
#include <vector>
int main() { return 0; }`,
			expected: 2,
		},
		{
			name: "mixed includes",
			content: `#include <iostream>
#include "myheader.hpp"
#include <string>
#include "another.h"`,
			expected: 4,
		},
		{
			name: "no includes",
			content: `int main() {
    return 0;
}`,
			expected: 0,
		},
		{
			name: "includes with comments",
			content: `// This is a comment
#include "test.h" // inline comment
/* block comment */
#include <vector>`,
			expected: 2,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			file, err := ParseCppFile(tt.content)
			if err != nil {
				t.Fatalf("ParseCppFile() error = %v", err)
			}
			
			if len(file.Includes) != tt.expected {
				t.Errorf("ParseCppFile() got %d includes, want %d", len(file.Includes), tt.expected)
			}
		})
	}
}

func TestIncludeTypeDetection(t *testing.T) {
	content := `#include <iostream>
#include "local.h"`
	
	file, err := ParseCppFile(content)
	if err != nil {
		t.Fatalf("ParseCppFile() error = %v", err)
	}
	
	if len(file.Includes) != 2 {
		t.Fatalf("Expected 2 includes, got %d", len(file.Includes))
	}
	
	// First include should be system
	if !file.Includes[0].IsSystem {
		t.Error("Expected first include to be system include")
	}
	if file.Includes[0].Header != "iostream" {
		t.Errorf("Expected header 'iostream', got '%s'", file.Includes[0].Header)
	}
	
	// Second include should be local
	if file.Includes[1].IsSystem {
		t.Error("Expected second include to be local include")
	}
	if file.Includes[1].Header != "local.h" {
		t.Errorf("Expected header 'local.h', got '%s'", file.Includes[1].Header)
	}
}

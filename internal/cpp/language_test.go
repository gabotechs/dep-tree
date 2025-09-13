package cpp

import (
	"os"
	"path/filepath"
	"testing"
)

func TestLanguageParseFile(t *testing.T) {
	// Create a temporary test file
	tmpDir := t.TempDir()
	testFile := filepath.Join(tmpDir, "test.cpp")
	
	content := `#include <iostream>
#include "header.h"

int main() {
    return 0;
}`
	
	err := os.WriteFile(testFile, []byte(content), 0644)
	if err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}
	
	lang := NewLanguage(nil)
	fileInfo, err := lang.ParseFile(testFile)
	if err != nil {
		t.Fatalf("ParseFile() error = %v", err)
	}
	
	if fileInfo.AbsPath != testFile {
		t.Errorf("Expected AbsPath %s, got %s", testFile, fileInfo.AbsPath)
	}
	
	if fileInfo.Size == 0 {
		t.Error("Expected file size > 0")
	}
	
	if fileInfo.Loc == 0 {
		t.Error("Expected lines of code > 0")
	}
}

func TestLanguageParseImports(t *testing.T) {
	// Create a temporary directory structure
	tmpDir := t.TempDir()
	testFile := filepath.Join(tmpDir, "main.cpp")
	headerFile := filepath.Join(tmpDir, "myheader.h")
	
	// Create the header file
	err := os.WriteFile(headerFile, []byte("// Header content"), 0644)
	if err != nil {
		t.Fatalf("Failed to create header file: %v", err)
	}
	
	// Create main file that includes the header
	content := `#include <iostream>
#include "myheader.h"

int main() { return 0; }`
	
	err = os.WriteFile(testFile, []byte(content), 0644)
	if err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}
	
	lang := NewLanguage(nil)
	fileInfo, err := lang.ParseFile(testFile)
	if err != nil {
		t.Fatalf("ParseFile() error = %v", err)
	}
	
	imports, err := lang.ParseImports(fileInfo)
	if err != nil {
		t.Fatalf("ParseImports() error = %v", err)
	}
	
	// Should only have 1 import (system includes are excluded)
	if len(imports.Imports) != 1 {
		t.Errorf("Expected 1 import, got %d", len(imports.Imports))
	}
	
	if len(imports.Imports) > 0 {
		expectedPath, _ := filepath.Abs(headerFile)
		if imports.Imports[0].AbsPath != expectedPath {
			t.Errorf("Expected import path %s, got %s", expectedPath, imports.Imports[0].AbsPath)
		}
		
		if !imports.Imports[0].All {
			t.Error("Expected All to be true for C++ includes")
		}
	}
}

func TestLanguageParseExports(t *testing.T) {
	tmpDir := t.TempDir()
	
	// Test header file (should export everything)
	headerFile := filepath.Join(tmpDir, "test.h")
	err := os.WriteFile(headerFile, []byte("class MyClass {};"), 0644)
	if err != nil {
		t.Fatalf("Failed to create header file: %v", err)
	}
	
	lang := NewLanguage(nil)
	fileInfo, err := lang.ParseFile(headerFile)
	if err != nil {
		t.Fatalf("ParseFile() error = %v", err)
	}
	
	exports, err := lang.ParseExports(fileInfo)
	if err != nil {
		t.Fatalf("ParseExports() error = %v", err)
	}
	
	// Header files should export everything
	if len(exports.Exports) != 1 {
		t.Errorf("Expected 1 export entry for header file, got %d", len(exports.Exports))
	}
	
	if len(exports.Exports) > 0 && !exports.Exports[0].All {
		t.Error("Expected header file to export all")
	}
	
	// Test source file (should export nothing by default)
	sourceFile := filepath.Join(tmpDir, "test.cpp")
	err = os.WriteFile(sourceFile, []byte("int main() { return 0; }"), 0644)
	if err != nil {
		t.Fatalf("Failed to create source file: %v", err)
	}
	
	fileInfo, err = lang.ParseFile(sourceFile)
	if err != nil {
		t.Fatalf("ParseFile() error = %v", err)
	}
	
	exports, err = lang.ParseExports(fileInfo)
	if err != nil {
		t.Fatalf("ParseExports() error = %v", err)
	}
	
	// Source files should not export by default
	if len(exports.Exports) != 0 {
		t.Errorf("Expected 0 export entries for source file, got %d", len(exports.Exports))
	}
}

func TestExtensions(t *testing.T) {
	expectedExtensions := []string{"cpp", "cc", "cxx", "c++", "hpp", "hh", "hxx", "h++", "h"}
	
	if len(Extensions) != len(expectedExtensions) {
		t.Errorf("Expected %d extensions, got %d", len(expectedExtensions), len(Extensions))
	}
	
	for _, ext := range expectedExtensions {
		found := false
		for _, actualExt := range Extensions {
			if ext == actualExt {
				found = true
				break
			}
		}
		if !found {
			t.Errorf("Expected extension %s not found", ext)
		}
	}
}

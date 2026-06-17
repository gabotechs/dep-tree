package dart

import (
	"path/filepath"
	"strings"

	"github.com/gabotechs/dep-tree/internal/utils"
)

const packageScheme = "package:"
const dartScheme = "dart:"

// ResolvePath resolves an unresolved Dart import/export/part URI to an absolute file path
// based on the directory where the directive was found.
//
// It returns an empty string (without an error) when the URI points to something that should
// not be considered part of the dependency graph, for example the Dart SDK (`dart:` imports)
// or third party packages that do not belong to the current package.
func (l *Language) ResolvePath(unresolved string, dir string) (string, error) {
	switch {
	case unresolved == "":
		return "", nil

	// `dart:async`, `dart:io`, ... point to the Dart SDK, ignore them.
	case strings.HasPrefix(unresolved, dartScheme):
		return "", nil

	// `package:my_app/foo/bar.dart` points to the `lib` directory of a package.
	case strings.HasPrefix(unresolved, packageScheme):
		return l.resolvePackage(strings.TrimPrefix(unresolved, packageScheme), dir), nil

	// Anything else is a relative URI, resolved against the directory of the current file.
	default:
		return resolveFile(filepath.Join(dir, unresolved)), nil
	}
}

// resolvePackage resolves a `package:` URI without its scheme, e.g. `my_app/foo/bar.dart`.
//
// Only URIs that belong to the current package (the one declared in the closest pubspec.yaml)
// can be resolved, as third party packages are not part of the analyzed source tree.
func (l *Language) resolvePackage(pkgPath string, dir string) string {
	slashIndex := strings.Index(pkgPath, "/")
	if slashIndex < 0 {
		return ""
	}
	pkgName := pkgPath[:slashIndex]
	relPath := pkgPath[slashIndex+1:]

	pubspec := findClosestPubspec(dir)
	if pubspec == nil || pubspec.Name != pkgName {
		// Third party package or package outside of the current source tree.
		return ""
	}
	return resolveFile(filepath.Join(pubspec.absDir, "lib", relPath))
}

// resolveFile returns the absolute path of the given path if it points to an existing file,
// otherwise it returns an empty string.
func resolveFile(path string) string {
	if !utils.FileExists(path) {
		return ""
	}
	abs, err := filepath.Abs(path)
	if err != nil {
		return ""
	}
	return abs
}

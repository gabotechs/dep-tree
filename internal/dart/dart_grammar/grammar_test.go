package dart_grammar

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestParse(t *testing.T) {
	a := require.New(t)

	src := `// a leading comment
library my.lib;

import 'dart:async';
import 'package:my_app/foo.dart';
import '../utils/helper.dart' as h;
import 'helper.dart' show foo, bar;
import 'helper.dart' deferred as d hide baz;

export 'src/foo.dart';
export 'src/bar.dart' show A, B;

part 'gen.dart';
part of 'main.dart';
part of my.lib.name;

@annotation
abstract class Foo extends Bar {
  void method() {}
}

base mixin Mixy on Base {}
enum Color { red, green }
typedef IntList = List<int>;
final answer = 42;
void topLevelFn() {}
`
	dir := t.TempDir()
	path := filepath.Join(dir, "sample.dart")
	a.NoError(os.WriteFile(path, []byte(src), 0o600))

	fileInfo, err := Parse(path)
	a.NoError(err)
	file := fileInfo.Content.(*File)

	var imports []*Import
	var exports []*Export
	var parts []*Part
	var partOfs []*PartOf
	var declarations []*Declaration
	for _, stmt := range file.Statements {
		switch {
		case stmt.Import != nil:
			imports = append(imports, stmt.Import)
		case stmt.Export != nil:
			exports = append(exports, stmt.Export)
		case stmt.PartOf != nil:
			partOfs = append(partOfs, stmt.PartOf)
		case stmt.Part != nil:
			parts = append(parts, stmt.Part)
		case stmt.Declaration != nil:
			declarations = append(declarations, stmt.Declaration)
		}
	}

	a.Len(imports, 5)
	a.Equal("dart:async", imports[0].Path)
	a.Equal("package:my_app/foo.dart", imports[1].Path)
	a.Equal("../utils/helper.dart", imports[2].Path)
	a.Equal("h", imports[2].As)
	a.Equal([]string{"foo", "bar"}, imports[3].Combinators[0].Show)
	a.Equal("d", imports[4].As)
	a.Equal([]string{"baz"}, imports[4].Combinators[0].Hide)

	a.Len(exports, 2)
	a.Equal("src/foo.dart", exports[0].Path)
	a.Equal("src/bar.dart", exports[1].Path)
	a.Equal([]string{"A", "B"}, exports[1].Combinators[0].Show)

	a.Len(parts, 1)
	a.Equal("gen.dart", parts[0].Path)

	a.Len(partOfs, 2)
	a.Equal("main.dart", partOfs[0].Path)

	a.Equal([]string{"Foo", "Mixy", "Color", "IntList"}, declarationNames(declarations))
}

func declarationNames(declarations []*Declaration) []string {
	names := make([]string, len(declarations))
	for i, d := range declarations {
		names[i] = d.Name
	}
	return names
}

//nolint:govet
package dart_grammar

// Combinator represents a `show` or `hide` clause in an import or export directive,
// e.g. `show foo, bar` or `hide baz`.
type Combinator struct {
	Show []string `(  "show" @Ident ("," @Ident)*`
	Hide []string `|  "hide" @Ident ("," @Ident)* )`
}

// Import represents an import directive, e.g.
//
//	import 'package:foo/bar.dart';
//	import '../utils/helper.dart' as h;
//	import 'helper.dart' show foo, bar;
//	import 'helper.dart' deferred as h;
type Import struct {
	Path        string        `"import" @String`
	As          string        `"deferred"? ("as" @Ident)?`
	Combinators []*Combinator `@@* ";"`
}

// Export represents an export directive, e.g.
//
//	export 'src/foo.dart';
//	export 'src/foo.dart' show foo;
type Export struct {
	Path        string        `"export" @String`
	Combinators []*Combinator `@@* ";"`
}

// Part represents a `part` directive, e.g. `part 'foo.dart';`.
type Part struct {
	Path string `"part" @String ";"`
}

// PartOf represents a `part of` directive, e.g.
//
//	part of 'foo.dart';
//	part of my.library.name;
type PartOf struct {
	Path string `"part" "of" (@String | (Ident ("." Ident)*)) ";"`
}

// Declaration represents a top-level declaration that exposes a public symbol, e.g.
//
//	class Foo {}
//	abstract class Bar {}
//	enum Color { red, green }
//	mixin Baz {}
//	typedef IntList = List<int>;
type Declaration struct {
	Name string `("abstract" | "base" | "final" | "sealed" | "interface")* ("class" | "enum" | "mixin" | "typedef") @Ident`
}

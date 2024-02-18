# Implementing new languages in Dep Tree

Implementing a new language in Dep Tree boils down to writing code that satisfies the
[`Language` interface](../internal/language/language.go) and wiring it up to the appropriate
file extensions. There's three core methods that should be provided:

- ParseFile: parses a file object given its path.
- ParseImports: given a parsed file, retrieves the imported symbols, like functions, classes or variables.
- ParseExports: given a parsed file, retrieves the exported symbols.

As long as implementations are able to satisfy this interface, they will be compatible with dep-tree's
machinery for creating graphs and analyzing dependencies.

## Learn by example

First, clone Dep Tree's repository, as we will work directly committing files to it:

```shell
git clone https://github.com/gabotechs/dep-tree
```

Then, ensure you have Golang set-up in your machine. Dep Tree is written in Golang, so this
tutorial will assume that you have the compiler installed and that you have some basic knowledge
of the language.

### The Dummy Language

In order to keep it simple, we will create a fictional programming language
that only has `import` and `export` statements, we will call it "Dummy Language", and its file
extension will be `.dl`.

Dummy Language files will have statements like this:
```js
import foo from ./file.dl

export bar
```

In this file, the first statement imports symbol `foo` from the file `file.dl` located in the
same folder, and the second statement is exporting the symbol `bar`. We can expect `file.dl`
to contain something like this:
```js
export foo
```
Where foo is the symbol that the other file is trying to import.

### 1. Parsing files

First, we will need to create a parser for our Dummy Language. There are many tools in Golang for
creating language parsers, but most language implementations in Dep Tree use https://github.com/alecthomas/participle,
which allows writing parsers with very few lines of code.

Navigate to Dep Tree's cloned repository, and create a directory under the `internal` folder called `dummy`.
Create a file called `parser.go` inside `internal/dummy`, where we will place our parser code:

```go
package dummy

import (
	"github.com/alecthomas/participle/v2"
	"github.com/alecthomas/participle/v2/lexer"
)

type ImportStatement struct {
	Symbols []string `"import" @Ident ("," @Ident)*`
	From    string   `"from" @(Ident|Punctuation|"/")*`
}

type ExportStatement struct {
	Symbol string `"export" @Ident`
}

type Statement struct {
	Import *ImportStatement `@@ |`
	Export *ExportStatement `@@`
}

type File struct {
	Statements []Statement `@@*`
}

var (
	lex = lexer.MustSimple(
		[]lexer.SimpleRule{
			{"KewWord", "(export|import|from)"},
			{"Punctuation", `[,\./]`},
			{"Ident", `[a-zA-Z]+`},
			{"Whitespace", `\s+`},
		},
	)
	parser = participle.MustBuild[File](
		participle.Lexer(lex),
		participle.Elide("Whitespace"),
	)
)
```

We will not cover here how [participle](https://github.com/alecthomas/participle) works, but
it's important to note that using it is not required. If you are implementing a new language
for Dep Tree, feel free to choose the parsing mechanism that you find most suitable.

We now need to implement the `ParseFile` method from the `Language` interface. 

We will place all our methods in a file called `language.go` inside the `internal/dummy` dir:

```go
package dummy

import (
	"bytes"
	"os"
	"path/filepath"

	"github.com/gabotechs/dep-tree/internal/language"
)

type Language struct{}

func (l *Language) ParseFile(path string) (*language.FileInfo, error) {
	//TODO implement me
	panic("implement me")
}
```

The ultimate goal of the `ParseFile` method is to output a `FileInfo` struct, that contains
information about the source file itself, like its size, the amount of lines of code it has, its
parsed statements, it's path on the disk...

A fully working implementation of this method would look like this:
```go
func (l *Language) ParseFile(path string) (*language.FileInfo, error) {
	content, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	file, err := parser.ParseBytes(path, content)
	if err != nil {
		return nil, err
	}
	currentDir, _ := os.Getwd()
	relPath, _ := filepath.Rel(currentDir, path)
	return &language.FileInfo{
		Content: file.Statements, // dump the parsed statements into the FileInfo struct.
		Loc:     bytes.Count(content, []byte("\n")), // get the amount of lines of code.
		Size:    len(content), // get the size of the file in bytes.
		AbsPath: path, // provide its absolute path.
		RelPath: relPath, // provide the path relative to the current dir.
	}, nil
}
```
The `RelPath` attribute is important as it's what ultimately will be shown while rendering the graph.
Some language implementations choose to provide a path not relative to the current working directory, 
but to its closest `package.json` for example. Language implementation are free to choose what `RelPath`
should look like.

### 2. Parsing Import statements

Parsing imports is far simpler, as we have everything in place already. 

This method accepts the same `FileInfo` structure that we created previously in the `ParseFile` method,
and returns an `ImportResult` structure with all the import statements gathered from the file.

We will place our method implementation in the same `language.go` file, just below the `ParseFile` method:

```go
func (l *Language) ParseImports(file *language.FileInfo) (*language.ImportsResult, error) {
	var result language.ImportsResult

	for _, statement := range file.Content.([]Statement) {
		if statement.Import != nil {
			result.Imports = append(result.Imports, language.ImportEntry{
				Symbols: statement.Import.Symbols,
				// in our Dummy language, imports are always relative to source file.
				AbsPath: filepath.Join(filepath.Dir(file.AbsPath), statement.Import.From),
			})
		}
	}

	return &result, nil
}
```

### 3. Parsing Export statements

The `ParseExports` method is very similar to the `ParseImports` method, but it gathers export statements rather
than import statements.

```go
func (l *Language) ParseExports(file *language.FileInfo) (*language.ExportsResult, error) {
	var result language.ExportsResult

	for _, statement := range file.Content.([]Statement) {
		if statement.Export != nil {
			result.Exports = append(result.Exports, language.ExportEntry{
				// our Dummy Language only allows exporting 1 symbol at a time, and does not support aliasing.
				Symbols: []language.ExportSymbol{{Original: statement.Export.Symbol}},
				AbsPath: file.AbsPath,
			})
		}
	}

	return &result, nil
}
```

### 4. Wiring up the language with Dep Tree

Now that the `Language` interface is fully implemented, we need to wire it up so that it's recognized by
Dep Tree. For that, let's declare the array of extensions that the Dummy Language supports in the 
`internal/dummy/language.go` file:

```go
var Extensions = []string{"dl"}
```

Now, we will need to go to `cmd/root.go` and tweak the `inferLang` function in order to also take `.dl` files
into account. Beware that this function is highly susceptible to changing, so the following instructions
might not be accurate:

- Add one more entry to the `score` struct:
```go
	score := struct {
		js     int
		python int
		rust   int
		+ dummy  int // <- add this
	}{}
```
- Add one case branch in the `for` loop:
```go
		+ case utils.EndsWith(file, dummy.Extensions):
		+     score.dummy += 1
		+     if score.dummy > top.v {
		+         top.v = score.dummy
		+         top.lang = "dummy"
		+     }
```
- Add one case branch at the bottom of the function
```go
	+ case "dummy":
    +     return &dummy.Language{}, nil
```

### 5. Running Dep Tree on the Dummy Language

You have everything in place to start playing with the Dummy Language and Dep Tree.
- Compile Dep Tree by running `go build` in the root directory of the project
- Create some Dummy Language files that import each other
- use the generated binary `./dep-tree` and run them on one of the Dummy Language files

If everything went correctly, you should be seeing a graph that renders your files.

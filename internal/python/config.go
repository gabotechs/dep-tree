package python

type Config struct {
	PythonPath                 []string `yaml:"pythonPath"`
	ExcludeConditionalImports  bool     `yaml:"excludeConditionalImports"`
	IgnoreFromImportsAsExports bool
	IgnoreDirectoryImports     bool
}

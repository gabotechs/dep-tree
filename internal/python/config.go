package python

type Config struct {
	ExcludeConditionalImports  bool `yaml:"excludeConditionalImports"`
	IgnoreFromImportsAsExports bool
	IgnoreDirectoryImports     bool
}

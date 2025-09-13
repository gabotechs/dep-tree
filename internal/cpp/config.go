package cpp

type Config struct {
	IncludeDirs []string `yaml:"include_dirs"`

	ExcludeSystemHeaders bool `yaml:"exclude_system_headers"`

	HeaderExtensions []string `yaml:"header_extensions"`

	SourceExtensions []string `yaml:"source_extensions"`
}

func DefaultConfig() *Config {
	return &Config{
		IncludeDirs:          []string{},
		ExcludeSystemHeaders: true,
		HeaderExtensions:     []string{".h", ".hpp", ".hh", ".hxx", ".h++"},
		SourceExtensions:     []string{".cpp", ".cc", ".cxx", ".c++"},
	}
}

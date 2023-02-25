package language

type ExportsResult struct {
	// Exports: map from exported name to the absolute path from where it is exported
	//  NOTE: even though it could work returning a path relative to the file, it should return absolute
	Exports map[string]string
	// Errors: errors while parsing exports
	Errors []error
}

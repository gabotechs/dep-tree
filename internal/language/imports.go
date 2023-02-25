package language

import "github.com/elliotchance/orderedmap/v2"

type ImportsResult struct {
	// Imports: ordered map from absolute imported path to the array of names that where imported.
	//  if one of the names is *, then all the names are imported
	Imports *orderedmap.OrderedMap[string, []string]
	// Errors: errors while parsing imports
	Errors []error
}

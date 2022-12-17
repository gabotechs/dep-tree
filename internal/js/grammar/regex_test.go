package grammar

import (
	"os"
	"path"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestImportRegex(t *testing.T) {
	tests := []struct {
		Name     string
		Content  string
		Expected []string
	}{
		{
			Name:    "import regex",
			Content: "import-regex.js",
			Expected: []string{
				"import {\n    Component\n} from '@angular2/core';",
				"import defaultMember from \"module-name\";",
				"import   *    as name from \"module-name  \";",
				"import   {  member }   from \"  module-name\";",
				"import { member as alias } from \"module-name\";",
				"import { member1 , member2 } from \"module-name\";",
				"import { member1 , member2 as alias2 , member3 as alias3 } from \"module-name\";",
				"import {\n    Component\n} from '@angular2/core';",
				"import defaultMember from \"$module-name\";",
				"import defaultMember, { member, member } from \"module-name\";",
				"import defaultMember, * as name from \"module-name\";",
				"import   *    as name from \"module-name  \"",
				"import   {  member }   from \"  module-name\"",
				"import { member as alias } from \"module-name\"",
				"import { member1 , member2 } from \"module-name\"",
				"import { member1 , member2 as alias2 , member3 as alias3 } from \"module-name\"",
				"import {\n    Component\n} from '@angular2/core'",
				"import defaultMember from \"$module-name\"",
				"import defaultMember, { member, member } from \"module-name\"",
				"import defaultMember, * as name from \"module-name\"",
				"import \"module-name\";",
				"import React from \"react\"",
				"import { Field } from \"redux-form\"",
				"import \"module-name\";",
				"import {\n    PlaneBufferGeometry,\n    OctahedronGeometry,\n    TorusBufferGeometry\n} from '../geometries/Geometries.js';",
				"import {\n    PlaneBufferGeometry,\n    OctahedronGeometry,\n    TorusBufferGeometry\n} from '../geometries/Geometries.js'",
				"import { Field } from \"redux-form\";",
				"import MultiContentListView from \"./views/ListView\";",
				"import MultiContentAddView from \"./views/AddView\";",
				"import MultiContentEditView from \"./views/EditView\";",
				"import { Field } from \"redux-form\"",
				"import MultiContentListView from \"./views/ListView\"",
				"import MultiContentAddView from \"./views/AddView\"",
				"import MultiContentEditView from \"./views/EditView\"",
				"import(\"whatever.js\");",
				"import('whatever.js')",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.Name, func(t *testing.T) {
			a := require.New(t)
			var content []byte
			if strings.HasSuffix(tt.Content, ".js") {
				var err error
				content, err = os.ReadFile(path.Join("regex_test", tt.Content))
				a.NoError(err)
			} else {
				content = []byte(tt.Content)
			}
			results := make([]string, 0)
			for _, importMatch := range ParseImport(content) {
				results = append(results, string(importMatch))
			}
			a.Equal(tt.Expected, results)
		})
	}
}

func TestImportPathRegex(t *testing.T) {
	tests := []struct {
		Name     string
		Content  string
		Expected string
	}{
		{
			Name:     "import path regex 1",
			Content:  "import {\n    Component\n} from '@angular2/core';",
			Expected: "'@angular2/core'",
		},
		{
			Name:     "import path regex 2",
			Content:  "import { Field } from \"redux-form\"",
			Expected: "\"redux-form\"",
		},
	}

	for _, tt := range tests {
		t.Run(tt.Name, func(t *testing.T) {
			a := require.New(t)
			content := []byte(tt.Content)
			matches := ParsePathFromImport(content)
			a.Equal(tt.Expected, string(matches[0]))
		})
	}
}

package grammar

import (
	"os"
	"path"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestExport(t *testing.T) {
	tests := []struct {
		Name                string
		ExpectedDeclaration []string
		ExpectedList        []AliasedName
		ExpectedDefault     []bool
		ExpectedProxy       []*ProxyExport
	}{
		{
			Name:                `export let name1, name2;`,
			ExpectedDeclaration: []string{"name1"},
		},
		{
			Name:                `export const name1 = 1`,
			ExpectedDeclaration: []string{"name1"},
		},
		{
			Name:                `export function functionName()`,
			ExpectedDeclaration: []string{"functionName"},
		},
		{
			Name:                `export async function functionName()`,
			ExpectedDeclaration: []string{"functionName"},
		},
		{
			Name:                `export class ClassName { /* … */ }`,
			ExpectedDeclaration: []string{"ClassName"},
		},
		{
			Name:                `export function* generatorFunctionName() { /* … */ }`,
			ExpectedDeclaration: []string{"generatorFunctionName"},
		},
		{
			Name:                `export type MyType`,
			ExpectedDeclaration: []string{"MyType"},
		},
		{
			Name:                `export interface MyInterface`,
			ExpectedDeclaration: []string{"MyInterface"},
		},
		{
			Name: `export { name1, nameN };`,
			ExpectedList: []AliasedName{
				{
					Original: "name1",
				},
				{
					Original: "nameN",
				},
			},
		},
		{
			Name: `export { variable1 as name1, variable2 as name2, nameN };`,
			ExpectedList: []AliasedName{
				{
					Original: "variable1",
					Alias:    "name1",
				},
				{
					Original: "variable2",
					Alias:    "name2",
				},
				{
					Original: "nameN",
				},
			},
		},
		{
			Name: `export { name1 as default /*, … */ };`,
			ExpectedList: []AliasedName{
				{
					Original: "name1",
					Alias:    "default",
				},
			},
		},
		{
			Name:            `export default expression;`,
			ExpectedDefault: []bool{true},
		},
		{
			Name:            `export default function functionName() { /* … */ }`,
			ExpectedDefault: []bool{true},
		},
		{
			Name:            `export default class ClassName { /* … */ }`,
			ExpectedDefault: []bool{true},
		},
		{
			Name:            `export default function* generatorFunctionName() { /* … */ }`,
			ExpectedDefault: []bool{true},
		},
		{
			Name:            `export default function () { /* … */ }`,
			ExpectedDefault: []bool{true},
		},
		{
			Name:            `export default class { /* … */ }`,
			ExpectedDefault: []bool{true},
		},
		{
			Name:            `export default function* () { /* … */ }`,
			ExpectedDefault: []bool{true},
		},
		{
			Name: `export * from "module-name";`,
			ExpectedProxy: []*ProxyExport{{
				ExportAll: true,
				From:      "module-name",
			}},
		},
		{
			Name: `export * as name1 from "module-name";`,
			ExpectedProxy: []*ProxyExport{{
				ExportAll:      true,
				ExportAllAlias: "name1",
				From:           "module-name",
			}},
		},
		{
			Name: `export { name1, /* …, */ nameN } from "module-name";`,
			ExpectedProxy: []*ProxyExport{{
				ExportDeconstruction: &ExportDeconstruction{
					Names: []AliasedName{
						{
							Original: "name1",
						},
						{
							Original: "nameN",
						},
					},
				},
				From: "module-name",
			}},
		},
		{
			Name: `export { import1 as name1, import2 as name2, /* …, */ nameN } from "module-name";`,
			ExpectedProxy: []*ProxyExport{{
				ExportDeconstruction: &ExportDeconstruction{
					Names: []AliasedName{
						{
							Original: "import1",
							Alias:    "name1",
						},
						{
							Original: "import2",
							Alias:    "name2",
						},
						{
							Original: "nameN",
						},
					},
				},
				From: "module-name",
			}},
		},
		{
			Name: `export { default } from "module-name"; `,
			ExpectedProxy: []*ProxyExport{{
				ExportDeconstruction: &ExportDeconstruction{
					Names: []AliasedName{
						{
							Original: "default",
						},
					},
				},
				From: "module-name",
			}},
		},
	}

	for _, tt := range tests {
		t.Run(tt.Name, func(t *testing.T) {
			a := require.New(t)
			var content []byte
			if strings.HasSuffix(tt.Name, ".js") {
				var err error
				content, err = os.ReadFile(path.Join(".export_test", tt.Name))
				a.NoError(err)
			} else {
				content = []byte(tt.Name)
			}
			parsed, err := parser.ParseBytes("", content)
			a.NoError(err)

			var declarationResults []string
			var listResults []AliasedName
			var defaultResults []bool
			var proxyResult []*ProxyExport
			for _, stmt := range parsed.Statements {
				switch {
				case stmt.DeclarationExport != nil:
					declarationResults = append(declarationResults, stmt.DeclarationExport.Name)
				case stmt.ListExport != nil:
					listResults = append(listResults, stmt.ListExport.ExportDeconstruction.Names...)
				case stmt.DefaultExport != nil:
					defaultResults = append(defaultResults, stmt.DefaultExport.Default)
				case stmt.ProxyExport != nil:
					proxyResult = append(proxyResult, stmt.ProxyExport)
				}
			}
			a.Equal(tt.ExpectedDeclaration, declarationResults)
			a.Equal(tt.ExpectedList, listResults)
			a.Equal(tt.ExpectedDefault, defaultResults)

			if tt.ExpectedProxy == nil {
				a.Equal(tt.ExpectedProxy, proxyResult)
			} else {
				a.NotEqual(nil, proxyResult)
				for i, expectedProxy := range tt.ExpectedProxy {
					a.Greater(len(proxyResult), i)
					actualProxy := proxyResult[i]
					a.Equal(expectedProxy.From, actualProxy.From)
					a.Equal(expectedProxy.ExportAll, actualProxy.ExportAll)
					a.Equal(expectedProxy.ExportAllAlias, actualProxy.ExportAllAlias)
					if expectedProxy.ExportDeconstruction != nil {
						a.NotEqual(nil, actualProxy.ExportDeconstruction)
						a.Equal(expectedProxy.ExportDeconstruction.Names, actualProxy.ExportDeconstruction.Names)
					}
				}
			}
		})
	}
}

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
		ExpectedProxy       []string
	}{
		{
			Name:                `export let name1, name2;`,
			ExpectedDeclaration: []string{"name1"},
		},
		{
			Name:                `export const name1 = 1`,
			ExpectedDeclaration: []string{},
		},
		{
			Name:                `export function functionName() { /* … */ }`,
			ExpectedDeclaration: []string{},
		},
		{
			Name:                `export class ClassName { /* … */ }`,
			ExpectedDeclaration: []string{},
		},
		{
			Name:                `export function* generatorFunctionName() { /* … */ }`,
			ExpectedDeclaration: []string{},
		},
		{
			Name:                `export const { name1, name2: bar } = o;`,
			ExpectedDeclaration: []string{},
		},
		{
			Name:                `export const [ name1, name2 ] = array;`,
			ExpectedDeclaration: []string{},
		},
		{
			Name:                `export { name1, /* …, */ nameN };`,
			ExpectedDeclaration: []string{},
		},
		{
			Name:                `export { variable1 as name1, variable2 as name2, /* …, */ nameN };`,
			ExpectedDeclaration: []string{},
		},
		{
			Name:                `export { variable1 as "string name" };`,
			ExpectedDeclaration: []string{},
		},
		{
			Name:                `export { name1 as default /*, … */ };`,
			ExpectedDeclaration: []string{},
		},
		{
			Name:                `export default expression;`,
			ExpectedDeclaration: []string{},
		},
		{
			Name:                `export default function functionName() { /* … */ }`,
			ExpectedDeclaration: []string{},
		},
		{
			Name:                `export default class ClassName { /* … */ }`,
			ExpectedDeclaration: []string{},
		},
		{
			Name:                `export default function* generatorFunctionName() { /* … */ }`,
			ExpectedDeclaration: []string{},
		},
		{
			Name:                `export default function () { /* … */ }`,
			ExpectedDeclaration: []string{},
		},
		{
			Name:                `export default class { /* … */ }`,
			ExpectedDeclaration: []string{},
		},
		{
			Name:                `export default function* () { /* … */ }`,
			ExpectedDeclaration: []string{},
		},
		{
			Name:                `export * from "module-name";`,
			ExpectedDeclaration: []string{},
		},
		{
			Name:                `export * as name1 from "module-name";`,
			ExpectedDeclaration: []string{},
		},
		{
			Name:                `export { name1, /* …, */ nameN } from "module-name";`,
			ExpectedDeclaration: []string{},
		},
		{
			Name:                `export { import1 as name1, import2 as name2, /* …, */ nameN } from "module-name";`,
			ExpectedDeclaration: []string{},
		},
		{
			Name:                `export { default, /* …, */ } from "module-name";	`,
			ExpectedDeclaration: []string{},
		},
		{
			Name:                "export.js",
			ExpectedDeclaration: []string{},
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
			var proxyResults []string
			for _, stmt := range parsed.Statements {
				if stmt.DeclarationExport != nil {
					declarationResults = append(declarationResults, stmt.StaticImport.Path)
				} else if stmt.DynamicImport != nil {
					proxyResults = append(proxyResults, stmt.DynamicImport.Path)
				}
			}
			a.Equal(tt.ExpectedDeclaration, declarationResults)
			a.Equal(tt.ExpectedProxy, proxyResults)
		})
	}
}

// Package errormsg adds context field path to UserType validation errors.
package errormsg

import (
	"fmt"
	"path/filepath"
	"strings"

	"github.com/umisama/go-regexpcache"
	"goa.design/goa/v3/codegen"
	"goa.design/goa/v3/eval"
	"goa.design/goa/v3/expr"
	httpgen "goa.design/goa/v3/http/codegen"

	"github.com/keboola/keboola-as-code/internal/pkg/utils/strhelper"
)

// nolint: gochecknoinits
func init() {
	codegen.RegisterPluginFirst("errormsg", "gen", prepare, generate)
}

func prepare(_ string, roots []eval.Root) error {
	for _, root := range roots {
		root.WalkSets(func(s eval.ExpressionSet) error {
			for _, e := range s {
				if v, ok := e.(*expr.HTTPServiceExpr); ok {
					httpData := httpgen.HTTPServices.Get(v.Name())

					// Endpoint requests
					for _, e := range httpData.Endpoints {
						if e.Payload != nil && e.Payload.Request != nil && e.Payload.Request.ServerBody != nil {
							modifyTypeValidation(e.Payload.Request.ServerBody)
						}
					}

					// User defined types
					for _, t := range httpData.ServerBodyAttributeTypes {
						modifyTypeValidation(t)
					}
				}
			}
			return nil
		})
	}
	return nil
}

func generate(_ string, _ []eval.Root, files []*codegen.File) ([]*codegen.File, error) {
	for _, f := range files {
		// nolint: forbidigo
		if filepath.Base(f.Path) == "types.go" {
			for _, s := range f.SectionTemplates {
				if s.Name == "source-header" {
					codegen.AddImport(s, &codegen.ImportSpec{Path: "fmt"}, &codegen.ImportSpec{Path: "strings"})
				}
				if s.Name == "server-validate" {
					s.Source = strings.ReplaceAll(
						s.Source,
						"func Validate{{ .VarName }}(body {{ .Ref }}) (err error)",
						"func Validate{{ .VarName }}(body {{ .Ref }}, errContext []string) (err error)",
					)
				}
			}
		}
	}
	return files, nil
}

func modifyTypeValidation(t *httpgen.TypeData) {
	// Call the type validation with errContext []string, add parameter
	t.ValidateRef = strings.ReplaceAll(t.ValidateRef, `(v)`, `(v, errContext)`)
	t.ValidateRef = strings.ReplaceAll(t.ValidateRef, `(&body)`, `(&body, []string{"body"})`)

	// Use errContext in goa.*Error constructors
	t.ValidateDef = regexpcache.
		MustCompile(`goa\.[a-zA-Z0-9]+Error\([^\n]+\)`).
		ReplaceAllStringFunc(t.ValidateDef, func(call string) string {
			return regexpcache.
				MustCompile(`"body(\.[^"]+)?"`).
				ReplaceAllStringFunc(call, func(param string) string {
					param = strings.TrimPrefix(param, `"body`)
					param = strings.TrimPrefix(param, `.`)
					param = strings.TrimSuffix(param, `"`)
					if len(param) == 0 {
						return `strings.Join(errContext, ".")`
					} else {
						return `strings.Join(append(errContext, "` + param + `"), ".")`
					}
				})
		})

	// Add context argument to nested Validate* calls
	t.ValidateDef = regexpcache.
		MustCompile(`:= Validate[^()]+\([^()]+\)`).
		ReplaceAllStringFunc(t.ValidateDef, func(s string) string {
			s = strings.TrimSuffix(s, `)`)
			return s + ", errContext)"
		})

	// Add errContext to nested object validation calls
	{
		re := regexpcache.MustCompile(`(if err2 := Validate[^()]+\(body.)([^ {}]+)(, errContext)(\); err2 != nil {)`)
		t.ValidateDef = re.ReplaceAllStringFunc(t.ValidateDef, func(s string) string {
			m := re.FindStringSubmatch(s)
			field := strhelper.FirstLower(m[2])
			return m[1] + m[2] + fmt.Sprintf(", append(errContext, \"%s\")", field) + m[4]
		})
	}

	// Add errContext to nested array validation calls
	{
		re := regexpcache.MustCompile(`(for _, e := range body\.)([^ {}]+)( {)`)
		t.ValidateDef = re.ReplaceAllStringFunc(t.ValidateDef, func(s string) string {
			m := re.FindStringSubmatch(s)
			field := strhelper.FirstLower(m[2])
			return fmt.Sprintf("for i, e := range body.%s {\nerrContext := append(errContext, fmt.Sprintf(`%s[%%d]`, i))", m[2], field)
		})
	}
}

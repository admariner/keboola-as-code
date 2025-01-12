// Package dependencies contains extension to inject dependencies to the service endpoint handlers.
package dependencies

import (
	"path/filepath"
	"strings"

	"goa.design/goa/v3/codegen"
	"goa.design/goa/v3/eval"
)

func RegisterPlugin(pkgPath string) {
	addPackageImport := func(s *codegen.SectionTemplate) {
		data := s.Data.(map[string]interface{})
		imports := data["Imports"].([]*codegen.ImportSpec)
		imports = append(imports, &codegen.ImportSpec{Name: "dependencies", Path: pkgPath})
		data["Imports"] = imports
	}

	generate := func(genpkg string, roots []eval.Root, files []*codegen.File) ([]*codegen.File, error) {
		for _, f := range files {
			// nolint: forbidigo
			switch filepath.Base(f.Path) {
			case "service.go":
				for _, s := range f.SectionTemplates {
					switch s.Name {
					case "source-header":
						// Import dependencies package
						addPackageImport(s)
					case "service":
						// Add dependencies to the service interface, instead of context (it is included in dependencies)
						search := `{{ .VarName }}(context.Context`
						replace := `{{ .VarName }}(context.Context, 
{{- $authFound := false}}
{{- range .Requirements }}
	{{- range .Schemes }}
		{{- if eq .Type "APIKey" -}}
			dependencies.ProjectRequestScope
			{{- $authFound = true}}
			{{- break}}
		{{- end }}
	{{- end }}
{{- end }}
{{- if eq $authFound false -}}
dependencies.PublicRequestScope
{{- end -}}
	`
						s.Source = strings.ReplaceAll(s.Source, search, replace)
					}
				}
			case "endpoints.go":
				for _, s := range f.SectionTemplates {
					switch s.Name {
					case "source-header":
						// Import dependencies package
						addPackageImport(s)
					case "endpoint-method":

						search := `
{{- if .ServerStream }}
	`
						replace := `
{{- $authFound := false}}
{{- range .Requirements }}
	{{- range .Schemes }}
		{{- if eq .Type "APIKey" }}
			deps := ctx.Value(dependencies.ProjectRequestScopeCtxKey).(dependencies.ProjectRequestScope)
			{{- $authFound = true}}
			{{- break}}
		{{- end }}
	{{- end }}
{{- end }}
{{- if eq $authFound false }}
	deps := ctx.Value(dependencies.PublicRequestScopeCtxKey).(dependencies.PublicRequestScope)
{{- end }}
{{- if .ServerStream }}
	`
						s.Source = strings.ReplaceAll(s.Source, search, replace)

						// Add dependencies to the service method call
						s.Source = strings.ReplaceAll(
							s.Source,
							"s.{{ .VarName }}(ctx",
							"s.{{ .VarName }}(ctx, deps",
						)
					}
				}
			}
		}
		return files, nil
	}

	codegen.RegisterPluginFirst("api-dependencies", "gen", nil, generate)
}

{{/* ---- Schema meta rendering ---------------------------------------------------------------*/}}
// GenSchemaMeta{{ .Name }} returns schema for {{.Name}}
//
{{.Comment}}
func GenSchemaMeta{{ .Name }}() map[string]*accessors.SchemaMeta {
	return {{ template "fieldsSchema" .Fields -}}
}

{{- define "fieldsSchema" -}}
map[string]*accessors.SchemaMeta {
{{- range $index, $field := . }}
	{{.Comment}}
	"{{ .NameSnake }}": {{ template "fieldSchema" . }}    
{{- end }}
}
{{- end -}}

{{- define "fieldSchema" -}}
{
	Name: {{.Name|quote}},
	IsTime: {{.IsTime}},
	IsDuration: {{.IsDuration}},
	{{ if eq .Kind "CUSTOM_TYPE" -}}
	FromTerraform: FromTerraform{{.Suffix}},
	ToTerraform: ToTerraform{{.Suffix}},
	{{- end }}
	{{- if eq .Kind "SINGULAR_MESSAGE" "REPEATED_MESSAGE" -}}
	Nested: {{ template "fieldsSchema" .Message.Fields }},
	{{- end }}
	{{- if eq .Kind "MESSSAGE_MAP" -}}
	Nested: {{ template "fieldsSchema" .MapValueField.Message.Fields }},
	{{- end }}
},
{{- end -}}

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
    {{if .Comment}}{{.Comment}}{{else}}{{if .Message}}{{.Message.Comment}}{{end}}{{end}}
	"{{ .NameSnake }}": {{ template "fieldSchema" . }}    
{{- end }}
}
{{- end -}}

{{- define "fieldSchema" -}}
{{- if eq .Kind "REPEATED_MESSAGE" }}
{
    {{ template "repeatedMessage" . }}
},
{{- end }}

{{- if eq .Kind "REPEATED_ELEMENTARY" }}
{
    {{ template "repeatedElementary" . }}
},
{{- end }}

{{- if eq .Kind "MAP" }}
{
    {{ template "map" . }}
},
{{- end }}

{{- if eq .Kind "MESSSAGE_MAP" }}
{
    {{ template "messageMap" . }}
},
{{- end }}

{{- if eq .Kind "SINGULAR_MESSAGE" }}
{
    Name: {{.Name|quote}},
    Nested: {{ template "fieldsSchema" .Message.Fields }},
},
{{- end }}

{{- if eq .Kind "SINGULAR_ELEMENTARY" }}
{
    {{ template "singularElementary" . }}
},
{{- end }}

{{- if eq .Kind "CUSTOM_TYPE" }}
{
    Name: {{.Name|quote}},
},
{{- end }}
{{- end -}}

{{- define "singularElementary" -}}
Name: {{.Name|quote}},
IsTime: {{.IsTime}},
IsDuration: {{.IsDuration}},
{{- end -}}

{{- define "repeatedMessage" -}}
Name: {{.Name|quote}},
Nested: {{ template "fieldsSchema" .Message.Fields }},
{{- end -}}

{{- define "repeatedElementary" -}}
Name: {{.Name|quote}},
IsTime: {{.IsTime}},
IsDuration: {{.IsDuration}},
{{- end -}}

{{- define "map" -}}
Name: {{.Name|quote}},
IsTime: {{.IsTime}},
IsDuration: {{.IsDuration}},
{{- end -}}

{{- define "messageMap" -}}
Name: {{.Name|quote}},
Nested: {{ template "fieldsSchema" .MapValueField.Message.Fields }},
{{- end -}}
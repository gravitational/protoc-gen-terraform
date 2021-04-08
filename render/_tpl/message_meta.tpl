{{/* ---- Schema meta rendering ---------------------------------------------------------------*/}}
// GenSchemaMeta{{ .Name }} returns schema for {{.Name}}
//
{{.Comment}}
func GenSchemaMeta{{ .Name }}() map[string]*SchemaMeta {
	return {{ template "fieldsSchema" .Fields -}}
}

{{- define "fieldsSchema" -}}
map[string]*SchemaMeta {
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

{{/* {{- if eq .Kind "OBJECT_MAP" }}
{
    {{ template "required" . }}
    {{ template "objectMap" . }}
},
{{- end }} */}}

{{- if eq .Kind "SINGULAR_MESSAGE" }}
{
    name: {{.Name|quote}},
    nested: {{ template "fieldsSchema" .Message.Fields }},
},
{{- end }}

{{- if eq .Kind "SINGULAR_ELEMENTARY" }}
{
    {{ template "singularElementary" . }}
},
{{- end }}

{{/* {{- if eq .Kind "CUSTOM_TYPE" }}
Schema{{.CustomTypeMethodInfix}}(),
{{- end  */}}
{{- end -}}

{{- define "singularElementary" -}}
name: {{.Name|quote}},
isTime: {{.IsTime}},
isDuration: {{.IsDuration}},
{{- end -}}

{{- define "repeatedMessage" -}}
name: {{.Name|quote}},
nested: {{ template "fieldsSchema" .Message.Fields }},
{{- end -}}

{{- define "repeatedElementary" -}}
name: {{.Name|quote}},
isTime: {{.IsTime}},
isDuration: {{.IsDuration}},
{{- end -}}

{{- define "map" -}}
name: {{.Name|quote}},
isTime: {{.IsTime}},
isDuration: {{.IsDuration}},
{{- end -}}

{{- define "objectMap" -}}
Type: schema.TypeList,
Description: {{ .RawComment | quote }},
{{ template "configMode" . }}
Elem: &schema.Resource {
    Schema: map[string]*schema.Schema{
        "key": {
            Type: schema.TypeString,
            Required: true,
        },
        "value": {{ template "fieldSchema" .MapValueField }}
    },
},
{{- end -}}
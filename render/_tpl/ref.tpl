/*

## {{.NameSnake}}

{{.RawComment}}

{{ template "fieldsRef" .Fields -}}
{{ template "fieldsNestedRef" .Fields -}}

{{- define "fieldsRef" -}}
| Name                    | Type       | Description                               |
|-------------------------|------------|-------------------------------------------|
{{- range $index, $field := . }}
| `{{.NameSnake}}`{{if .IsRequired}}*{{end}} | {{template "refType" .}} | {{.RawComment}} |
{{- end }}
{{- end -}}

{{- define "fieldsNestedRef" -}}
{{- range $index, $field := . }}
{{- if eq .Kind "SINGULAR_MESSAGE" "REPEATED_MESSAGE" "MESSAGE_MAP" }}

**{{.NameSnake}}**

{{.RawComment}}

{{ if eq .Kind "SINGULAR_MESSAGE" "REPEATED_MESSAGE" -}}
{{ template "fieldsRef" .Message.Fields }}
{{ template "fieldsNestedRef" .Message.Fields }}
{{ else -}}
{{ template "fieldsRef" .MapValueField.Fields }}
{{ template "fieldsNestedRef" .MapValueField.Fields }}
{{- end }}
{{- end }}
{{- end }}
{{- end}}

{{- define "refType" -}}
{{- if eq .Kind "SINGULAR_ELEMENTARY" -}}
{{.SchemaRawType }}
{{- else if eq .Kind "REPEATED_ELEMENTARY" -}}
{{.SchemaRawType }} list
{{- else if eq .Kind "SINGULAR_MESSAGE" -}}
object
{{- else if eq .Kind "REPEATED_MESSAGE" -}}
object list
{{- else if eq .Kind "MAP" -}}
map
{{- else if eq .Kind "MESSSAGE_MAP" -}}
set
{{- end -}}
{{- end -}}


*/
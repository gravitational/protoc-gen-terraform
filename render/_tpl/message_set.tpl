func Set{{.Name}}ToResourceData(d *schema.ResourceData, t *{{.GoTypeName}}) error {
    p := ""

    {{ template "fields" .Fields }}

    return nil
}

{{- define "fields" -}}
{{ range $index, $field := . }}
{{- template "field" $field }}
{{ end }}
{{- end -}}

{{- define "field" -}}
{{- if eq .Kind "SINGULAR_ELEMENTARY" -}}
{
    {{ template "singularElementary" . }}
}
{{- end -}}

{{- if eq .Kind "REPEATED_ELEMENTARY" -}}
{
    {{ template "repeatedElementary" . }}
}
{{- end -}}
{{- end -}}

{{/* Renders setter for singular value of any type */}}
{{- define "singularElementary" -}}
_v := t.{{.Name}}

{{- if .GoTypeIsPtr }}
if _v != nil {
{{- end }}

{{ template "rawToValue" . }}
err := d.Set(p+{{.NameSnake | quote }}, _value)
if err != nil {
    return err
}
{{- if .GoTypeIsPtr }}
}
{{- end }}
{{- end -}}

{{/* Renders setter for elementary array of any type */}}
{{- define "repeatedElementary" -}}
_arr := t.{{.Name}}
_raw := make([]{{.SchemaRawType}}, len(_arr))

for i, _v := range _arr {
    {{- template "rawToValue" . }}
    _raw[i] = _value
}

d.Set(p+{{ .NameSnake | quote }}, _raw)
{{- end -}}

{{/* Converts elementary value from from target struct type to raw data type */}}
{{/* Input: _raw */}}
{{/* Output: _value */}}
{{- define "rawToValue" }}
{{- if .IsTime }}
_value := _v.Format(time.RFC3339)
{{- else if .IsDuration }}
_value := _v.String()
{{- else }}
_value := {{.SchemaRawType}}({{if .GoTypeIsPtr}}*{{end}}_v)
{{- end }}
{{- end -}}

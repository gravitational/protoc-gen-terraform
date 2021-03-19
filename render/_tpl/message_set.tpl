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
{{- end -}}

{{/* Renders unmarshaller for singular value of any type */}}
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

{{/* Converts elementary value from from target struct type to raw data type */}}
{{/* Input: _raw */}}
{{/* Output: _value */}}
{{- define "rawToValue" }}
{{- if .IsTime }}
_value := _v.Format(time.RFC3339)
{{- else if .IsDuration }}
_value := _v.String()
{{- else }}
_value := {{.SchemaRawType}}({{.SchemaGoType}}(_v))
{{- end }}
{{- end -}}

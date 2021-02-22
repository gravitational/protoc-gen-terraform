func Unmarshal{{.Name}}(d *schema.ResourceData, t *{{.GoTypeName}}) error {
    p := ""

    {{- template "fields" . }}

    return nil
}

{{- define "fields" -}}
{{ range $index, $field := .Fields }}
{
    {{- template "field" $field }}
}
{{ end }}
{{- end -}}

{{- define "field" -}}
{{- if eq .Kind "SINGULAR_ELEMENTARY" -}}
    {{ template "singularElementary" . }}
{{- end -}}
{{- end -}}

{{/* Renders unmarshaller for singular value of any type */}}
{{- define "singularElementary" -}}
{{- template "getOk" . }}
if ok {
    {{- template "rawToValue" . }}
    t.{{.Name}} = {{if .GoTypeIsPtr }}&{{end}}_value
}
{{- end -}}

{{/* Converts elementary value from raw form to target struct type */}}
{{/* Input: _raw */}}
{{/* Output: _value */}}
{{- define "rawToValue" -}}
{{- if .IsTime }}
_value, err := time.Parse(time.RFC3339, _raw.({{.SchemaRawType}}))
if err != nil {
    return fmt.Errorf("Malformed time value for field {{.Name}} : %w", err)
}
{{- else if .IsDuration }}
_valued, err := time.ParseDuration(_raw.({{.SchemaRawType}}))
if err != nil {
    return fmt.Errorf("Malformed duration value for field {{.Name}} : %w", err)
}
_value := {{.GoType}}(_valued)
{{- else }}
{{- if eq .SchemaRawType .SchemaGoType }}
_value := _raw.({{.SchemaRawType}})
{{- else }}
_value := {{.GoType}}({{.SchemaGoType}}(_raw.({{.SchemaRawType}})))
{{- end }}
{{- end }}
{{- end -}}

{{/* Generates schema getter statement */}}
{{/* Input: p */}}
{{/* Output: _raw */}}
{{- define "getOk" -}}
{{- if eq .SchemaRawType "bool" }}
_raw, ok := d.GetOkExists(p + {{ .NameSnake | quote }})
{{- else }}
_raw, ok := d.GetOk(p + {{ .NameSnake | quote  }})
{{- end }}
{{- end }}

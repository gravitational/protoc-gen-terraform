func Unmarshal{{.Name}}(d *schema.ResourceData, t *{{.GoTypeName}}) error {
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

{{- if eq .Kind "CUSTOM_TYPE" -}}
{
    {{ template "custom" . }}
}
{{- end -}}

{{- if eq .Kind "SINGULAR_MESSAGE" -}}
{
    {{ template "singularMessage" . }}
}
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

{{/* Renders unmarshaller for elementary array of any type */}}
{{- define "repeatedElementary" -}}
_rawi, ok := d.GetOk(p + {{ .NameSnake | quote }})
if ok {
    _rawi := _rawi.([]interface{})
    t.{{.Name}} = make({{.RawGoType}}, len(_rawi))
    for i := 0; i < len(_rawi); i++ {
        _raw := _rawi[i]
        {{- template "rawToValue" . }}
        t.{{.Name}}[i] = {{if .GoTypeIsPtr }}&{{end}}_value
    }
}
{{- end -}}

{{/* Renders custom unmarshaller custom type */}}
{{- define "custom" -}}
err := Unmarshal{{.CustomTypeMethodInfix}}(p + {{.NameSnake | quote}}, d, &t.{{.Name}})
if err != nil {
    return err
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

{{/* Generates schema getter statement with an exception for bool, which must not be parsed if not set */}}
{{/* Input: p */}}
{{/* Output: _raw */}}
{{- define "getOk" -}}
{{- if eq .SchemaRawType "bool" }}
_raw, ok := d.GetOkExists(p + {{ .NameSnake | quote }})
{{- else }}
_raw, ok := d.GetOk(p + {{ .NameSnake | quote  }})
{{- end }}
{{- end }}

{{/* Singular message */}}
{{- define "singularMessage" -}}
p := p + {{.NameSnake | quote }} + ".0."

{{ if .GoTypeIsPtr }}
_obj := {{.GoType}}{}
t.{{ .Name }} = &_obj
t := &_obj
{{ else }}
t := &t.{{.Name}}
{{ end }}

{{ template "fields" .Message.Fields }}
{{- end -}}

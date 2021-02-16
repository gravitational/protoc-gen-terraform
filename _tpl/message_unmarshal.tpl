{{/* ---- Unmarshalling ------------------------------------------------------------------*/}}
{{/* Made KISS as possible, for the price of DRY */}}

// Type full name: {{ .Name }}
func Unmarshal{{ .Name }}(d *schema.ResourceData, t *{{ .GoTypeName }}, p string) error {
    {{ template "fieldsUnmarshal" .Fields }}
    return nil
}

{{- define "fieldsUnmarshal" -}}
{{- range $index, $field := . }}
    {{ template "fieldUnmarshal" $field }}
{{- end }}
{{- end -}}

{{- define "fieldUnmarshal" -}}

// schema["{{ .NameSnake }}"] => {{ .Name }}, {{ .RawGoType }}, {{ .GoType }}
// {{ .Kind }}
{
{{- if eq .Kind "REPEATED_MESSAGE" }}
    {{ template "repeatedMessage" . }}
{{- else if eq .Kind "REPEATED_ELEMENTARY" }}
    {{ template "repeatedElementary" . }}
{{- else if eq .Kind "SINGULAR_MESSAGE" }}
    {{ template "singularMessage" . }}
{{- else if eq .Kind "SINGULAR_ELEMENTARY" }}
    {{ template "singularElementary" . }}
{{- end }}
}
{{- end -}}

{{- define "rawToValue" -}}
{{- if .IsTime }}
_value, err := time.Parse(time.RFC3339, _raw.({{.TFSchemaRawType}}))
if err != nil {
    return fmt.Errorf("Malformed time value for field {{.Name}} : %w", err)
}
{{- else if .IsDuration }}
_value, err := time.ParseDuration(_raw.({{.TFSchemaRawType}}))
if err != nil {
    return fmt.Errorf("Malformed duration value for field {{.Name}} : %w", err)
}
{{- else }}
_value := {{.TFSchemaGoType}}(_raw.({{.TFSchemaRawType}}))
{{- end }}
{{- end -}}

{{- define "repeatedMessage" -}}
p := p + "{{.NameSnake}}"
_rawi, ok := d.GetOk(p)
if ok {
    _rawi := _rawi.([]interface{})
    t.{{.Name}} = make([]{{.GoType}}, len(_rawi))
    for i := 0; i < len(_rawi); i++ {
        t := &t.{{ .Name }}[i]
        p := p + fmt.Sprintf(".%v.", i)
        {{- template "fieldsUnmarshal" .Message.Fields }}            
    }
}
{{- end -}}

{{- define "repeatedElementary" -}}
_rawi, ok := d.GetOk(p + "{{ .NameSnake}}")
if ok {
    _rawi := _rawi.([]interface{})
    t.{{.Name}} = make([]{{.GoType}}, len(_rawi))
    for i := 0; i < len(_rawi); i++ {
        _raw := _rawi[i]
        {{- template "rawToValue" . }}
        _final := {{.GoType}}(_value)
        t.{{.Name}}[i] = {{if .GoTypeIsPtr }}&{{end}}_final            
    }
}
{{- end -}}

{{- define "singularMessage" -}}
p := p + "{{ .NameSnake }}.0."
{{ if .GoTypeIsPtr }}
_obj := {{ .GoType}}{}
t.{{ .Name }} = &_obj
t := &_obj
{{- else -}}
t := &t.{{ .Name }}
{{ end }}
{{- template "fieldsUnmarshal" .Message.Fields }}
p = p
t = t
{{- end -}}

{{- define "singularElementary" -}}
_raw, ok := d.GetOk(p + "{{ .NameSnake}}")
if ok {
    {{- template "rawToValue" . }}
    _final := {{if .GoTypeIsSlice }}[]{{end}}{{.GoType}}(_value)
    t.{{.Name}} = {{if .GoTypeIsPtr }}&{{end}}_final
}
{{- end -}}
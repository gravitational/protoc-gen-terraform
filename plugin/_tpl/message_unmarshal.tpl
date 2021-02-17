{{/* ---- Unmarshalling ------------------------------------------------------------------*/}}
{{/* Made KISS as possible, for the exchange of DRY */}}

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
{{- else if eq .Kind "SINGULAR_MESSAGE_FOLD" }}
    {{ template "singularMessageFold" . }}
{{- else if eq .Kind "SINGULAR_ELEMENTARY" }}
    {{ template "singularElementary" . }}
{{- end }}
}
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
p := p + "{{.NameSnake}}"
_rawi, ok := d.GetOk(p)
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
p := p + "{{.NameSnake}}.0."

{{ if .GoTypeIsPtr }}
_obj := {{.GoType}}{}
t.{{ .Name }} = &_obj
t := &_obj
{{- else -}}
t := &t.{{.Name}}
{{ end }}

{{ template "fieldsUnmarshal" .Message.Fields }}
// NOTE: remove
p = p
t = t
{{- end -}}

{{- define "singularMessageFold" -}}
{{ $folded := .Message.Fields | first }}
p := p + "{{.NameSnake}}"

{{ template "getOk" $folded }}
if ok {
    {{ if .GoTypeIsPtr }}
    _obj := {{.GoType}}{}
    t.{{ .Name }} = {{ if.GoTypeIsPtr}}&{{end}}_obj
    t := &_obj
    {{ else }}
    t := &t.{{.Name}}
    {{ end }}

    {{ template "rawToValue" $folded }}
    {{ template "assignSingularElementary" $folded }}
}
// NOTE: remove
p = p
t = t
{{- end -}}

{{- define "singularElementary" -}}
p := p + "{{.NameSnake}}"

{{ template "getOk" . }}
if ok {
    {{ template "rawToValue" . }}
    {{ template "assignSingularElementary" . }}
}
{{- end -}}

{{- define "rawToValue" -}}
{{- if .IsTime }}
_value, err := time.Parse(time.RFC3339, _raw.({{.SchemaRawType}}))
if err != nil {
    return fmt.Errorf("Malformed time value for field {{.Name}} : %w", err)
}
{{- else if .IsDuration }}
_value, err := time.ParseDuration(_raw.({{.SchemaRawType}}))
if err != nil {
    return fmt.Errorf("Malformed duration value for field {{.Name}} : %w", err)
}
{{- else }}
_value := {{.SchemaGoType}}(_raw.({{.SchemaRawType}}))
{{- end }}
{{- end -}}

{{- define "assignSingularElementary" -}}
_final := {{if .GoTypeIsSlice }}[]{{end}}{{.GoType}}(_value)
t.{{.Name}} = {{if .GoTypeIsPtr }}&{{end}}_final
{{- end -}}

{{- define "getOk" -}}
{{- if eq .SchemaRawType "bool" }}
_raw, ok := d.GetOkExists(p)
{{- else }}
_raw, ok := d.GetOk(p)
{{- end }}
{{- end }}
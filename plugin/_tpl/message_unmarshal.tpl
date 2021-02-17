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
// {{ .Kind }}, Map: {{ .IsMap }}, List: {{ .IsRepeated }}, Container: {{ .IsContainer }}
{
{{- if eq .Kind "REPEATED_MESSAGE" }}
    {{ template "repeatedMessage" . }}
{{- else if eq .Kind "REPEATED_ELEMENTARY" }}
    {{ template "repeatedElementary" . }}
{{- else if eq .Kind "SINGULAR_MESSAGE" }}
{{- if .IsContainer }}
    {{ template "container" . }}
{{- else }}
    {{ template "singularMessage" . }}
{{- end }}
{{- else if eq .Kind "SINGULAR_ELEMENTARY" }}
    {{ template "singularElementary" . }}
{{- else if eq .Kind "MAP" }}
    {{ template "map" . }}
{{- else if eq .Kind "ARTIFICIAL_OBJECT_MAP" }}
    t = t
    p = p
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

// TODO: Check if section exists to skip pointer initialization

{{ template "initMessage" . }}
{{ template "fieldsUnmarshal" .Message.Fields }}
{{- end -}}

{{- define "container" -}}
{{ $folded := .Message.Fields | first }}
p := p + "{{.NameSnake}}"
{{ template "initMessage" . }}
{{ template "fieldUnmarshal" $folded }}
{{- end -}}

{{- define "singularElementary" -}}
p := p + "{{.NameSnake}}"

{{ template "getOk" . }}
if ok {
    {{ template "rawToValue" . }}
    {{ template "assignSingularElementary" . }}
}
{{- end -}}

{{- define "map" -}}
{{ $m := .MapValueField }}
p := p + "{{.NameSnake}}"
_rawm, ok := d.GetOk(p)
if ok {
    _rawm := _rawm.(map[string]interface{})
    t.{{.Name}} = make(map[string]{{$m.GoType}}, len(_rawm))
    for _k, _v := range _rawm {
        _raw := _v
        {{- template "rawToValue" $m }}
        _final := {{if $m.GoTypeIsSlice }}[]{{end}}{{$m.GoType}}(_value)
        t.{{.Name}}[_k] = {{if $m.GoTypeIsPtr }}&{{end}}_final
    }   
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

{{- define "initMessage" -}}
{{ if .GoTypeIsPtr }}
_obj := {{.GoType}}{}
t.{{ .Name }} = &_obj
t := &_obj
{{ else }}
t := &t.{{.Name}}
{{ end }}
{{- end }}
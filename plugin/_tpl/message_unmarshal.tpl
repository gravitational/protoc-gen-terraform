{{/* ---- Unmarshalling ------------------------------------------------------------------*/}}

// Type: {{ .GoTypeName }}
func Unmarshal{{ .Name }}(d *schema.ResourceData, t *{{ .GoTypeName }}) error {
    var p string
    {{ template "fieldsUnmarshal" .Fields }}
    return nil
}

{{- define "fieldsUnmarshal" -}}
{{- range $index, $field := . }}
    {{ template "fieldUnmarshal" $field }}
{{- end }}
{{- end -}}

{{- define "fieldUnmarshal" -}}
{
    {{ if eq .Kind "SINGULAR_MESSAGE" }}
    p := p + "{{.NameSnake}}{{ if not .IsContainer }}.0.{{ end }}"
    {{ else }}
    p := p + "{{.NameSnake}}"
    {{ end }}
    {{ template "fieldUnmarshalBody" . }}
}
{{- end -}}

{{- define "fieldUnmarshalBody" -}}
{{- if eq .Kind "SINGULAR_ELEMENTARY" }}
// SINGULAR_ELEMENTARY {{ .Name }}
{{ template "singularElementary" . }}
{{- else if eq .Kind "REPEATED_ELEMENTARY" }}
// REPEATED_ELEMENTARY {{ .Name }}
{{ template "repeatedElementary" . }}
{{- else if eq .Kind "REPEATED_MESSAGE" }}
// REPEATED_MESSAGE {{ .Name }}
{{ template "repeatedMessage" . }}
{{- else if eq .Kind "SINGULAR_MESSAGE" }}
// SINGULAR_MESSAGE {{ .Name }}
{{ if .IsContainer  }}
// NORMAL CONTAINER
{{ template "singularContainerMessage" . }}
{{ else }}
{{ template "singularMessage" . }}
{{ end }}
{{- else if eq .Kind "MAP" }}
// MAP {{ .Name }}
{{ template "map" . }}
{{- else if eq .Kind "ARTIFICIAL_OBJECT_MAP" }}
// ARTIFICIAL_OBJECT_MAP {{ .Name }}
{{ template "artificialObjectMap" . }}
{{- end }}
{{- end -}}

{{- define "singularElementary" -}}
{{ template "getOk" . }}
if ok {
    {{ template "rawToValue" . }}
    t.{{.Name}} = {{if .GoTypeIsPtr }}&{{end}}_value
}
{{- end -}}

{{- define "repeatedElementary" -}}
_rawi, ok := d.GetOk(p)
if ok {
    _rawi := _rawi.([]interface{})
    t.{{.Name}} = make([]{{.GoType}}, len(_rawi))
    for i := 0; i < len(_rawi); i++ {
        _raw := _rawi[i]
        {{- template "rawToValue" . }}
        t.{{.Name}}[i] = {{if .GoTypeIsPtr }}&{{end}}_value
    }
}
{{- end -}}

{{- define "repeatedMessage" -}}
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

{{- define "singularMessage" -}}
{{ template "initTarget" . }}
{{ template "fieldsUnmarshal" .Message.Fields }}
{{- end -}}

{{- define "singularContainerMessage" -}}
{{ $folded := .Message.Fields | first }}
{{ template "initTarget" . }}
{{ template "fieldUnmarshalBody" $folded }}
{{- end -}}

{{- define "map" -}}
{{ $m := .MapValueField }}
_rawm, ok := d.GetOk(p)
if ok {
    _rawm := _rawm.(map[string]interface{})
    t.{{.Name}} = make(map[string]{{$m.GoType}}, len(_rawm))
    for _k, _v := range _rawm {
        _raw := _v
        {{- template "rawToValue" $m }}
        t.{{.Name}}[_k] = {{if $m.GoTypeIsPtr }}&{{end}}_value
    }   
}
{{- end -}}

{{- define "artificialObjectMap" -}}
{{ $m := .MapValueField }}
// value: {{ .MapValueField.GoType }} {{ .MapValueField.IsContainer }}
_rawi, ok := d.GetOk(p)
if ok {
    _rawi := _rawi.([]interface{})
    for _, _artificialItem := range _rawi {
        _item := _artificialItem.(map[string]interface{})

        p := p + "." + _item["key"].(string)

        {{ template "fieldUnmarshalBody" .MapValueField }}
    }
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
_value, err := time.ParseDuration(_raw.({{.SchemaRawType}}))
if err != nil {
    return fmt.Errorf("Malformed duration value for field {{.Name}} : %w", err)
}
{{- else }}
_value := {{.GoType}}({{.SchemaGoType}}(_raw.({{.SchemaRawType}})))
{{- end }}
{{- end -}}

{{/* Generates schema getter statement */}}
{{/* Input: p */}}
{{/* Output: _raw */}}
{{- define "getOk" -}}
{{- if eq .SchemaRawType "bool" }}
_raw, ok := d.GetOkExists(p)
{{- else }}
_raw, ok := d.GetOk(p)
{{- end }}
{{- end }}

{{/* Initializes target struct field if needed */}}
{{/* Output: new t */}}
{{- define "initTarget" -}}
{{ if .GoTypeIsPtr }}
_obj := {{.GoType}}{}
t.{{ .Name }} = &_obj
t := &_obj
{{ else }}
t := &t.{{.Name}}
{{ end }}
{{- end }}
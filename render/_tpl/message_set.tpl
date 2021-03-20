func Set{{.Name}}ToResourceData(d *schema.ResourceData, t *{{.GoTypeName}}) error {
    obj := make(map[string]interface{})

    {{ template "fields" .Fields }}

    for key, value := range obj {
        err := d.Set(key, value)
        if err != nil {
            return err
        }
    }

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

{{- if eq .Kind "REPEATED_MESSAGE" -}}
{
    {{ template "repeatedMessage" . }}
}
{{- end -}}

{{- if eq .Kind "MAP" -}}
{
    {{ template "map" . }}
}
{{- end -}}

{{- if eq .Kind "OBJECT_MAP" -}}
{
    {{ template "objectMap" . }}
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
obj[{{.NameSnake | quote}}] = _value
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

obj[{{.NameSnake | quote}}] = _raw
{{- end -}}

{{/* Renders custom getter custom type */}}
{{- define "custom" -}}
_v, err := Set{{.CustomTypeMethodInfix}}ToResourceData(&t.{{.Name}})
if err != nil {
    return err
}
obj[{{.NameSnake | quote}}] = _v
{{- end -}}

{{/* Converts elementary value from from target struct type to raw data type */}}
{{/* Input: _raw */}}
{{/* Output: _value */}}
{{- define "rawToValue" }}
{{- if .IsTime }}
_value := _v.Format(time.RFC3339)
{{- else if .IsDuration }}
_value := time.Duration(_v).String()
{{- else }}
_value := {{.SchemaRawType}}({{if .GoTypeIsPtr}}*{{end}}_v)
{{- end }}
{{- end -}}

{{/* Singular message */}}
{{- define "singularMessage" -}}
msg := make(map[string]interface{})
obj[{{.NameSnake | quote }}] = []interface{}{msg}
{
    obj := msg
    t := t.{{.Name}}

    {{ template "fields" .Message.Fields }}
}
{{- end -}}

{{/* Repeated message */}}
{{- define "repeatedMessage" -}}
arr := make([]interface{}, len(t.{{.Name}}))

for i, t := range t.{{.Name}} {
    obj := make(map[string]interface{})
    {{ template "fields" .Message.Fields }}
    arr[i] = obj
}

if len(arr) > 0 {
    obj[{{.NameSnake | quote }}] = arr
}
{{- end -}}

{{/* String -> elementary value map */}}
{{- define "map" -}}
{{ $m := .MapValueField }}
m := make(map[string]interface{})

for key, _v := range t.{{.Name}} {
    {{- template "rawToValue" $m }}
    m[key] = _value
}

if len(m) > 0 {
    obj[{{.NameSnake | quote}}] = m
}
{{- end -}}

{{/* String -> object map */}}
{{- define "objectMap" -}}
{{ $m := .MapValueField }}
a := make([]interface{}, len(t.{{.Name}}))
n := 0

for k, v := range t.{{.Name}} {
    i := make(map[string]interface{})
    i["key"] = k
    
    obj := make(map[string]interface{})
    t := v
    {{ template "fields" $m.Message.Fields }}
    i["value"] = []interface{}{obj}

    a[n] = i
    n++
}

if len(a) > 0 {
    obj[{{.NameSnake | quote}}] = a
}
{{- end -}}

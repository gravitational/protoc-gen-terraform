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
obj[{{.NameSnake | quote}}] = Set{{.CustomTypeMethodInfix}}ToResourceData(&t.{{.Name}})
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
arr := make([]map[string]interface{}, len(t.{{.Name}}))

for i, t := range t.{{.Name}} {
    arr[i] = make(map[string]interface{})
    obj := arr[i]
    {{ template "fields" .Message.Fields }}
}

obj[{{.NameSnake | quote }}] = arr
{{- end -}}

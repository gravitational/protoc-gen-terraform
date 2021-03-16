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
    _rawi, ok := _rawi.([]interface{})
    if !ok {
        return fmt.Errorf("count not convert %T to []interface{}", _rawi)
    }
    t.{{.Name}} = make({{.GoTypeFull}}, len(_rawi))
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
{{- define "rawToValue" }}
_raws, ok := _raw.({{.SchemaRawType}})
if !ok {
    return fmt.Errorf("can not convert %T to {{.SchemaRawType}}", _raws)
}

{{- if .IsTime }}
_value, err := time.Parse(time.RFC3339, _raws)
if err != nil {
    return fmt.Errorf("malformed time value for field {{.Name}} : %w", err)
}
{{- else if .IsDuration }}
_valued, err := time.ParseDuration(_raws)
if err != nil {
    return fmt.Errorf("malformed duration value for field {{.Name}} : %w", err)
}
_value := {{.GoType}}(_valued)
{{- else }}
_value := {{.GoType}}({{.SchemaGoType}}(_raws))
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

{{/* Repeated message */}}
{{- define "repeatedMessage" -}}
p := p + {{.NameSnake | quote }}

_rawi, ok := d.GetOk(p)
if ok {
    _rawi := _rawi.([]interface{})
    t.{{.Name}} = make({{.GoTypeFull}}, len(_rawi))
    for i := 0; i < len(_rawi); i++ {
        {{ if .GoTypeIsPtr }}
        _obj := {{.GoType}}{}
        t.{{.Name }}[i] = &_obj
        {{ end }}

        {
            t := {{if not .GoTypeIsPtr}}&{{end}}t.{{.Name }}[i]
            p := p + fmt.Sprintf(".%v.", i)
            {{ template "fields" .Message.Fields }}
        }
    }
}
{{- end -}}

{{/* String -> elementary value map */}}
{{- define "map" -}}
{{ $m := .MapValueField }}
p := p + {{.NameSnake | quote }}
_rawm, ok := d.GetOk(p)
if ok {
    _rawmi, ok := _rawm.(map[string]interface{})
    if !ok {
        return fmt.Errorf("can not convert %T to map[string]interface{}", _rawm)
    }
    t.{{.Name}} = make(map[string]{{$m.GoType}}, len(_rawmi))
    for _k, _v := range _rawmi {
        _raw := _v
        {{- template "rawToValue" $m }}
        t.{{.Name}}[_k] = {{if $m.GoTypeIsPtr }}&{{end}}_value
    }   
}
{{- end -}}

{{/* String -> object map */}}
{{- define "objectMap" -}}
p := p + {{.NameSnake | quote }}

{{ $m := .MapValueField }}
_rawi, ok := d.GetOk(p)
if ok {
    _rawi := _rawi.([]interface{})
    _value := make(map[string]{{$m.GoTypeFull}})

    for i, _ := range _rawi {
        _rawkey := d.Get(fmt.Sprintf("%v.%v.", p, i) + "key")
        _key, ok := _rawkey.(string)
        if !ok {
            return fmt.Errorf("can not convert %T to string", _rawkey)
        }
        if _key == "" {
            return fmt.Errorf("missing key field in object map {{.Name}}")
        }

        {{ if $m.GoTypeIsPtr }}
        _obj := {{$m.GoType}}{}
        _value[_key] = &_obj
        t := &_obj
        {{ else }}
        t := &_value[_key]
        {{ end }}

        {
            p := fmt.Sprintf("%v.%v.value.0.", p, i)
            {{ template "fields" $m.Message.Fields }}
        }
    }

    t.{{.Name}} = _value
}
{{- end -}}

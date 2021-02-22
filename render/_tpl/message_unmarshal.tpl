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

{{/* Repeated message */}}
{{- define "repeatedMessage" -}}
p := p + {{.NameSnake | quote }}

_rawi, ok := d.GetOk(p)
if ok {
    _rawi := _rawi.([]interface{})
    t.{{.Name}} = make({{.RawGoType}}, len(_rawi))
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
    _rawm := _rawm.(map[string]interface{})
    t.{{.Name}} = make(map[string]{{$m.GoType}}, len(_rawm))
    for _k, _v := range _rawm {
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
    _value := make(map[string]{{$m.RawGoType}})

    for i, _ := range _rawi {
        key := d.Get(fmt.Sprintf("%v.%v.", p, i) + "key").(string)
        
        if key == "" {
            return fmt.Errorf("Missing key field in object map {{.Name}}")
        }

        {{ if $m.GoTypeIsPtr }}
        _obj := {{$m.GoType}}{}
        _value[key] = &_obj
        t := &_obj
        {{ else }}
        t := &_value[key]
        {{ end }}

        {
            p := fmt.Sprintf("%v.%v.value.0.", p, i)
            {{ template "fields" $m.Message.Fields }}
        }
    }

    t.{{.Name}} = _value
}
{{- end -}}

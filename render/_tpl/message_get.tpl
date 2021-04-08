func Get{{.Name}}FromResourceData(d *schema.ResourceData, t *{{.GoTypeName}}) error {
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

{{- if eq .Kind "MESSSAGE_MAP" -}}
{
    {{ template "messageMap" . }}
}
{{- end -}}
{{- end -}}

{{/* Renders unmarshaller for singular value of any type */}}
{{- define "singularElementary" -}}
{{- if eq .SchemaRawType "bool" }}
_raw, ok := d.GetOkExists(p + {{ .NameSnake | quote }})
{{- else }}
_raw, ok := d.GetOk(p + {{ .NameSnake | quote  }})
{{- end }}

if ok {
    {{- template "rawToValue" . }}
    t.{{.Name}} = {{if .GoTypeIsPtr }}&{{end}}_value
}
{{- end -}}

{{/* Renders unmarshaller for elementary array of any type */}}
{{- define "repeatedElementary" -}}
_a, ok := d.GetOk(p + {{ .NameSnake | quote }})
if ok {
    a, ok := _a.([]interface{})
    if !ok {
        return fmt.Errorf("count not convert %T to []interface{}", _a)
    }
    if len(a) > 0 {
        t.{{.Name}} = make({{.GoTypeFull}}, len(a))
        for i := 0; i < len(a); i++ {
            _raw := a[i]
            {{- template "rawToValue" . }}
            t.{{.Name}}[i] = {{if .GoTypeIsPtr }}&{{end}}_value
        }
    }
} else {
    t.{{.Name}} = make({{.GoTypeFull}}, 0)
}
{{- end -}}

{{/* Renders custom unmarshaller custom type */}}
{{- define "custom" -}}
err := Get{{.CustomTypeMethodInfix}}FromResourceData(p + {{.NameSnake | quote}}, d, &t.{{.Name}})
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
_value, err := time.Parse(time.RFC3339Nano, _raws)
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

{{/* Singular message */}}
{{- define "singularMessage" -}}
n := d.Get(p + {{.NameSnake | quote }} + ".#")

p := p + {{.NameSnake | quote }} + ".0"
_, ok := d.GetOk(p)
if ok && n != nil && n.(int) != 0 {
    p := p + "."

    {{ if .GoTypeIsPtr }}
    _obj := {{.GoType}}{}
    t.{{ .Name }} = &_obj
    t := &_obj
    {{ else }}
    t := &t.{{.Name}}
    {{ end }}

    {{ template "fields" .Message.Fields }}
} else {
    {{ if .GoTypeIsPtr }}
    t.{{ .Name }} = nil
    {{ end }}
}
{{- end -}}

{{/* Repeated message */}}
{{- define "repeatedMessage" -}}
p := p + {{.NameSnake | quote }}

_a, ok := d.GetOk(p)
if ok {
    a, ok := _a.([]interface{})
    if !ok {
        return fmt.Errorf("can not convert %T to []interface{}", _a)
    }

    if len(a) > 0 {
        t.{{.Name}} = make({{.GoTypeFull}}, len(a))
     
        for i := 0; i < len(a); i++ {
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
} else {
    t.{{.Name}} = make({{.GoTypeFull}}, 0)
}
{{- end -}}

{{/* String -> elementary value map */}}
{{- define "map" -}}
{{ $m := .MapValueField }}
p := p + {{.NameSnake | quote }}
_m, ok := d.GetOk(p)
if ok {
    m, ok := _m.(map[string]interface{})
    if !ok {
        return fmt.Errorf("can not convert %T to map[string]interface{}", _m)
    }
    if len(m) > 0 {
        t.{{.Name}} = make(map[string]{{$m.GoTypeFull}}, len(m))
        for _k, _v := range m {
            _raw := _v
            {{- template "rawToValue" $m }}
            t.{{.Name}}[_k] = {{if $m.GoTypeIsPtr }}&{{end}}_value
        }   
    }
} else {
    t.{{.Name}} = make(map[string]{{$m.GoTypeFull}})
}
{{- end -}}

{{/* String -> object map */}}
{{- define "messageMap" -}}
p := p + {{.NameSnake | quote }}

{{ $m := .MapValueField }}
_m, ok := d.GetOk(p)
if ok {
    m, ok := _m.([]interface{})
    if !ok {
        return fmt.Errorf("can not convert %T to []interface{}", _m)
    }

    if len(m) > 0 {
        _value := make(map[string]{{$m.GoTypeFull}})

        for i := range m {
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
} else {
    t.{{.Name}} = make(map[string]{{$m.GoTypeFull}})
}
{{- end -}}

// Type full name: {{ .Name }}
func Schema{{ .Name }}() map[string]*schema.Schema {
	return {{ template "fieldsSchema" .Fields -}}
}

func Unmarshal{{ .Name }}(d *schema.ResourceData, t *{{ .GoTypeName }}, p string) error {
    {{ template "fieldsUnmarshal" .Fields }}
    return nil
}

{{/* ---- Schema rendering ---------------------------------------------------------------*/}}

{{- define "fieldsSchema" -}}
map[string]*schema.Schema {
{{- range $index, $field := . }}
	// {{ .Name }}
	"{{ .NameSnake }}": {{ template "fieldSchema" . }}    
{{- end }}
}
{{- end -}}

{{- define "fieldSchema" -}}
{   
	Type: schema.{{ coalesce .TFSchemaAggregateType .TFSchemaType }},

	{{- if .IsRequired }}
	Required: true,
	{{- else }}
	Optional: true,
	{{- end }}

	{{- if .TFSchemaMaxItems }}
	MaxItems: {{ .TFSchemaMaxItems }},
	{{- end }}

	{{- if .TFSchemaValidate }}
	ValidateFunc: {{ .TFSchemaValidate }},
	{{- end }}

    {{- if .IsMessage }}
    Elem: &schema.Resource {
        Schema: Schema{{ .Message.Name }}(),
    },
    {{- else if .IsAggregate }}
    Elem: &schema.Schema {
        Type: schema.{{ .TFSchemaType }},
    },
    {{- end }}
},
{{- end -}}

{{/* ---- Unmarshalling ------------------------------------------------------------------*/}}
{{/* Made KISS as possible, for the price of DRY */}}

{{- define "fieldsUnmarshal" -}}
{{- range $index, $field := . }}
    {{ template "fieldUnmarshal" $field }}
{{- end }}
{{- end -}}

{{- define "fieldUnmarshal" -}}

// schema["{{ .NameSnake }}"] => {{ .Name }}, {{ .RawGoType }}, {{ .GoType }}
{
{{- if and .IsAggregate .IsRepeated .IsMessage }}
    {{ template "repeatedMessage" . }}
{{- else if and .IsAggregate .IsRepeated }}
    {{ template "repeatedElementary" . }}
{{- else if .IsMessage }}
    {{ template "singularMessage" . }}
{{- else }}
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
// repeated message
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
// repeated elementary
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
// singular message
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
// singular elementary
_raw, ok := d.GetOk(p + "{{ .NameSnake}}")
if ok {
    {{- template "rawToValue" . }}
    _final := {{.GoType}}(_value)
    t.{{.Name}} = {{if .GoTypeIsPtr }}&{{end}}_final
}
{{- end -}}
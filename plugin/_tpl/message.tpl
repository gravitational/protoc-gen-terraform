// Type full name: {{ .Name }}
func Schema{{ .Name }}() map[string]*schema.Schema {
	return {{ template "fieldsSchema" .Fields -}}
}

func Unmarshal{{ .Name }}(d *schema.ResourceData, t *{{ .GoTypeName }}, p string) error {
    {{- template "fieldsUnmarshal" .Fields }}

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
{{- define "fieldsUnmarshal" -}}
{{- range $index, $field := . }}
{
	// schema["{{ .NameSnake }}"] => {{ .Name }}, {{ .RawGoType }}, {{ .GoType }}
    _raw, ok := d.GetOk(p + "{{ .NameSnake}}")
    if ok {
        {{- if .IsAggregate }}
            {{- if .IsRepeated }}
                _rawi := _raw.([]interface{})
                t.{{.Name}} = make([]{{.GoType}}, len(_rawi))
                for i := 0; i < len(_rawi); i++ {
                    {{- if .IsMessage }}
                    Unmarshal{{ .Message.Name }}(d, &t.{{ .Name }}[i], p+fmt.Sprintf("{{ .NameSnake }}.%v.", i))
                    _raw = _raw
                    {{- else }}
                    _currentRaw := _rawi[i]
                    {{- template "rawToValue" dict "raw" "_currentRaw" "field" . }}
                    _tmp := {{.GoType}}(_value)
                    t.{{.Name}}[i] = {{if .GoTypeIsPtr }}&{{end}}_tmp
                    {{- end }}
                }
            {{- end }}
        {{- else -}}
            {{- if .IsMessage -}}
                Unmarshal{{ .Message.Name }}(d, &t.{{ .Name }}, "{{ .NameSnake }}.0.")
            {{- else -}}
                {{/* We convert from schema type to real type */}}
                {{ template "rawToValue" dict "raw" "_raw" "field" . }}
                _tmp := {{.GoType}}(_value)
                t.{{.Name}} = {{if .GoTypeIsPtr }}&{{end}}_tmp
            {{- end }}                
        {{- end }}
    }
}
{{- end }}
{{- end -}}

{{- define "rawToValue" -}}
{{- if .field.IsTime }}
_value, err := time.Parse(time.RFC3339, {{.raw}}.({{.field.TFSchemaRawType}}))
if err != nil {
    return fmt.Errorf("Malformed time value for field {{.field.Name}} : %w", err)
}
{{- else if .field.IsDuration }}
_value, err := time.ParseDuration({{.raw}}.({{.field.TFSchemaRawType}}))
if err != nil {
    return fmt.Errorf("Malformed duration value for field {{.field.Name}} : %w", err)
}
{{- else }}
_value := {{.field.TFSchemaGoType}}({{.raw}}.({{.field.TFSchemaRawType}}))
{{- end }}
{{- end -}}
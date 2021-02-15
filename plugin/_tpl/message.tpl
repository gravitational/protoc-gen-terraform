// Type full name: {{ .Name }}
func Schema{{ .Name }}() map[string]*schema.Schema {
	return {{ template "fieldsSchema" .Fields -}}
}

func Unmarshal{{ .Name }}(d *schema.ResourceData, t *types.{{ .Name }}, p string) error {
    {{- template "fieldsUnmarshal" .Fields }}

    return nil
}

{{/* Schema rendering */}}

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

{{/* Marshal rendering */}}
{{- define "fieldsUnmarshal" -}}
{{- range $index, $field := . }}
{
	// schema["{{ .NameSnake }}"] => {{ .Name }}, {{ .GoType }}
    _raw, ok := d.GetOk(prefix + "{{ .NameSnake}}")
    if ok {
        {{- if .IsMessage }} // and not IsAggregate
            Unmarshal{{ .Message.Name }}(r, &d.{{ .Name }}, "{{ .NameSnake }}.0.")
        {{- else if .IsAggregate }}
        {{- else -}}
            {{/* We convert from schema type to real type */}}
            {{- if .IsTime }}
            _value, ok := time.Parse(time.RFC3339, _raw.({{.TFSchemaRawType}}))
            if !ok {
                return fmt.Errorf("Malformed time value for field {{.Name}}")
            }
            {{- else if .IsDuration }}
            _value, ok := time.ParseDuration(_raw.({{.TFSchemaRawType}}))
            if !ok {
                return fmt.Errorf("Malformed duration value for field {{.Name}}")
            }
            {{- else }}
            _value := {{.TFSchemaGoType}}(_raw.({{.TFSchemaRawType}}))
            {{- end }}
            t.{{.Name}} = {{- if .IsNullable -}}&{{- end -}}{{.GoType}}(_value)
        {{- end }}
    }
}
{{- end }}
{{- end -}}

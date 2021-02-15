// Type full name: {{ .Name }}
func Schema{{ .Name }}() map[string]*schema.Schema {
	return {{ template "fieldsSchema" .Fields -}}
}

func Unmarshal{{ .Name }}(d *schema.ResourceData, t *types.{{ .Name }}, p string) error {
    {{ template "fieldsUnmarashal" .Fields -}}

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

    {{- if .IsAggregate }}
    Elem: &schema.Schema {
        Type: schema.{{ .TFSchemaType }},
    },
    {{- else if .IsMessage }}
    Elem: &schema.Resource {
        Schema: Schema{{ .Message.Name }}(),
    },
    {{- end }}
},
{{- end -}}

{{/* Marshal rendering */}}
{{- define "fieldsUnmarshal" -}}
{{- range $index, $field := . }}
	// {{ .NameSnake }}
    _{{ .NameSnake }}_raw, ok := d.GetOk(prefix + "{{ .NameSnake}}")
    if ok {
    
    }
{{- end }}
}
{{- end -}}

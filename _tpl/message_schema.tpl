{{/* ---- Schema rendering ---------------------------------------------------------------*/}}
// Type full name: {{ .Name }}
func Schema{{ .Name }}() map[string]*schema.Schema {
	return {{ template "fieldsSchema" .Fields -}}
}

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
        Schema: {{ template "fieldsSchema" .Message.Fields }},
    },
    {{- else if .IsAggregate }}
    Elem: &schema.Schema {
        Type: schema.{{ .TFSchemaType }},
    },
    {{- end }}
},
{{- end -}}

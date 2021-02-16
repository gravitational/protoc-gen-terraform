{{/* ---- Schema rendering ---------------------------------------------------------------*/}}
// Type full name: {{ .Name }}
func Schema{{ .Name }}() map[string]*schema.Schema {
	return {{ template "fieldsSchema" .Fields -}}
}

{{- define "fieldsSchema" -}}
map[string]*schema.Schema {
{{- range $index, $field := . }}
	// {{ .Name }} {{ .Kind }}
	"{{ .NameSnake }}": {{ template "fieldSchema" . }}    
{{- end }}
}
{{- end -}}

{{- define "fieldSchema" -}}
{   
    {{- if eq .Kind "REPEATED_MESSAGE" }}
    Type: schema.{{ .TFSchemaAggregateType }},
    Elem: &schema.Resource {
        Schema: {{ template "fieldsSchema" .Message.Fields }},
    },
    {{- end }}

    {{- if eq .Kind "REPEATED_ELEMENTARY" }}
    Type: schema.{{ .TFSchemaAggregateType }},
    Elem: &schema.Schema {
        Type: schema.{{ .TFSchemaType }},
    },
    {{- end }}

    {{- if eq .Kind "SINGULAR_MESSAGE" }}
    Type: schema.{{ .TFSchemaAggregateType }},
    MaxItems: 1,
    Elem: &schema.Resource {
        Schema: {{ template "fieldsSchema" .Message.Fields }},
    },
    {{- end }}

    {{- if eq .Kind "SINGULAR_MESSAGE_FOLD" }}
    {{ template "singularElementary" .Message.Fields | first }}
    {{- end }}

    {{- if eq .Kind "SINGULAR_ELEMENTARY" }}
    {{ template "singularElementary" . }}
    {{- end }}

    {{- if .IsRequired }}
    Required: true,
    {{- else }}
    Optional: true,
    {{- end }}
},
{{- end -}}

{{- define "singularElementary" -}}
Type: schema.{{ .TFSchemaType }},
{{- if .TFSchemaValidate }}
ValidateFunc: {{ .TFSchemaValidate }},
{{- end }}
{{- end -}}
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
    Type: schema.TypeList,
    Elem: &schema.Resource {
        Schema: {{ template "fieldsSchema" .Message.Fields }},
    },
    {{- end }}

    {{- if eq .Kind "REPEATED_ELEMENTARY" }}
    Type: schema.TypeList,
    Elem: &schema.Schema {
        Type: {{ template "type" .SchemaRawType }},
    },
    {{- end }}

    {{- if eq .Kind "SINGULAR_MESSAGE" }}
    Type: schema.TypeList,
    MaxItems: 1,
    Elem: &schema.Resource {
        Schema: {{ template "fieldsSchema" .Message.Fields }},
    },
    {{- end }}

    {{- if eq .Kind "ELEMENTARY_CONTAINER" }}
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
Type: {{ template "type" .SchemaRawType }},
{{- if .IsTime }}
ValidateFunc: validation.IsRFC3339Time,
{{- end }}
{{- end -}}

{{- define "type" -}}
{{- if eq . "float64" -}}
schema.TypeFloat
{{- else if eq . "int" -}}
schema.TypeInt
{{- else if eq . "bool" -}}
schema.TypeBool
{{- else if eq . "string" -}}
schema.TypeString
{{- end -}}
{{- end -}}

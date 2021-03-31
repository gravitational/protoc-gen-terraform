{{/* ---- Schema rendering ---------------------------------------------------------------*/}}
// Schema{{ .Name }} returns schema for {{.Name}}
//
{{.Comment}}
func Schema{{ .Name }}() map[string]*schema.Schema {
	return {{ template "fieldsSchema" .Fields -}}
}

{{- define "fieldsSchema" -}}
map[string]*schema.Schema {
{{- range $index, $field := . }}
    {{if .Comment}}{{.Comment}}{{else}}{{if .Message}}{{.Message.Comment}}{{end}}{{end}}
	"{{ .NameSnake }}": {{ template "fieldSchema" . }}    
{{- end }}
}
{{- end -}}

{{- define "fieldSchema" -}}
{{- if eq .Kind "REPEATED_MESSAGE" }}
{
    {{ template "required" . }}
    {{ template "repeatedMessage" . }}
},
{{- end }}

{{- if eq .Kind "REPEATED_ELEMENTARY" }}
{
    {{ template "required" . }}
    {{ template "repeatedElementary" . }}
},
{{- end }}

{{- if eq .Kind "MAP" }}
{
    {{ template "required" . }}
    {{ template "map" . }}
},
{{- end }}

{{- if eq .Kind "OBJECT_MAP" }}
{
    {{ template "required" . }}
    {{ template "objectMap" . }}
},
{{- end }}

{{- if eq .Kind "SINGULAR_MESSAGE" }}
{
    {{ template "required" . }}
    Type: schema.TypeList,
    Description: {{ .Message.RawComment | quote }},
    MaxItems: 1,
    Elem: &schema.Resource {
        Schema: {{ template "fieldsSchema" .Message.Fields }},
    },
},
{{- end }}

{{- if eq .Kind "SINGULAR_ELEMENTARY" }}
{
    {{ template "singularElementary" . }}
    {{ template "required" . }}    
},
{{- end }}

{{- if eq .Kind "CUSTOM_TYPE" }}
Schema{{.CustomTypeMethodInfix}}(),
{{- end }}
{{- end -}}

{{- define "singularElementary" -}}
Type: {{ template "type" .SchemaRawType }},
Description: {{ .RawComment | quote }},
{{- if .IsTime }}
ValidateFunc: validation.IsRFC3339Time,
{{- end }}
{{- if .Default }}
Default: {{.Default | quote}},
{{- end }}
{{- end -}}

{{- define "repeatedMessage" -}}
Type: schema.TypeList,
Description: {{ .Message.RawComment | quote }},
Elem: &schema.Resource {
    Schema: {{ template "fieldsSchema" .Message.Fields }},
},
{{- end -}}

{{- define "repeatedElementary" -}}
Type: schema.TypeList,
Description: {{ .RawComment | quote }},
Elem: &schema.Schema {
    Type: {{ template "type" .SchemaRawType }},
},
{{- end -}}

{{- define "map" -}}
Type: schema.TypeMap,
Description: {{ .RawComment | quote }},
Elem: &schema.Schema {
    Type: {{ template "type" .MapValueField.SchemaRawType }},
},
{{- end -}}

{{- define "objectMap" -}}
Type: schema.TypeList,
Description: {{ .MapValueField.RawComment | quote }},
Elem: &schema.Resource {
    Schema: map[string]*schema.Schema{
        "key": {
            Type: schema.TypeString,
            Required: true,
        },
        "value": {{ template "fieldSchema" .MapValueField }}
    },
},
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

{{- define "required" -}}
{{- if .IsComputed -}}
Computed: true,
{{- else -}}
{{- if .IsRequired -}}
Required: true,
{{- else -}}
Optional: true,
{{- end }}
{{- end -}}
{{- end -}}
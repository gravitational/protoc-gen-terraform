{{/* ---- Schema rendering ---------------------------------------------------------------*/}}
// Schema{{ .Name }} returns schema for {{.Name}}
//
{{.Comment}}
func GenSchema{{ .Name }}() map[string]*schema.Schema {
	return {{ template "fieldsSchema" .Fields -}}
}

{{- define "fieldsSchema" -}}
map[string]*schema.Schema {
{{- range $index, $field := . }}
    {{.Comment}}
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

{{- if eq .Kind "MESSSAGE_MAP" }}
{
    {{ template "required" . }}
    {{ template "messageMap" . }}
},
{{- end }}

{{- if eq .Kind "SINGULAR_MESSAGE" }}
{
    Type: schema.TypeList,
    MaxItems: 1,
    Description: {{ .Message.RawComment | quote }},
    {{ template "configMode" . -}}    
    {{- template "required" . }}
    Elem: &schema.Resource {
        Schema: {{ template "fieldsSchema" .Message.Fields }},
    },
},
{{- end }}

{{- if eq .Kind "SINGULAR_ELEMENTARY" }}
{
    {{ template "singularElementary" . }}
    {{- template "required" . }}    
},
{{- end }}

{{- if eq .Kind "CUSTOM_TYPE" }}
Schema{{.Suffix}}(),
{{- end }}
{{- end -}}

{{- define "singularElementary" -}}
{{- template "configMode" . -}}
Type: {{ template "type" .SchemaRawType }},
Description: {{ .RawComment | quote }},
{{- if .IsTime }}
ValidateFunc: validation.IsRFC3339Time,
{{- end }}
{{- if .IsDuration }}
DiffSuppressFunc: SuppressDurationChange,
{{- end }}
{{- if .StateFunc }}
StateFunc: {{ .StateFunc }},
{{- end }}
{{- end -}}

{{- define "repeatedMessage" -}}
Type: schema.TypeList,
Description: {{ .RawComment | quote }},
Elem: &schema.Resource {
    Schema: {{ template "fieldsSchema" .Message.Fields }},
},
{{ template "configMode" . -}}
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

{{- define "messageMap" -}}
Type: schema.TypeSet,
Description: {{ .RawComment | quote }},
Elem: &schema.Resource {
    Schema: map[string]*schema.Schema{
        "key": {
            Type: schema.TypeString,
            Required: true,
        },
        "value": {{ template "fieldSchema" .MapValueField }}
    },
},
{{ template "configMode" . -}}
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
{{- if .IsRequired }}
Required: true,
{{- else }}
Optional: true,
{{- end }}
{{- if .IsComputed }}
Computed: true,
{{- end }}
{{- if .Default }}
Default: {{.Default }},
{{- end }}
{{- if .IsForceNew }}
ForceNew: true,
{{- end }}
{{- end -}}

{{- define "configMode" -}}
{{- if .ConfigMode }}
ConfigMode: schema.{{.ConfigMode}},
{{- end }}
{{- end -}}
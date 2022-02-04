var (
{{- range $index, $message := . }}
    // Schema{{.Name}} is schema for {{.RawComment}}
    Schema{{.Name}} = GenSchema{{.Name}}()
    // SchemaMeta{{.Name}} is schema metadata for {{.RawComment}}
    SchemaMeta{{.Name}} = GenSchemaMeta{{.Name}}()
    
{{- end }}
)

// SuppressDurationChange supresses change for equal durations written differently, ex.: "1h" and "1h0m"
func SuppressDurationChange(k string, old string, new string, d *schema.ResourceData) bool {
    o, err := time.ParseDuration(old)
    if err != nil {
        return false
    }

    n, err := time.ParseDuration(new)
    if err != nil {
        return false
    }

    return o == n
}

{{- range $index, $message := . }}
func FromTerraform{{.Name}}(data *schema.ResourceData, obj *{{.GoTypeName}}) error {
    return accessors.FromTerraform(obj, data, Schema{{.Name}}, SchemaMeta{{.Name}})
}

func ToTerraform{{.Name}}(obj *{{.GoTypeName}}, data *schema.ResourceData) error {
    return accessors.ToTerraform(obj, data, Schema{{.Name}}, SchemaMeta{{.Name}})
}
{{- end }}
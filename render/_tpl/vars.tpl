var (
{{- range $index, $message := . }}
    // Schema{{.Name}} is schema for {{.RawComment}}
    Schema{{.Name}} = GenSchema{{.Name}}()

    // SchemaMeta{{.Name}} is schema metadata for {{.RawComment}}
    SchemaMeta{{.Name}} = GenSchemaMeta{{.Name}}()
{{- end }}
)

// SchemaMeta represents schema metadata struct
type SchemaMeta struct {
	name        string
	isTime      bool
	isDuration  bool
	nested      map[string]*SchemaMeta
}

// SupressDurationChange supresses change for equal durations written differently, ex.: "1h" and "1h0m"
func SupressDurationChange(k string, old string, new string, d *schema.ResourceData) bool {
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
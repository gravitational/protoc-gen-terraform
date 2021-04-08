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

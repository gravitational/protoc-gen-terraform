func Unmarshal{{.Name}}(d *schema.ResourceData, t *{{.GoTypeName}}) error {
	meta := []struct {
		Name          string
		SchemaName    string
		SchemaRawType string
		SchemaGoType  string
	}{
        {{ range $index, $field := .Fields }}
            {
                Name: {{$field.Name | quote}},
                SchemaName: {{$field.NameSnake | quote}},
                SchemaRawType: {{$field.SchemaRawType | quote}},
                SchemaGoType: {{$field.SchemaGoType | quote}},
            },
        {{ end }}
    }
}
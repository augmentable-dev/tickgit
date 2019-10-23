package todos

import (
	"io"
	"text/template"
)

// DefaultTemplate is the default report template
const DefaultTemplate = `
{{- range . }}
=== 📋 {{ .String }}
  --- {{ .FilePath }}:{{ .Line }}:{{ .Position }}
{{- else }}
no todos 🎉
{{- end }}
`

// WriteTodos renders a report of todos
func WriteTodos(todos []*TODO, writer io.Writer) error {

	t, err := template.New("todos").Parse(DefaultTemplate)
	if err != nil {
		return err
	}

	t.Execute(writer, todos)

	return nil
}

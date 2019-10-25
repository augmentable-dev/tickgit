package todos

import (
	"io"
	"text/template"
)

// DefaultTemplate is the default report template
const DefaultTemplate = `
{{- range $index, $todo := . }}
{{ print "\u001b[33m" }}TODO{{ print "\u001b[0m" }}{{ .String }}
  => {{ with .Comment }}{{ .FilePath }}:{{ .StartLocation.Line }}:{{ .StartLocation.Pos }}{{ end }}
{{ else }}
no todos 🎉
{{- end }}
{{ .Count }} TODOs Found 📝
`

// WriteTodos renders a report of todos
func WriteTodos(todos ToDos, writer io.Writer) error {

	t, err := template.New("todos").Parse(DefaultTemplate)
	if err != nil {
		return err
	}

	err = t.Execute(writer, todos)
	if err != nil {
		return err
	}

	return nil
}

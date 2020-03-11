package todos

import (
	"io"
	"strings"
	"text/template"
)

const defaultTemplate = `
{{- range $index, $todo := . }}
{{ .String }}
  => {{ .Comment.FilePath }}:{{ .Comment.StartLocation.Line }}:{{ .Comment.StartLocation.Pos }}
  {{- if .Blame }}
  => added {{ .TimeAgo }} by {{ .Blame.Author }} in {{ .Blame.SHA }}
  {{- end }}
{{ else }}
no todos 🎉
{{- end }}
{{ len . }} TODOs Found 📝
`

// WriteTodos renders a report of todos
func WriteTodos(todos ToDos, writer io.Writer) error {

	t, err := template.New("todos").Parse(defaultTemplate)
	if err != nil {
		return err
	}

	// replace the phrase in the todo string with a "highlighted" version for console output
	// TODO eventually make this configurable, for NO_COLOR output (or customization of color?)
	for _, todo := range todos {
		todo.String = strings.Replace(todo.String, todo.Phrase, "\u001b[33m"+todo.Phrase+"\u001b[0m", 1)
	}

	err = t.Execute(writer, todos)
	if err != nil {
		return err
	}

	return nil
}

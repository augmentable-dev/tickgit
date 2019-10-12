package tickgit

import (
	"io"
	"text/template"

	"gopkg.in/src-d/go-git.v4/plumbing/object"
)

// DefaultTemplate is the default report template
const DefaultTemplate = `
{{- range . }}
=== {{ .Title }} {{ if .Completed }}✅{{ else }}⏳{{ end }}
  --- {{ .Summary.Completed }}/{{.Summary.Total}} tasks completed ({{ .Summary.Pending }} remaining)
  --- {{ .Summary.PercentCompleted }}% completed

  {{ range .Tasks }}
  {{- if .Completed }}✅{{ else }}⏳{{ end }} {{ .Title }}:
    {{- if not .Description -}}{{ else }}
    > {{ .Description }}
    {{ end }}
  {{ else }}
  {{- end }}
{{- else }}
no goals
{{- end }}
`

// WriteStatus renders a status report to the passed in writer
func WriteStatus(commit *object.Commit, writer io.Writer) error {
	goals, err := GoalsFromCommit(commit, nil)
	if err != nil {
		return err
	}

	t, err := template.New("status").Parse(DefaultTemplate)
	if err != nil {
		return err
	}

	t.Execute(writer, goals)

	return nil
}

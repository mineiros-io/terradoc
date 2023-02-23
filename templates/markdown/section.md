{{define "section"}}{{if .Title}}{{repeat "#" .Level}} {{.Title}}{{- end -}}
{{if .Content}}{{- newline -}}{{.Content}}{{end}}
{{ end }}

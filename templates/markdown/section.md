{{define "section"}}{{if .Title}}{{repeat "#" .Level}} {{.Title}}{{- end -}}
{{- newline -}}
{{if .Content}}{{.Content}}{{- newline -}}{{end}}{{- end -}}

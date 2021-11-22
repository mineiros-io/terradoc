{{- define "section"}}{{- newline -}}{{- repeat "#" .Level}} {{.Title -}}

{{- if .Description}}{{- newline}}{{.Description}}{{end}}{{- newline -}}{{- end -}}

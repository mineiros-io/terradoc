{{- define "section"}}
{{- repeat "#" .Level}} {{.Title -}}

{{if .Description}}{{newline}}{{.Description}}{{end}}{{- newline -}}
{{end -}}

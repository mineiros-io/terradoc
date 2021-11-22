{{- define "toc"}}
{{- range . }}
{{ indent .IndentLevel "-"}} [{{.Label}}](#{{.Value}}){{end -}}
{{- newline -}}{{- end -}}

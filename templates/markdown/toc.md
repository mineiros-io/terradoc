{{- define "toc"}}{{- range . }}
{{ indent .IndentLevel "-"}} [{{.Label}}](#{{urlfragment .Label}}){{end -}}
{{- newline -}}{{- end -}}

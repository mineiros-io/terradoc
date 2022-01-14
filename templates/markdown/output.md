{{define "output"}}- [**`{{.Name}}`**](#output-{{.Name}}): *(`{{template "variableType" .Type}}`)*<a name="output-{{.Name}}"></a>

{{- if .Description}}{{- newline}}{{indent 2 .Description}}{{end}}
{{- newline -}}
{{end}}

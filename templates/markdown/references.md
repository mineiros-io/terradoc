{{define "references"}}
<!-- References -->
{{range . }}
[{{.Name}}]: {{.Value}}{{end}}

{{end}}

{{define "variable"}}- **`{{.Name}}`**: *({{if .Required}}**Required**{{else}}Optional{{end}} `{{template "variableType" .Type}}`{{if .ForcesRecreation}}, Forces new resource{{end}})*

{{- if .Description}}{{- newline}}  {{print .Description}}{{- end}}

{{- if .Default}}{{- newline}}  Default is `{{printf "%s" .Default}}`.{{- end}}

{{- if .ReadmeExample}}{{- newline}}  Example:

{{printf "```hcl\n%s\n```" .ReadmeExample | indent 2}}{{end -}}
{{- newline -}}
{{end}}

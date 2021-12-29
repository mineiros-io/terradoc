{{define "variable"}}- [**`{{.Name}}`**](#var-{{.Name}}): *({{if .Required}}**Required**{{else}}Optional{{end}} `{{template "variableType" .Type}}`{{if .ForcesRecreation}}, Forces new resource{{end}})*<a name="var-{{.Name}}"></a>

{{- if .Description}}{{- newline}}{{indent 2 .Description}}{{end}}

{{- if .Default}}{{- newline}}  Default is `{{printf "%s" .Default}}`.{{- end}}

{{- if .ReadmeExample}}{{- newline}}  Example:

{{printf "```hcl\n%s\n```" .ReadmeExample | indent 2}}{{end -}}
{{- newline -}}
{{end}}

{{define "attribute"}}{{indent (multiply .Level 2) "-"}} [**`{{.Name}}`**](#attr-{{.Name}}-{{.ParentName}}): *({{if .Required}}**Required**{{else}}Optional{{end}} `{{template "variableType" .Type}}`{{if .ForcesRecreation}}, Forces new resource{{end}})*<a name="attr-{{.Name}}-{{.ParentName}}"></a>

{{- if .Description}}{{- newline}}{{indent (getIndent .Level) .Description}}{{end}}

{{- if .Default}}{{- newline}}{{indent (getIndent .Level) "Default"}} is `{{printf "%s" .Default}}`.{{end -}}

{{- if .ReadmeExample}}{{- newline}}{{indent (getIndent .Level) "Example:"}}

{{printf "```hcl\n%s\n```" .ReadmeExample | indent (getIndent .Level)}}{{end -}}
{{- newline -}}
{{end}}

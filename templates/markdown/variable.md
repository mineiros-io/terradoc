{{- define "variable"}}- **`{{.Name}}`**: *({{if .Required}}**Required**{{else}}Optional{{end}} `{{if .Type.ReadmeType}}{{.Type.ReadmeType}}{{else}}{{.Type.TerraformType.Type}}{{end}}`{{if .ForcesRecreation}}, Forces new resource{{end}})*

{{- if .Description}}{{- newline}}  {{print .Description}}{{- end}}

{{- if .Default}}{{- newline}}  Default is `{{printf "%s" .Default}}`.{{- end}}

{{- if .ReadmeExample}}{{- newline}}  Example:

{{printf "```terraform\n%s\n```" .ReadmeExample | indent 2}}{{end -}}
{{- newline -}}{{- end -}}

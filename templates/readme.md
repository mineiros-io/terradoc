{{- range .Sections}}{{- template "section" . }}{{end}}

{{- define "section"}}{{repeat "#" .Level}} {{.Title}}

{{- if .Description}}

{{ indent 2 .Description }}
{{- end}}

{{range .Variables}}{{include "variable" . }}{{end}}
{{range .Sections}}{{include "section" . }}{{end}}

{{- end}}

{{- define "variable"}}- **`{{.Name}}`**: *({{if .Required}}Required{{else}}Optional{{end}} `{{if .ReadmeType}}{{.ReadmeType}}{{else}}{{.Type}}{{end}}`{{if .ForcesRecreation}}, Forces new resource{{end}})*
{{ include "variableDescription" . | indent 2 }}
{{end}}


{{define "variableDescription" }}
{{.Description}}

{{if .Default}}Default is `{{.Default}}`.{{- end}}

{{- if .ReadmeExample}}

```terraform
{{.ReadmeExample | indent 2}}
```

{{end}}

{{- if .HasAttributes}}
{{- if .IsCollection}}Each element of {{end}}`{{.ReadmeType}}` is an object with the following attributes:

{{- range .Attributes}}
{{include "attribute" .}}
{{- end}}
{{- end}}
{{- end}}

{{- define "attribute"}}
- **`{{.Name}}`**: *({{if .Required}}Required{{else}}Optional{{end}} `{{.Type}}`{{if .ForcesRecreation}}, Forces new resource{{end}})*

{{indent 2 .Description}}

{{- range .Attributes}}
{{include "attribute" . | indent (multiply .Level 2)}}
{{- end}}

{{end}}

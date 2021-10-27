{{define "typeDescription"}}
{{- if .TerraformType.HasNestedType }}{{indent .IndentLevel "`"}}{{.ReadmeType}}` is a `{{.TerraformType.Type}}` of `{{.TerraformType.NestedType}}` with the following attributes:{{newline}}
{{- else -}}
{{indent .IndentLevel "`" }}{{.ReadmeType}}` is a `{{.TerraformType.Type}}` with the following attributes:{{newline}}{{end}}
{{- end -}}

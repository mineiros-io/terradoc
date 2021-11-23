{{define "typeDescription"}}
{{- if .TerraformType.HasNestedType }}{{indent .IndentLevel "Each object in the"}} {{.TerraformType.Type}} accepts the following attributes:{{newline}}
{{- else -}}
{{indent .IndentLevel "The object accepts the following attributes:" }}{{newline}}{{end}}
{{- end -}}

{{- define "variableType" -}}
{{- if .ReadmeType -}}
  {{- .ReadmeType -}}
{{- else -}}
    {{- if .TerraformType.HasNestedType -}}
        {{- .TerraformType.Type -}}({{.TerraformType.NestedType}})
    {{- else -}}
        {{- .TerraformType.Type -}}
    {{- end -}}
{{- end -}}
{{- end -}}

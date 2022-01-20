{{- define "variableType" -}}
    {{- if .HasNestedType -}}
        {{template "nestedVariableType" .}}
    {{- else -}}
        {{- if .Label -}}
            {{- .TFType -}}({{ .Label }})
        {{- else -}}
            {{- .TFType -}}
        {{- end -}}
    {{- end -}}
{{- end -}}

{{- define "nestedVariableType" -}}
    {{- if .Nested.Label -}}
        {{- .TFType -}}({{.Nested.Label}})
    {{- else -}}
        {{- .TFType -}}({{.Nested.TFType}})
    {{- end -}}
{{- end -}}

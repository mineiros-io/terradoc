{{- define "variableType" -}}
    {{- if .HasNestedType -}}
        {{template "nestedVariableType" .}}
    {{- else -}}
        {{- if .TFTypeLabel -}}
            {{- .TFType -}}({{ .TFTypeLabel }})
        {{- else -}}
            {{- .TFType -}}
        {{- end -}}
    {{- end -}}
{{- end -}}

{{- define "nestedVariableType" -}}
    {{- if .NestedTFTypeLabel -}}
        {{- .TFType -}}({{.NestedTFTypeLabel}})
    {{- else -}}
        {{- .TFType -}}({{.NestedTFType}})
    {{- end -}}
{{- end -}}

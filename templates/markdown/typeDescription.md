{{define "typeDescription"}}
    {{- if .HasNestedType }}{{indent .IndentLevel "Each"}} `{{.NestedTFTypeLabel}}` object in the {{.Type.TFType}} accepts the following attributes:{{newline}}
    {{- else -}}
        {{- if .TFTypeLabel -}}
            {{indent .IndentLevel "The"}} `{{.TFTypeLabel}}` object accepts the following attributes:{{newline}}
        {{- else -}}
            {{indent .IndentLevel "The object accepts the following attributes:" }}{{newline}}
        {{- end -}}
    {{- end -}}
{{- end -}}

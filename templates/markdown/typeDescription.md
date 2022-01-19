{{define "typeDescription"}}
    {{- if .HasNestedType }}{{indent .IndentLevel "Each"}} `{{.Nested.Label}}` object in the {{.Type.TFType}} accepts the following attributes:{{newline}}
    {{- else -}}
        {{- if .Label -}}
            {{indent .IndentLevel "The"}} `{{.Label}}` object accepts the following attributes:{{newline}}
        {{- else -}}
            {{indent .IndentLevel "The object accepts the following attributes:" }}{{newline}}
        {{- end -}}
    {{- end -}}
{{- end -}}
